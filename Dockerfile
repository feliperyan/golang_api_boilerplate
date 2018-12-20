FROM golang:latest
RUN mkdir /app
WORKDIR /go/src/github.com/feliperyan/go_api_example_1
COPY . .
RUN go build -o main .
CMD ["./main"]