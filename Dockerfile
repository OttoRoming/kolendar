FROM golang:1.26.4-alpine3.23

WORKDIR /app

COPY . .

RUN go mod download

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@1.31.1

RUN make

EXPOSE 8080
CMD ["./kolendar"]
