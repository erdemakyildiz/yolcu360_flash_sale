FROM golang:1.23.1

COPY . /app

WORKDIR /app/

RUN go mod download
RUN go test -count=1 ./...
RUN go build -o main

CMD ["./main"]