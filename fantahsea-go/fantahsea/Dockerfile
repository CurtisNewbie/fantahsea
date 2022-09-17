FROM golang:1.18-alpine 
LABEL author="yongjie.zhuang"
LABEL descrption="Fantahsea - A Gallery Service"

WORKDIR /usr/src/fantahsea

# for golang env
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# for convert (change the cdn if necessary)
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add file
RUN apk add imagemagick

# dependencies
COPY go.mod .
COPY go.sum .
# RUN go mod download

RUN go mod tidy 

# build executable
COPY . .

RUN go build -o ./main

# script (for io redirection stuff)
# COPY run.sh ./ 
RUN chmod +x run.sh

CMD ["sh", "run.sh"]
