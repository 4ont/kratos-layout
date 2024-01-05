FROM golang:1.20.5-alpine3.18 AS builder
ARG KIT_REPO_PRIVATE_KEY
ARG APP_NAME

ENV NAME=${APP_NAME} \
    #GOPROXY=https://goproxy.cn \
    GO111MODULE="on" \
    GOPRIVATE="github.com/Taskon-xyz"

WORKDIR /data

COPY . .

RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl bash git binutils vim gdb openssh-client gcc g++ make libffi-dev openssl-dev libtool protobuf&& \
    echo 'set auto-load safe-path /' > /root/.gdbinit && \
    git config --global --add url."git@github.com:".insteadOf "https://github.com"
RUN go mod download
RUN make build

FROM alpine:3.18.3
WORKDIR /app

ARG APP_NAME
ARG ENV

ENV PATH="/usr/local/go/bin:${PATH}" \
    GOPRIVATE="github.com/Taskon-xyz" \
    GO111MODULE="on" \
    ENV=${ENV} \
    NAME=${APP_NAME}

RUN apk update && \
    apk add --no-cache curl netcat-openbsd bind-tools

COPY --from=builder /data/bin /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./kratos-layout", "-conf", "/data/conf"]
