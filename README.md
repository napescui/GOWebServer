# GOWebServer

GOWebServer is a web server for Growtopia Private Server. It is written in Go and is designed to be fast and efficient.

## To-Do

- [x] Basic HTTP server
- [x] Handle `/growtopia/server_data.php` requests
- [x] Implementing rate limiting requests
- [ ] Implementing cache server for Growtopia Client
- [ ] Handling missing cache files
- [ ] Geo Location checker to block certain countries

## Build

The following are required to build and run GOWebServer:

- [Golang](https://golang.org/dl/) (1.16+) - The Go Programming Language
- and little bit of brain cells (optional)

Building the server is simple, just run the following command:

- 1. Clone the repository:

```bash
git clone https://github.com/yoruakio/GOWebServer.git
```

- 2. Build the server:

```bash
go build

# or running the go file directly
go run main.go
```

- 3. Run the server:

```bash
./GOWebServer
```

## Configuration

The server can be configured using the `config.json` file. The following are the default configuration:

```json
{
    "host": "127.0.0.1",
    "port": "17091",
    "loginUrl": "gtlogin-backend.vercel.app",
    "isLogging": true,
    "rateLimit": 150,
    "rateLimitDuration": 5
}
```

## Contributing

Contributions are welcome! If you would like to contribute to the project, please fork the repository and submit a pull request.

## Aknowledgements

- [GTPSWebServer](https://github.com/yoruakio/GTPSWebServer) - The original GOWebServer was inspired by this project.
- [Golang](https://golang.org/) - The Go Programming Language

## License

GOWebServer is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
