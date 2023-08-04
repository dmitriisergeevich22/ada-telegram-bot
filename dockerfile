FROM golang:1.20
WORKDIR /ada
COPY . .
RUN go build -o ada ./cmd/main.go
CMD ["./ada"]