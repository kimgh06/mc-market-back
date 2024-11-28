FROM golang:alpine

WORKDIR /build

COPY ./ ./

RUN go mod download
RUN go build -o maple .

WORKDIR /dist
RUN cp /build/maple .

WORKDIR /

COPY /dist/maple .
COPY schema ./schema
