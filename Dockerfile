ARG ZSW_CHAIN_LISHI_DEB_URL="https://github.com/invisible-train-40/zswchain-lishi/releases/download/2.0.8-prod-1.1.0/zswchain-lishi_2.0.8-dm.12.0_amd64.deb"
FROM ubuntu:18.04 AS base
ARG ZSW_CHAIN_LISHI_DEB_URL
RUN apt update && apt-get -y install curl ca-certificates libicu60 libusb-1.0-0 libcurl3-gnutls
RUN mkdir -p /var/cache/apt/archives/
RUN echo "dl $ZSW_CHAIN_LISHI_DEB_URL" && curl -sL -o/var/cache/apt/archives/zswchain.deb "$ZSW_CHAIN_LISHI_DEB_URL"
RUN dpkg -i /var/cache/apt/archives/zswchain.deb
RUN rm -rf /var/cache/apt/*
RUN rm -rf /var/cache/apt/*

FROM node:12 AS dlauncher
WORKDIR /work
ADD go.mod /work
RUN apt update && apt-get -y install git
RUN cd /work && git clone https://github.com/streamingfast/dlauncher.git dlauncher &&\
    grep -w github.com/streamingfast/dlauncher go.mod | sed 's/.*-\([a-f0-9]*$\)/\1/' |head -n 1 > dlauncher.hash &&\
    cd dlauncher &&\
    git checkout "$(cat ../dlauncher.hash)" &&\
    cd dashboard/client &&\
    yarn install && yarn build

FROM node:12 AS eosq
ADD eosq /work
WORKDIR /work
RUN yarn install && yarn build

FROM golang:1.14 as dfuse
ARG COMMIT
ARG VERSION
RUN go get -u github.com/GeertJohan/go.rice/rice && export PATH=$PATH:$HOME/bin:/work/go/bin
RUN mkdir -p /work/build
ADD . /work
WORKDIR /work
COPY --from=eosq      /work/ /work/eosq
# The copy needs to be one level higher than work, the dashboard generates expects this file layout
COPY --from=dlauncher /work/dlauncher /dlauncher
RUN cd /dlauncher/dashboard && go generate
RUN cd /work/eosq/app/eosq && go generate
RUN cd /work/dashboard && go generate
RUN cd /work/dgraphql && go generate
RUN go test ./...
RUN go build -ldflags "-s -w" -v -o /work/build/dfuseeos ./cmd/dfuseeos

FROM base
RUN mkdir -p /app/ && curl -Lo /app/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.2.2/grpc_health_probe-linux-amd64 && chmod +x /app/grpc_health_probe
COPY --from=dfuse /work/build/dfuseeos /app/dfuseeos
COPY --from=dfuse /work/tools/manageos/motd /etc/motd
COPY --from=dfuse /work/tools/manageos/scripts /usr/local/bin/
RUN echo cat /etc/motd >> /root/.bashrc
