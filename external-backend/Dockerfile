FROM golang:alpine3.18

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY .. ./

RUN go build -v -o go-backend

CMD ["/app/go-backend"]