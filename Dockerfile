FROM golang:1.17-alpine
COPY . /app/
WORKDIR /app/src
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build
CMD ["./diorama"]