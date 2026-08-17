# benzhi.Dockerfile —— 保留完整 Go 工具链的验证镜像（禁止多阶段只留二进制）
FROM golang:1.22

# 工具链固定为镜像自带的 1.22，不允许自动下载更高版本
ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=1

WORKDIR /src
COPY . /src

# 本项目仅用标准库，无 go.sum，故不需要 go mod download
RUN go build ./...

CMD ["bash"]
