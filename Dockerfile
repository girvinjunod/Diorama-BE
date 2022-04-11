FROM golang:1.18-alpine as builder
COPY .env ./
WORKDIR /build
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ ./
RUN CGO_ENABLED=0 go build

FROM alpine
WORKDIR /
COPY .env ./
COPY --from=builder /build/diorama ./bin
WORKDIR /bin
CMD ["./diorama"]