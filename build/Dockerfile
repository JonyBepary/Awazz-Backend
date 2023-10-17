FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

ENV PATH=$PATH:/usr/local/bin

RUN if ! command -v protoc &> /dev/null; then \
        curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v25.0-rc1/protoc-25.0-rc-1-linux-x86_64.zip && \
        unzip protoc-25.0-rc-1-linux-x86_64.zip && \
        sudo mv bin/protoc /usr/local/bin && \
        sudo mv include/google /usr/local/include; \
    fi

RUN go mod tidy && \
    go mod vendor

CMD ["go", "build", "-v", "./..."]
