FROM golang:1.24.2-alpine AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/local/bin/app ./main.go

FROM alpine:3.21.3
COPY --from=builder /usr/local/bin/app /usr/local/bin/ical-converter
CMD ["ical-converter"]
