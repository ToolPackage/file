FROM golang:1.14 as builder
LABEL stage=builder
ENV GOPROXY https://goproxy.cn
WORKDIR  /src
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /fse/bin/server cmd/server/main.go


FROM alpine:3.11 as prod
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates
WORKDIR /fse/bin/
COPY --from=builder /fse/bin/ .
EXPOSE 8000
CMD ["./server"]
