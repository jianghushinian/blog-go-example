### 此项目是 [onex](https://github.com/onexstack/onex) 项目中 [nightwatch](https://github.com/onexstack/onex/tree/master/internal/nightwatch) 组件的 copy 实现

- references: https://github.com/onexstack/onex/tree/master/internal/nightwatch

### 项目功能

此项目实现了定时同步 MariaDB 和 K8s 之间的任务状态。

启动项目后，主要干了两件事：

1. 将 MariaDB 表中 Normal 状态的 task 记录在 K8s 中启动对应的 Job。

2.  同步在 K8s 中已经启动但还未完成的 Job 状态到 MariaDB 表对应的 task 记录中。

### 快速开始

1. 准备一个 K8s 集群，并创建名称为 `demo` 的 namespace

```bash
# 创建 K8s 集群步骤省略
# ...

# 创建 demo namespace
$ kubectl create ns demo
```

2. 使用 docker compose 启动 MariaDB 和 Redis

```bash
$ docker compose -p nightwatch -f nightwatch/assets/docker-compose.yaml up -d
```

3. 在 MariaDB 上执行 `nightwatch/assets/schema.sql` 文件中的 SQL 语句准备测试数据

4. 启动 nightwatch 项目

```bash
$ go run cmd/main.go
```

5. 查看 K8s 中 job 运行情况

```bash
$ kubectl -n demo get job
NAME          COMPLETIONS   DURATION   AGE
demo-task-2   0/1           22s        22s
```
