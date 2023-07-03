FROM golang:alpine3.18 AS builder

ARG app

WORKDIR /build/

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o application ./cmd/$app

FROM scratch

ARG app

WORKDIR /app/

COPY --from=builder /build/application .

ENTRYPOINT ["/app/application"]