This program replaces docker images with images from the last successful Jenkins build.

Before first run:
- Create a `config.json` file using the `example.config.json` file.
- The files `config.json` and `main.go` must be located in same directory.

Run:
- In the terminal go to the directory with the `main.go` file.
- Execute: `go run main.go`


# TODO:
1. Make an example docker-compose.yaml
2. Make a better readme to explain how it works usaing real data