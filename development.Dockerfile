FROM golang:alpine AS builder

WORKDIR /build

COPY ./ ./

RUN go mod download
RUN go build -o mapi .

WORKDIR /dist
RUN cp /build/mapi .

FROM alpine AS runtime
COPY --from=builder /dist/maple .
COPY schema ./schema