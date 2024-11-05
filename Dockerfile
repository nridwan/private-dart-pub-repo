FROM golang:1.23.2-alpine3.20 AS builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR .
COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /out/pubserver .

FROM alpine:3.20.3
COPY --from=builder /out/pubserver /

ENTRYPOINT ["/pubserver"]
CMD ["fx"]