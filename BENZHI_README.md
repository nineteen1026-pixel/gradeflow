# BENZHI_README · gradeflow

## 项目简介

`gradeflow` 是一个作业评审流转服务，纯标准库实现，无外部依赖。主要能力：计分、评分表、评审队列、审计流水与分页。

- module：`gradeflow`
- go.mod：`go 1.22`（无 toolchain 指令）
- 工具链：`GOTOOLCHAIN=local`，容器内固定使用 golang:1.22 自带的 1.22

## 目录

```
cmd/gradeflow/      可执行入口，演示一遍完整业务流程
internal/        业务包，每个包自带单元测试
```

## 标准构建 / 运行 / 测试命令

本机：

```bash
export GOTOOLCHAIN=local
go build ./...
go test ./...
go run ./cmd/gradeflow
```

容器（双架构）：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh benzhi-gradeflow linux/arm64
./build_benzhi_docker.sh benzhi-gradeflow linux/amd64

# 注意用 bash -c，不要用 bash -lc（登录 shell 会重置 PATH 导致找不到 go）
docker run --rm --platform linux/arm64 benzhi-gradeflow \
  bash -c "go build ./... && go test ./... && echo CONTAINER_OK"
```

看到 `CONTAINER_OK` 即为交付物可用。
