#!/usr/bin/env bash
#
# 验证 GOPROXY 逗号(,)与管道(|)分隔符回退语义的自动化测试。
#
#   ,  -> 仅当上一代理返回 404/410(模块不存在)时回退到下一个;
#         其余错误(如 500、连接失败)立即停止。
#   |  -> 上一代理返回任何错误都继续尝试下一个。
#
# 断言依据(交叉验证):
#   1) go mod download 的退出码是否符合预期(逗号模式应失败、管道模式应成功);
#   2) "健康代理"(200)是否真的收到了请求(逗号模式不应收到、管道模式应收到)。
#
# 环境要求: go(推荐 1.21+)、bash。脚本只访问 127.0.0.1,不依赖外网。

set -u # 不用 set -e,便于逐个用例断言并汇总

# 使用本机已安装的 go(脚本目录同级若存在 .local/go 则优先)
if [ -x "$HOME/.local/go/bin/go" ]; then
    export PATH="$HOME/.local/go/bin:$PATH"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR" || { echo "无法进入脚本目录: $SCRIPT_DIR" >&2; exit 1; }

# ---- 配置 ---------------------------------------------------------------
P1_PORT=18081; P1_STATUS=404 # 模拟"该代理上找不到模块"
P2_PORT=18082; P2_STATUS=500 # 模拟"代理服务端出错"
P3_PORT=18083; P3_STATUS=200 # 模拟"健康代理",能真实提供模块
P1_URL="http://127.0.0.1:$P1_PORT"
P2_URL="http://127.0.0.1:$P2_PORT"
P3_URL="http://127.0.0.1:$P3_PORT"
TARGET="example.com/dep@v1.0.0"

PASS=0
FAIL=0
WORK="$(mktemp -d "${TMPDIR:-/tmp}/goproxy-sep-test.XXXXXX")"
PIDFILE="$WORK/pids" # 记录本脚本启动的所有代理进程

say() { printf '%s\n' "$*"; }

# 精确结束由本脚本启动的代理进程(按 PID,不误伤其它进程)
stop_proxies() {
    if [ -f "$PIDFILE" ]; then
        while read -r pid; do
            if [ -n "$pid" ]; then
                kill "$pid" 2>/dev/null || true
            fi
        done <"$PIDFILE"
        sleep 0.3
        : > "$PIDFILE"
    fi
}

# 清理临时目录。go 未使用 -modcacherw 时,模块缓存里是只读文件(0444/0555),
# 直接 rm -rf 会 Permission denied,必须先放开写权限再删。
cleanup() {
    stop_proxies
    if [ -n "${WORK:-}" ] && [ -d "$WORK" ]; then
        chmod -R u+w "$WORK" 2>/dev/null || true
        rm -rf "$WORK"
    fi
}
# 任何退出路径(含编译失败、用例失败、Ctrl-C)都要清理
trap cleanup EXIT INT TERM

# 等待某端口就绪(bash 内置 /dev/tcp,兼容 macOS)
wait_port() {
    local port="$1"
    for _ in $(seq 1 100); do
        if (exec 3<>"/dev/tcp/127.0.0.1/$port") 2>/dev/null; then
            exec 3>&- 3<&- 2>/dev/null || true
            return 0
        fi
        sleep 0.1
    done
    return 1
}

# 端口是否已被占用(能连通即视为占用)
port_busy() {
    (exec 3<>"/dev/tcp/127.0.0.1/$1") 2>/dev/null
}

# 启动 3 个代理,把各自日志写到 run 目录,并把 PID 追加到 PIDFILE
start_proxies() {
    local run="$1"
    local pid
    PROXY_NAME=p1 PROXY_PORT=$P1_PORT PROXY_STATUS=$P1_STATUS ./proxy-server >"$run/p1.log" 2>&1 &
    pid=$!; echo "$pid" >>"$PIDFILE"
    wait_port $P1_PORT || { say "     [错误] p1 未能在 127.0.0.1:$P1_PORT 就绪"; return 1; }
    PROXY_NAME=p2 PROXY_PORT=$P2_PORT PROXY_STATUS=$P2_STATUS ./proxy-server >"$run/p2.log" 2>&1 &
    pid=$!; echo "$pid" >>"$PIDFILE"
    wait_port $P2_PORT || { say "     [错误] p2 未能在 127.0.0.1:$P2_PORT 就绪"; return 1; }
    PROXY_NAME=p3 PROXY_PORT=$P3_PORT PROXY_STATUS=$P3_STATUS ./proxy-server >"$run/p3.log" 2>&1 &
    pid=$!; echo "$pid" >>"$PIDFILE"
    wait_port $P3_PORT || { say "     [错误] p3 未能在 127.0.0.1:$P3_PORT 就绪"; return 1; }
    return 0
}

# 在独立临时模块目录里跑 go mod download(钉死版本),返回其退出码。
# 用版本号(而非 @latest)保证 go 只取 .info/.mod/.zip,请求序列规整可控,
# 回退行为最直观;go get 还会做版本解析与父模块探测,请求更杂、还会改 go.mod。
# 在子 shell 中隔离 GOMODCACHE/GOPATH/GOCACHE/GOENV,避免污染共享缓存与
# 外层 shell 环境;保证每个用例都会真正访问本地代理(而不是命中缓存)。
run_mod_download() {
    local moddir="$1" cfg="$2" outfile="$3"
    mkdir -p "$moddir"
    cp testmodule/go.mod "$moddir/go.mod"
    : > "$moddir/goenv"
    (
        export GOPROXY="$cfg"
        export GOSUMDB=off
        export GOINSECURE="*"
        export GOFLAGS="-mod=mod"
        export GOTOOLCHAIN=local
        export GONOPROXY='' GONOSUMDB='' GOPRIVATE=''
        export GOMODCACHE="$moddir/gomodcache"
        export GOPATH="$moddir/gopath"
        export GOCACHE="$moddir/gocache"
        export GOENV="$moddir/goenv"
        cd "$moddir" && go mod download "$TARGET"
    ) >"$outfile" 2>&1
    return $?
}

