FROM golang:1.18-alpine as build
LABEL author="Yongjie Zhuang"
LABEL descrption="Fantahsea - A Gallery Service"

RUN apk --no-cache add tzdata
WORKDIR /go/src/build/

# for golang env
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# dependencies
COPY go.mod .
COPY go.sum .

RUN go mod download

# build executable
COPY . .
RUN go build -o main

# ---------------------------------------------

FROM alpine:3.17

# for convert (change the source if necessary)
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
# RUN apk add file
# RUN apk add imagemagick

WORKDIR /usr/src/fantahsea
COPY --from=build /go/src/build/main ./main
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Shanghai

CMD ["./main", "app.name=fantahsea", "configFile=/usr/src/fantahsea/config/app-conf-prod.yml"]
