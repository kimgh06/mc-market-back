FROM golang:alpine AS builder

WORKDIR /build

COPY ./ ./

RUN go mod download
RUN go build -o maple .

WORKDIR /dist
RUN cp /build/maple .

FROM alpine AS runtime

COPY --from=builder /dist/maple .
COPY schema ./schema




