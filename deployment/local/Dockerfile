FROM golang:1.24.6-alpine AS builder
WORKDIR /app

COPY . .

ARG BUILD_TARGET
RUN go build -o /app/bin/app $BUILD_TARGET

FROM scratch

COPY --from=builder /app/bin/app /

ENTRYPOINT ["/app"]