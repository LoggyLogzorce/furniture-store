FROM golang:1.21.6

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . ./

RUN go build

EXPOSE 8080

