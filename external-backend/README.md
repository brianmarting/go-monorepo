# External backend

## Running the application

Run the `/cmd/main.go` file with the following:
- Environment variables: `ENVIRONMENT=dev;GRPC_PORT=8100;OTLP_ENDPOINT=127.0.0.1:4318;RABBITMQ_HOST=localhost;RABBITMQ_PASSWORD=guest;RABBITMQ_USERNAME=guest;REST_PORT=8000`
- Argument: `external-backend`