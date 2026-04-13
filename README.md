# irp-app-from-template

IRP scanner API

# Development

To run the formatters and linters automatically on commit, you can use <a href="https://pre-commit.com/">pre-commit</a>

```
pre-commit install
```

## Make commands

- `make format`: Run the code formatter
- `make lint`: Run linters
- `make test`: Run all tests
- `make test <package>`: Run tests for a specific package
- `make mod`: Run `go mod tidy` to ensure your go.sum is up to date
- `make update`: Update all project dependencies
- `make upgrade`: Get new changes from the template and update all project dependencies
- `make clean`: Clean project with `go clean`
- `make mock.gen`: generate mock packages for clients
- `make cov.render`: Show coverage stats in your browser
- `make docker.build`: Build the docker image for this project
- `make docker.run`: Build and start a docker container for this project
- `make docker.clean`: Remove docker containers from the current project
- `make docker.clean all=1`: Remove docker containers from all MiMe projects

## Generating mocks

We use mockery to generate mocks under <a>internal/testing/mocks</a>.
If you add or modify a client you should run `make mock.gen` to ensure those mocks are up to date.
