FROM golang as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
# 设置工作区
WORKDIR /usr/local/go
ADD . .
RUN go mod download
# 编译：把cmd/main.go编译成可执行的二进制文件，命名为app
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o  storm-aqi-server github.com/csnight/storm-aqi-server/cmd/server
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o  doc-gen github.com/csnight/storm-aqi-server/cmd/tools
RUN ./doc-gen

FROM centos
WORKDIR /usr/local/go

COPY --from=build /usr/local/go/storm-aqi-server .

CMD ["./storm-aqi-server","--config","conf.yml"]
