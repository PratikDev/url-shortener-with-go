# ShortURL

A simple URL shortener application built with Go.

## Endpoints

- `GET /` - Home page
- `GET /health` - Health check
- `GET /{id}` - Get the original URL by the short URL
- `POST /login` - Login to the application
- `POST /register` - Register a new user
- `POST /new` - Create a new short URL
- `GET /all` - List all the short URLs created by the user

## Run the application locally

The below instructions assumes that you have Go installed on your machine.

1. Clone the repository

```bash
git clone https://github.com/PratikDev/url-shortener-with-go.git
```

2. Change the directory

```bash
cd url-shortener-with-go
```

3. Install the dependencies

```bash
go mod tidy
```

4. Run the application

```bash
go run main.go
```
