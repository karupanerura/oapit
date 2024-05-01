FROM golang:1.22 as builder

WORKDIR /app

COPY ["./go.mod", "./go.sum", "./"]

RUN go mod download

COPY ./ ./

RUN go build -o oapit .

FROM gcr.io/distroless/base-debian12

COPY --from=builder /app/oapit /
ENTRYPOINT ["/oapit"]