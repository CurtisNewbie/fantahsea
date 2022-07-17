FROM golang:1.17-alpine 
LABEL author="yongjie.zhuang"
LABEL descrption="Fantahsea - A Gallery Service"

WORKDIR /usr/src/fantahsea

# for golang env
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# build executable
COPY . .
RUN go build -o ./main

# script (for io redirection stuff)
# COPY run.sh ./ 
RUN chmod +x run.sh

CMD ["sh", "run.sh"]
