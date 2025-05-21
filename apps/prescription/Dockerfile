FROM golang:1.23.8 as builder

ARG THIS_SERVICE=prescription

WORKDIR /app
COPY pkg/ ./pkg

WORKDIR /app/apps/$THIS_SERVICE
COPY apps/$THIS_SERVICE/go.mod apps/$THIS_SERVICE/go.sum ./
RUN go mod download

COPY apps/$THIS_SERVICE/ /app/apps/$THIS_SERVICE

WORKDIR /app/apps/$THIS_SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service ./cmd/main.go

RUN chmod +x /app/service


FROM gcr.io/distroless/static as runner

WORKDIR /app

COPY --from=builder /app/service /app/service

ENTRYPOINT ["/app/service"]
