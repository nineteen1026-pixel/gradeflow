#!/usr/bin/env bash
# build_benzhi_docker.sh —— benzhi 三件套之一（v2：每个架构一个独立 tag）
# 用法：./build_benzhi_docker.sh <镜像名> <平台>
#   ./build_benzhi_docker.sh benzhi-xxx linux/arm64   → 产出 benzhi-xxx:arm64
#   ./build_benzhi_docker.sh benzhi-xxx linux/amd64   → 产出 benzhi-xxx:amd64
# v2 修复：两架构不再共用同一个 tag。共用会被后一次 --load 覆盖，
#          docker run 找不到匹配架构就去拉远端，撞上镜像站限流(429)。
set -euo pipefail

IMG="${1:?用法: ./build_benzhi_docker.sh <镜像名> <平台，如 linux/arm64>}"
PLATFORM="${2:?用法: ./build_benzhi_docker.sh <镜像名> <平台，如 linux/arm64>}"
TAG="${PLATFORM##*/}"                      # arm64 / amd64

grep -qE '^go 1\.22$' go.mod || { echo "go.mod 缺少 'go 1.22'" >&2; exit 1; }
grep -q  '^toolchain '  go.mod && { echo "go.mod 不得含 toolchain 指令" >&2; exit 1; }

echo "构建 $IMG:$TAG ($PLATFORM)"
docker buildx build --platform "$PLATFORM" --load -t "$IMG:$TAG" -f benzhi.Dockerfile .

# 原生架构额外打一个 :latest，方便不带 tag 直接跑
NATIVE="$(uname -m)"
case "$NATIVE" in aarch64|arm64) NATIVE=arm64 ;; x86_64|amd64) NATIVE=amd64 ;; esac
if [ "$TAG" = "$NATIVE" ]; then
  docker tag "$IMG:$TAG" "$IMG:latest"
  echo "  （原生架构，已额外打 tag $IMG:latest）"
fi

echo "built  $IMG:$TAG ($PLATFORM)"
echo "验证：docker run --rm --pull=never --platform $PLATFORM $IMG:$TAG bash -c 'go build ./... && go test ./... && echo CONTAINER_OK'"
