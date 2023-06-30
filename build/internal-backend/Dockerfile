FROM golang:alpine3.18 AS builder

WORKDIR /build/

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o internal-backend ./cmd/internal-backend

FROM scratch

WORKDIR /app/

COPY --from=builder /build/internal-backend/internal-backend .

CMD ["/app/internal-backend"]