count_requests() { # 统计某个代理日志里收到的 GET 请求数
    local log="$1"
    if [ -f "$log" ]; then
        grep -c "GET " "$log" 2>/dev/null || true
    else
        echo 0
    fi
}

# 单个用例: title=描述, cfg=GOPROXY, expect=ok|fail
run_case() {
    local title="$1" cfg="$2" expect="$3" run
    run="$WORK/$(printf 'case_%02d' "$((PASS + FAIL + 1))")"
    mkdir -p "$run"

    say "-----------------------------------------------------------------"
    say "用例: $title"
    say "  GOPROXY=$cfg"
    say "  预期: go mod download $([ "$expect" = ok ] && echo 成功 || echo 失败)"

    if ! start_proxies "$run"; then
        stop_proxies
        FAIL=$((FAIL + 1))
        say "  [失败] 代理启动异常"
        return
    fi

    run_mod_download "$run/mod" "$cfg" "$run/go.log"
    local ec=$?
    stop_proxies

    # 交叉验证:健康代理 p3(200)是否收到请求
    local n1 n2 n3 p3_hit
    n1=$(count_requests "$run/p1.log")
    n2=$(count_requests "$run/p2.log")
    n3=$(count_requests "$run/p3.log")
    [ "$n3" -gt 0 ] && p3_hit=1 || p3_hit=0

    local expect_ok=0
    [ "$expect" = ok ] && expect_ok=1

    # 断言:退出码符合预期,且 p3 可达性与预期一致
    local ec_ok=0 p3_ok=0
    [ "$ec" -eq 0 ] && [ "$expect_ok" -eq 1 ] && ec_ok=1
    [ "$ec" -ne 0 ] && [ "$expect_ok" -eq 0 ] && ec_ok=1
    [ "$p3_hit" -eq "$expect_ok" ] && p3_ok=1

    if [ "$ec_ok" -eq 1 ] && [ "$p3_ok" -eq 1 ]; then
        PASS=$((PASS + 1))
        say "  [通过] go mod download 退出码=$ec(符合预期)"
        say "         健康代理 p3 收到请求数=$n3(符合预期: $([ "$expect_ok" = 1 ] && echo 应收到 || echo 不应收到))"
    else
        FAIL=$((FAIL + 1))
        say "  [未通过] 期望=$expect,实际退出码=$ec"
        say "  p1(404)请求数=$n1  p2(500)请求数=$n2  p3(200)请求数=$n3"
    fi

    # 打印证据:预期失败时展示 go 的报错;预期成功时展示 p3 收到的请求
    if [ "$expect" = fail ]; then
        say "  go 报错(末行): $(tail -n 1 "$run/go.log")"
    else
        say "  健康代理 p3 收到的请求:"
        grep -E "GET " "$run/p3.log" | sed 's/^/    /' || true
    fi
    say ""
}

# ---- 编译代理服务器 -------------------------------------------------------
say "编译代理服务器..."
if ! go build -o proxy-server main.go; then
    say "编译失败,请确认已安装 go。"
    exit 1
fi
say "编译完成。"
say ""

# 端口占用预检:残留进程会让测试连到过期代理,结果失真,直接报错退出
say "检查端口是否可用..."
for p in $P1_PORT $P2_PORT $P3_PORT; do
    if port_busy "$p"; then
        say "端口 $p 已被占用,可能残留 proxy-server 进程。请先清理再运行:"
        say "  lsof -nP -iTCP:$p -sTCP:LISTEN"
        exit 1
    fi
done
say "端口均可用。"
say ""

# ---- 用例 ---------------------------------------------------------------
# p1(404): 模块不存在; p2(500): 代理服务端出错; p3(200): 健康代理可提供模块。

# 逗号: p1 404 -> 回退 p2; p2 500 -> 逗号模式应停止
run_case "逗号模式: 404 后遇到 500 应停止" \
    "$P1_URL,$P2_URL,$P3_URL" "fail"

# 管道: p1 404 -> 回退 p2; p2 500 -> 继续 p3 -> 成功
run_case "管道模式: 404 后遇到 500 应继续尝试" \
    "$P1_URL|$P2_URL|$P3_URL" "ok"

# 逗号: 列表首位 p2 直接 500 -> 应立刻停止
run_case "逗号模式: 首位代理即 500 应停止" \
    "$P2_URL,$P1_URL,$P3_URL" "fail"

# 管道: 首位 p2 500 -> 继续 p1(404) -> 继续 p3 -> 成功
run_case "管道模式: 首位代理 500 应继续尝试" \
    "$P2_URL|$P1_URL|$P3_URL" "ok"

# ---- 汇总 ---------------------------------------------------------------
say "=============================================================="
say "结果: 通过 $PASS / $((PASS + FAIL))"
say "=============================================================="

cleanup # 显式清理;EXIT trap 会再跑一次(幂等,无副作用)

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
exit 0
