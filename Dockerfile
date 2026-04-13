ARG DOCKER_REPO_BASE_PATH=northamerica-northeast1-docker.pkg.dev/prj-c-mime-build-3ye3/docker-images
ARG DOCKER_IMAGE_TAG=20260220-183339-3f86bad

# <== Build Image ==>
FROM ${DOCKER_REPO_BASE_PATH}/base-build:${DOCKER_IMAGE_TAG} AS build

# Download dependencies as a separate step to take advantage of Docker's caching.
COPY go.mod go.sum ./

RUN --mount=type=ssh \
    --mount=type=cache,target=/go/pkg/mod/ \
    go mod download;

# Copy the source files.
# Using Dockers wildcard mechanism to include directories if they exist.
COPY /cm[d] ./cmd
COPY /interna[l] ./internal
COPY /pk[g] ./pkg

# Build the application.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOCACHE=/go/pkg/mod/ go build -o /usr/local/bin/ ./...;


# <== Application Image ==>
FROM ${DOCKER_REPO_BASE_PATH}/base-app:${DOCKER_IMAGE_TAG} AS app

EXPOSE 8080

# Copy migration files (if present)
COPY migration[s] /migrations
COPY atla[s].hcl .

# Copy the executable from the "build" stage.
COPY --from=build /usr/local/bin/* /usr/local/bin/
