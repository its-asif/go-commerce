# GO Commerce API

This is a RESTful API for an E-commerce platform built with Go. It supports user authentication, product management, cart operations, order processing, and category management.

## Features

- User registration and login with JWT authentication
- Admin and user roles
- Product CRUD (Create, Read, Update, Delete)
- Cart management (add, remove, view items)
- Order checkout and history
- Category management (admin only)
- PostgreSQL database integration

## Project Structure

- `handlers/` - HTTP handlers for API endpoints
- `models/` - Data models
- `middleware/` - Authentication and authorization middleware
- `db/` - Database connection logic
- `config/` - Environment variable loading
- `routes/` - API route definitions
- `migrations/` - SQL migration files
- `utils/` - Utility functions (JWT, hashing, etc.)

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

## Setup Instructions

### 1. Clone the Repository

```sh
git clone https://github.com/your-username/go-commerce.git
cd go-commerce
```

### 2. Configure Environment Variables

Copy `.env.example` to `.env` and update the values if needed:

```sh
cp .env.example .env
```

- `PORT`: The port the API will run on (default: 8000)
- `JWT_SECRET`: Secret key for JWT token signing

### 3. Build and Run with Docker Compose

This will start the Go API server, PostgreSQL database, and run migrations automatically.

```sh
docker-compose up --build
```

- The API will be available at `http://localhost:8000`
- PostgreSQL will be available at `localhost:5432` (internal to Docker network as `go_db`)

### 4. API Endpoints

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and receive JWT token
- `GET /api/products` - List all products
- `POST /api/cart` - Add item to cart (authenticated)
- `POST /api/orders/checkout` - Checkout cart (authenticated)
- ...and more (see [routes/routes.go](routes/routes.go))

### 5. Stopping the Services

Press `Ctrl+C` in the terminal or run:

```sh
docker-compose down
```

## Database Migrations

Migrations are automatically applied on startup using the `migrate` service in `docker-compose.yml`. Migration files are located in the [migrations/](migrations/) directory.


## API Testing

You can test the API using the public Postman collection:  
[GO Commerce API Postman Collection](https://www.postman.com/nemoh3618/public-workspace/collection/u2bc35o/ego-commerce?action=share&creator=31642937)

---

**Note:** For development without Docker, ensure you have Go and PostgreSQL installed, set up your `.env` file, run migrations manually, and start the server with:

```sh
go run main.go
```