# GO Commerce API

This is a RESTful API for an E-commerce platform built with Go. It supports user authentication, product management, cart operations, order processing, and category management.

## Features

- User registration and login with JWT authentication
- Admin and user roles
- Product CRUD (Create, Read, Update, Delete)
- Cart management (add, remove, view items)
- Order checkout and history
- Category management (admin only)
- PostgresSQL database integration
- Redis caching for improved performance
- Caching for products, categories, user data, cart items, and orders
- Swagger API Documentation

## Project Structure

- `handlers/` - HTTP handlers for API endpoints
- `models/` - Data models
- `middleware/` - Authentication and authorization middleware
- `db/` - Database connection logic (PostgreSQL & Redis)
- `config/` - Environment variable loading
- `routes/` - API route definitions
- `migrations/` - SQL migration files
- `utils/` - Utility functions (JWT, hashing, redis_cache)

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

## Setup Instructions

### 1. Clone the Repository

```sh
git clone https://github.com/its-asif/go-commerce.git
cd go-commerce
```

### 2. Configure Environment Variables

Copy `.env.example` to `.env` and update the values if needed:

```sh
cp .env.example .env
```

Environment Variables:
- `PORT`: The port the API will run on (default: 8000)
- `JWT_SECRET`: Secret key for JWT token signing
- `DB_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection URL (default: redis://localhost:6379)
- `REDIS_PASSWORD`: Redis password (if required)
- `REDIS_DB`: Redis database number (default: 0)

### 3. Build and Run with Docker Compose

This will start the Go API server, PostgreSQL database, Redis cache, and run migrations automatically.

```sh
docker-compose up --build
```

Services started:
- **API Server**: Available at `http://localhost:8000`
- **PostgreSQL**: Available at `localhost:5432` (internal to Docker network as `go_db`)
- **Redis**: Available at `localhost:6379` (internal to Docker network as `redis`)

### 4. Swagger API Documentation

The API includes Swagger documentation for all endpoints. After starting the server, you can access the Swagger UI at:

```
http://localhost:8000/swagger/index.html
```

Swagger provides an interactive interface to test and explore the API.

### 5. Redis Caching

The application implements Redis caching for:
- **Products**: Individual products cached for 5 minutes, all products for 10 minutes
- **Categories**: All categories cached for 30 minutes
- **User Data**: User login data cached for 15 minutes
- **Cart Items**: User cart items cached for 5 minutes
- **Orders**: User orders cached for 10 minutes

Cache keys are automatically invalidated when data is modified (create, update, delete operations).

### 6. API Endpoints

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and receive JWT token
- `GET /api/products` - List all products (cached)
- `GET /api/products/{id}` - Get single product (cached)
- `POST /api/products` - Create product (admin only, invalidates cache)
- `PUT /api/products/{id}` - Update product (admin only, invalidates cache)
- `DELETE /api/products/{id}` - Delete product (admin only, invalidates cache)
- `GET /api/categories` - List all categories (cached)
- `POST /api/categories` - Create category (admin only, invalidates cache)
- `GET /api/cart` - Get cart items (cached)
- `POST /api/cart` - Add item to cart (authenticated, invalidates cache)
- `DELETE /api/cart/{product_id}` - Remove item from cart (authenticated, invalidates cache)
- `GET /api/orders` - Get user orders (cached)
- `POST /api/orders/checkout` - Checkout cart (authenticated, invalidates cache)

### 7. Stopping the Services

Press `Ctrl+C` in the terminal or run:

```sh
docker-compose down
```

## Database Migrations

Migrations are automatically applied on startup using the `migrate` service in `docker-compose.yml`. Migration files are located in the [migrations/](migrations/) directory.

## Performance Optimization

The application uses Redis caching to improve performance:
- Frequently accessed data is cached to reduce database queries
- Cache invalidation ensures data consistency
- Different cache TTL (Time To Live) values based on data update frequency

## API Testing

You can test the API using the public Postman collection:  
[GO Commerce API Postman Collection](https://www.postman.com/nemoh3618/public-workspace/collection/u2bc35o/ego-commerce?action=share&creator=31642937)

---

**Note:** For development without Docker, ensure you have Go, PostgreSQL, and Redis installed, set up your `.env` file, run migrations manually, and start the server with:

```sh
go run main.go
```

Make sure Redis is running on your local machine:
```sh
redis-server
```