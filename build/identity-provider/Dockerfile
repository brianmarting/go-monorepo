FROM golang:alpine3.18 AS builder

WORKDIR /build/

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o identity-provider ./cmd/identity-provider

FROM scratch

WORKDIR /app/

COPY --from=builder /build/identity-provider/identity-provider .

CMD ["/app/external-backend"]