package outputtype

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// OutputType 只存在于 `go test -json` 的事件里，测试代码自身看不到。
// 因此这里分两层：
//  1. TestOutputTypeDemo 生成四种类型的输出；
//  2. TestOutputTypeVerifier 拉起 `go test -json` 子进程解析事件并断言。
//
// Demo 仅在 OUTPUTTYPE_DEMO=1（由 Verifier 注入）时才产生错误，
// 日常 go test ./... 保持全部通过。

func TestOutputTypeDemo(t *testing.T) {
	generate := os.Getenv("OUTPUTTYPE_DEMO") != ""

	t.Log("regular log line")          // OutputType: 空（普通输出）
	fmt.Println("regular stdout line") // OutputType: 空（普通输出）

	if !generate {
		return
	}

	t.Run("fatal", func(t *testing.T) {
		t.Fatal("fatal error stops the subtest") // OutputType: error
	})

	t.Errorf("single line error") // OutputType: error
	t.Errorf("multi\nline error") // 输出两行日志，第一行 OutputType 为 error，续行为 error-continue
}

type testEvent struct {
	Action     string `json:"Action"`
	Test       string `json:"Test"`
	Output     string `json:"Output"`
	OutputType string `json:"OutputType"`
}

func TestOutputTypeVerifier(t *testing.T) {
	cmd := exec.Command("go", "test", "-json", "-count=1", "-run", "^TestOutputTypeDemo$", ".")
	cmd.Env = append(os.Environ(), "OUTPUTTYPE_DEMO=1")
	out, err := cmd.Output()
	// err != nil 是预期的：demo 会失败。构建失败等严重问题表现为无输出。
	if len(out) == 0 {
		t.Fatalf("inner go test produced no JSON: err=%v", err)
	}

	// outputs 按 OutputType 归类收集所有 output 事件的文本
	outputs := make(map[string][]string)
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		var ev testEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			t.Fatalf("invalid JSON event: %v\nline: %s", err, line)
		}
		// error 标记字节（^O/^N）必须已被 test2json 抹掉
		if strings.ContainsAny(ev.Output, "\x0f\x0e") {
			t.Errorf("error markers leaked into output: %q", ev.Output)
		}
		if ev.Action == "output" {
			outputs[ev.OutputType] = append(outputs[ev.OutputType], ev.Output)
		}
	}

	want := map[string][]string{
		"":               {"regular log line", "regular stdout line"}, // 普通输出
		"frame":          {"=== RUN", "--- FAIL"},                     // 测试框架行
		"error":          {"fatal error", "single line error", "multi\n"},
		"error-continue": {"line error"}, // 多行错误的续行
	}
	order := []string{"", "frame", "error", "error-continue"}

	for _, typ := range order {
		for _, substr := range want[typ] {
			if !contains(outputs[typ], substr) {
				t.Errorf("OutputType %q: no output containing %q\ngot: %q", typ, substr, outputs[typ])
			}
		}
		if len(outputs[typ]) > 0 {
			name := typ
			if name == "" {
				name = "(blank)"
			}
			t.Logf("OutputType=%-16s count=%d sample=%q", name, len(outputs[typ]), outputs[typ][0])
		}
	}
	for typ, samples := range outputs {
		if _, ok := want[typ]; !ok {
			t.Errorf("unexpected OutputType %q: %q", typ, samples)
		}
	}
}

func contains(list []string, substr string) bool {
	for _, s := range list {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
