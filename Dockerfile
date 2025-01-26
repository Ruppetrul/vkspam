FROM golang:1.23.2

WORKDIR /

COPY go.mod main.go ./

RUN go mod download

COPY . .

RUN go build -o main

EXPOSE 80

CMD ["./main"]