FROM ubuntu:18.04 as base
# 環境変数
ENV CMAKE_INSTALL_DIR /local
ENV PATH $PATH:/grpc/cmake/build
ENV PATH $PATH:/grpc/cmake/build/third_party/protobuf
ENV PATH $PATH:$CMAKE_INSTALL_DIR/bin

RUN apt-get update && \
    apt-get install -y \
    git \
    wget \
    clang \
    ca-certificates \
    build-essential \
    libssl-dev \
    make \
    autoconf \
    automake \
    pkg-config \
    libtool \
    vim \
    cron \
    netcat \
    cmake \
    libboost-all-dev \
    curl

ENV PATH $PATH:/root/go/bin
ENV PATH $PATH:/usr/local/go/bin

# Install go1.17.3
RUN wget https://dl.google.com/go/go1.17.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.17.3.linux-amd64.tar.gz

# grpcサーバのhealthcheckをするためのツールを
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.14 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

FROM base as src
# Configの切り替えに使用する環境変数
ENV STAGE=src
# goの自動成型ツール
RUN go get -v github.com/ramya-rao-a/go-outline
# goの自動成型ツールimportが勝手に消されるのでお好みかも．．．
RUN go get -v golang.org/x/tools/gopls

# Install Protocolbuf
RUN apt-get install -y unzip curl && \
    PROTOC_VER=3.18.0 && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip && \
    unzip protoc-${PROTOC_VER}-linux-x86_64.zip -d /root/.local && \
    go get -v google.golang.org/protobuf/cmd/protoc-gen-go && \
    go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN wget https://github.com/fullstorydev/grpcurl/releases/download/v1.8.5/grpcurl_1.8.5_linux_x86_64.tar.gz && \
    tar -C /usr/local/bin -xzf grpcurl_1.8.5_linux_x86_64.tar.gz && \
    rm -f grpcurl_1.8.5_linux_x86_64.tar.gz

ENV PATH $PATH:/usr/local
ENV PATH $PATH:/root/.local/bin

# デバッグ用の環境変数
# ENV GODEBUG=http2debug=1
# ENV GRPC_TRACE=all

WORKDIR /QuickMPC-BTS

FROM base as builder
COPY src/BeaverTripleService/ /QuickMPC-BTS
# github actionsで使用されているdockerがCOPYを2連続で行うとたまにエラーが出る問題を抱えている
# その問題を回避するために、COPYの間にRun trueを挟んでいる
RUN true
COPY src/Proto/ /Proto
RUN true
COPY Config/ /QuickMPC-BTS/Config
WORKDIR /QuickMPC-BTS
RUN go mod vendor && go build

FROM builder as dev
WORKDIR /QuickMPC-BTS

# image for deploy
FROM alpine:latest as dep-runner
COPY --from=builder /QuickMPC-BTS/BeaverTripleService /QuickMPC-BTS/BeaverTripleService
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /QuickMPC-BTS

FROM builder as small
WORKDIR /QuickMPC-BTS

FROM builder as medium
WORKDIR /QuickMPC-BTS

FROM alpine:latest as large-runner
COPY --from=builder /QuickMPC-BTS/BeaverTripleService /QuickMPC-BTS/BeaverTripleService
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /QuickMPC-BTS
