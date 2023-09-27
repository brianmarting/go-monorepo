# External backend

## Running the application

Run the `/cmd/main.go` file with the following:
- Environment variables: `ENVIRONMENT=dev;GRPC_PORT=8100;OTLP_ENDPOINT=127.0.0.1:4318;RABBITMQ_HOST=localhost;RABBITMQ_PASSWORD=guest;RABBITMQ_USERNAME=guest;REST_PORT=8000`
- Argument: `external-backend`

Can also be run in a docker container with the docker compose file in the root folder

## Functionality

Application will accept POST request on /mineral/deposit in order to deposit a mineral. It will then publish this on a queue.
Requests need to be authenticated with a bearer token that can be fetched from the IDP service, by default the endpoint will timeout if no response is received from the idp after 10 seconds and will throw a 401.
Should the authentication succeed but should the request take long to process, after another 10 seconds the application will terminate and throw a 502

Grpc can also be sent by looking at the contract.