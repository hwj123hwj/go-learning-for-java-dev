# 多阶段构建：编译和运行分离。
# 对比 Java 的 Dockerfile：JDK 镜像打包 jar → JRE 镜像运行；
# Go 的优势在这里最直观——编译产物是静态二进制，最终镜像不需要任何运行时。

# ── 阶段1：构建（用完整工具链镜像）─────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 先拷依赖清单再下载：docker layer 缓存让改代码不重新下依赖
COPY go.mod go.sum ./
RUN go mod download

# 再拷源码编译；CGO_ENABLED=0 产出静态二进制，scratch/alpine 都能跑
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /user-api .

# ── 阶段2：运行（只要二进制和静态资源）─────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && adduser -D appuser
USER appuser

WORKDIR /app
COPY --from=builder /user-api .
COPY config.yaml .
COPY static ./static

EXPOSE 8080
ENTRYPOINT ["./user-api"]
