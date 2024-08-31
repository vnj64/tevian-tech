FROM golang:1.22.1 as build_base
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tevian

EXPOSE 8080

CMD ["go", "run", "main.go"]