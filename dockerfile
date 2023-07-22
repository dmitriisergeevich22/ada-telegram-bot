FROM golang:1.20
WORKDIR /ada
COPY . .
RUN go build -o adaApp ./cmd/main.go
CMD ["./adaApp"]