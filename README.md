# go-monorepo

## Structure

The application is built up of the following services:
- external-backend: Backend responsible for interacting with the external parties via REST or GRPC. This then sends the data
to a broker
- internal-backend: Will consume the data from the broker and save it to the DB

## Running the application

The application can be fully run with docker. Just start the docker compose files. The applications will be built when the image is not present.

## Tracing

Tracing is being done by using the `otlphttptrace` package, meaning that the application will send traces via http to the set location. \
This being done in the otlp format, and is being sent to an otel collector. In the collector it is being logged and being forwarded in the
same otlp format to the given ingester, in this case being jaeger. \
This can be viewed when running the `docker-compose.yml` and going to `http://localhost:16686/` (jaeger UI).


