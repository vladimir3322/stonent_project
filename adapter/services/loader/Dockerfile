FROM golang:1.14
ENV GOPROXY=direct
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o main ./main.go

CMD ["/app/main"]
