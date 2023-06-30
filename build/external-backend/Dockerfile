FROM golang:alpine3.18 AS builder

WORKDIR /build/

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o external-backend ./cmd/external-backend

FROM scratch

WORKDIR /app/

COPY --from=builder /build/external-backend/external-backend .

CMD ["/app/external-backend"]