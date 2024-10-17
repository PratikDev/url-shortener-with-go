# ShortURL

A simple URL shortener application built with Go.

## Endpoints

- `GET /` - Home page
- `GET /health` - Health check
- `GET /{id}` - Get the original URL by the short URL
- `GET /all` - List all the short URLs created by the user
- `POST /login` - Login to the application
- `POST /register` - Register a new user
- `POST /new` - Create a new short URL

## Request Body

- `POST /login`

```json
{
	"username": "username",
	"password": "password"
}
```

- `POST /register`

```json
{
	"username": "username",
	"password": "password"
}
```

- `POST /new`

```json
{
	"url": "https://www.google.com"
}
```

## Technologies Used

- Go
- MongoDB (Atlas)
- JWT

## Environment Variables

- `MONGODB_CONNECTION_STRING` - MongoDB connection string
- `JWT_SECRET` - Secret key for JWT token

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

3. Create a `.env` file and add the environment variables

```bash
touch .env
```

Add the following environment variables to the `.env` file

```env
MONGODB_CONNECTION_STRING=<your_mongodb_connection_string>
JWT_SECRET=<your_jwt_secret>
```

4. Install the dependencies

```bash
go mod tidy
```

5. Run the application

```bash
go run main.go
```
