FROM golang as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
# 设置工作区
WORKDIR /usr/local/go
ADD go.mod .
ADD go.sum .
RUN go mod download
# 把全部文件添加到/usr/local/go目录
ADD . .
# 编译：把cmd/main.go编译成可执行的二进制文件，命名为app
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o  aqi-server aqi-server/cmd

FROM centos
WORKDIR /usr/local/go

COPY --from=build /usr/local/go/aqi-server .
COPY assets .

CMD ["./aqi-server","--config","conf.yml"]
