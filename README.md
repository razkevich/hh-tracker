# personal-data.svc

## Introduction

This repository contains the Personal Data Service which is the main service responsible for delivering the Personal Data management functionality.
See [the architecture document](https://elasticpath.atlassian.net/wiki/spaces/DOCPD/pages/1556611106/Personal+Data+Management+in+Elastic+Path+Commerce+Cloud)
for more details

## Table of Contents

* Information about [Development is here](#development)
* Information about [Deployment is here](#deployment)
* Information about [Component Tests is here](component-tests/README.md)

## Development

### Prerequisites

The following table lists _hard_ dependencies you will need to use this project, in almost all cases using earlier versions is **known not to work**.

| Name                                                                                       | Version      | Notes                                        |
|--------------------------------------------------------------------------------------------|--------------|----------------------------------------------|
| [Go](https://golang.org/doc/)                                                              | 1.13.6       | Required to build and spin up services       |
| [docker](https://www.docker.com/products/docker-desktop)                                   | 18.02.0+     | Required to build and spin up services       |
| [docker-compose](https://docs.docker.com/compose/install/)                                 | 1.20.0+      | Required to build and spin up services       |
| [Node](https://nodejs.org/en/download/)                                                    | 14.16        | Required if running locally or running tests |
| [yarn](https://yarnpkg.com/)                                                               | 1.22.5       | Required if running locally or running tests |

### Useful Commands

| Command         | Description                                                                                   |
| ----------------| ----------------------------------------------------------------------------------------------|
| component-test  | Runs the component tests                                                                      |
| destroy         | Destroys the containers                                                                       |
| lint            | Performs lint checks on Go and TypeScript code                                                |
| start           | Starts the service and all necessary dependencies in the background                           |
| stop            | Stops running containers for either `make start` or `make watch`                              |
| test-coverage   | Runs `go test -v ./...` and displays test coverage in HTML report                             |
| watch           | Starts the service and all necessary dependencies in the foreground                           |
| prepare         | Prepares the code for pushing to git. (runs lint fixes)                                       |
| check           | Runs prepare followed by test then component-test                                             |

### Service Locations

| Service                                                                                                                                             | Location                                             | Username          | Password    |
| ----------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------- |-------------------| ------------|
| API Endpoint for PDS                                                                                                                                | http://localhost:8057                                | n/a               | n/a         |
| [Remote Debugging](https://github.com/go-delve/delve) for PDS                                                                                       | tcp://localhost:2545                                 | n/a               | n/a         |
| MongoDB                                                                                                                                             | tcp://localhost:57017                                |                   |             |
| Rabbit MQ (AMPQ)                                                                                                                                    | tcp://localhost:5675                                 | ???               | ???         |
| Rabbit MQ (User Interface)                                                                                                                          | http://localhost:15675                               | guest             | guest       |
| Wiremock (Used for testing)                                                                                                                         | http://localhost:8581                                | n/a               | n/a         |

### Debugging

The service uses Delve to enable remote debugging. Refer to [The Delve Repository](https://github.com/go-delve/delve) to learn about how to configure
your favorite IDE to connect to the application. This
[internal resource](https://elasticpath.atlassian.net/wiki/spaces/DOCPD/pages/422575207/Debugging+Multi-Tenant+Go+services#DebuggingMulti-TenantGoservices-ConfiguringIntelliJ)
might also be useful.

### Linting

This project uses [golangci-lint](https://github.com/golangci/golangci-lint) for linting any Go source code.

*The pipeline will fail if errors are not corrected*.

To install locally, please refer to
[https://golangci-lint.run/usage/install/#local-installation](https://golangci-lint.run/usage/install/#local-installation)

### Dependency Injection

This project uses the dependency injection framework [Google Wire](https://github.com/google/wire).

Our injectors for the app are specified in [cmd/app/inject_main.go](cmd/app/inject_main.go).

If additional providers are added, please make sure to regenerate the `wire_gen.go` file with:

```shell script
go generate ./... 
```

If the `wire_gen.go` file does not exist, you will need to run wire directly:

```shell script
wire gitlab.elasticpath.com/commerce-cloud/personal-data.svc/cmd/app
```

### Project Layout

This project roughly followed the layout of Go projects as described at
[https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout).

| Directory         | Description                                                                                    |
|-------------------|------------------------------------------------------------------------------------------------|
| `bin/`             | Any scripts for running or building the project are located in here                            |
| `cmd/`             | This Go package is where we use `main` for the executables of the project                      |
| `component-tests/` | Component tests for the service are located in here.                                           |
| `external/`        | Any Go package that are candidates as dependencies for projects are in here                    |
| `internal/`        | Application specific Go packages, e.g., they cannot be shared and are specific to this service |
| `local/`           | Any files relating to running the project locally should be in here                            |

## Layers and Folder Structure:

There are 3 main layers in this repo, including `Controller`, `Service`, and `Repository`. The only and only way for these layers to interact with
each other should be through their interfaces. And, the lower layer does not have any knowledge about the upper layers.

The `Entity` and `DTO` are shared between these layers. The entity represents the model in the database.  
DTO (Data Transfer Object) is used exclusively for interface definition in layers.

`helper` only does the final touches to adjust the data contract for the endpoint. The sole usage of the helper is the insider controller. It cannot
be used in Service or Repository. 

### Testing Against Built CI Production Containers

In our Continuous Integration pipeline we test _both_ production and development versions of Personal data, we test the production container as that
is the one that will be in production, but we also test the development container to prevent and detect breakages that impact development. Developers,
when running locally, will likely be using the development containers.

To test against the production containers (for instance to debug a failure in CI, and not locally do the following)

1. Find the commit hash of the built image
    * One way is to find your [Gitlab Pipeline](https://gitlab.elasticpath.com/commerce-cloud/personal-data.svc/-/pipelines/)
    * Open the `build-docker-image` job.
    * The image name will be in the output for
      example: `quay.io/elasticpath/personal-data.svc@sha256:ce32481c76c3df7a4119abce77e85e142154fc32622d011b2e5ab0ed225d3c34`
2. Shutdown all containers running locally (and also destroy their local data with `docker-compose down -v`)
3. Set the `PERSONAL_DATA_IMAGE` environment variable (
   e.g., `export PERSONAL_DATA_IMAGE=quay.io/elasticpath/personal-data.svc@sha256:ce32481c76c3df7a4119abce77e85e142154fc32622d011b2e5ab0ed225d3c34`)
4. Pull the image using `docker-compose pull`
5. Bring up the images using the following command `docker-compose up --no-build`
    * Note: `--no-build` is important as otherwise docker will build the local image.

At this point the production variants of the containers should be running, and you can run the component tests as usual.

### Locally building the Production Containers

Occasionally you will want to test against the production containers locally, the fastest way to do so is to run the following command:

```bash
 PERSONAL_DATA_DOCKERFILE=Dockerfile docker-compose up --build
```

This will change the Dockerfile being used to the production version, but allow everything else locally to work.

## Deployment

### Dependencies

#### Infrastructure Dependencies

Personal data depends on both RabbitMQ and MongoDB.

### Environment Variables

The following environment variables are [defined in Personal data](internal/config/env.go), and can be used to influence behaviour.

| Name                                                    | Default                                                                                                                               | Description                                                                                                                                                               |
| ------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `MONGO_DSN`                                             |                                                                                                                                       | MongoDB DSN                                                                                                                                                               |
| `MONGO_DATABASE_NAME`                                   | personal-data                                                                                                                         | The name of the MongoDB database to use                                                                                                                                   |
| `MONGO_TIMEOUT`                                         | 5000                                                                                                                                  | MongoDB connection timeout                                                                                                                                                |
| `ENFORCE_LIMITS`                                        | true                                                                                                                                  | The variable that enforces limitation for offsets of records when listing entities                                                                                        |
| `DEFAULT_PAGE_LIMIT`                                    | 20                                                                                                                                    | The variable that controls the pagination page limit                                                                                                                      |
| `LOGGING_LEVEL`                                         | `debug`                                                                                                                               | Log level                                                                                                                                                                 |
| `SVC_BASE_URL`                                          | This value depends on the environment, as is set consistently across services in an environment, ask the operator of the environment. | The variable that makes the application aware of the base url it's exposed at                                                                                             |
| `SVC_PORT`                                              | 8000                                                                                                                                  | The variable that makes the application aware of the port it's exposed at                                                                                                 |
| `RABBIT_HOSTS`                                          |                                                                                                                                       | The list of RabbitMQ hosts                                                                                                                                                |
| `RABBIT_CONSUME_QUEUE`                                  | personal-data-consumer-queue                                                                                                          | The RabbitMQ consume queue                                                                                                                                                |
| `RUN_MIGRATIONS`                                        |                                                                                                                                       | Set to `true` to run migrations before container start up, `only` to run migrations and exit, and `false` to skip migrations.                                             |

Additional environment variables for controlling component tests [are defined here](component-tests/shared/config.ts)
and can be used to influence behaviour of the component tests.

### Multiple Services With One Rabbitmq Instance

Occasionally you want to test multiple services can send/receive messages to another locally. You can change default RABBIT_HOSTS env variable when
running the service. For example, you can run following command:

```bash
RUN RABBIT_HOSTS=amqp://guest:guest@host.docker.internal:35672 make watch
```