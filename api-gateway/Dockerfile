FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]