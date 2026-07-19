# Golang Blog API

A RESTful blog API built with Go, Gin, and PostgreSQL. Designed with a clean architecture pattern separating concerns into handlers, repositories, and models.

## Tech Stack

- **Language:** Go 1.26+
- **Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL
- **Driver:** [pgx v5](https://github.com/jackc/pgx) (connection pool)
- **Migration:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Hot Reload:** [Air](https://github.com/air-verse/air)

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entrypoint
├── internal/
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── database/
│   │   └── postgres.go          # Database connection pool
│   ├── handlers/
│   │   └── blog_handler.go      # HTTP request handlers
│   ├── middleware/               # Middleware (planned)
│   ├── models/
│   │   └── blog.go              # Data models
│   └── repository/
│       └── blog_repository.go   # Database queries
├── migrations/                   # SQL migration files
├── scripts/
│   └── migrate.sh               # Migration helper script
├── .air.toml                     # Air hot reload config
├── .env                          # Environment variables
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.26 or higher
- PostgreSQL 14 or higher
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [Air](https://github.com/air-verse/air) (optional, for hot reload)

### Installation

1. Clone the repository

```bash
git clone git@github.com:Chilhan23/golang-blog-api.git
cd golang-blog-api
```

2. Install dependencies

```bash
go mod download
```

3. Set up environment variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgresql://<user>:<password>@localhost:5432/<dbname>?sslmode=disable
PORT=8080
```

4. Create the database

```bash
psql -U postgres -c "CREATE DATABASE blog_api;"
```

5. Run migrations

```bash
# Install golang-migrate (if not installed)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply all migrations
./scripts/migrate.sh up
```

6. Start the server

```bash
# With Air (hot reload)
air

# Without Air
go run ./cmd/api
```

The server will start at `http://localhost:8080`.

---

## API Documentation

### Health Check

```
GET /
```

**Response** `200 OK`

```json
{
  "message": "Go Gin API is running",
  "status": "success",
  "database": "connected"
}
```

---

### Blog Posts

| Method   | Endpoint     | Description          |
|----------|--------------|----------------------|
| `POST`   | `/blogs`     | Create a blog post   |
| `GET`    | `/blogs`     | Get all blog posts   |
| `GET`    | `/blogs/:id` | Get a blog post by ID|
| `PUT`    | `/blogs/:id` | Update a blog post   |
| `DELETE` | `/blogs/:id` | Delete a blog post   |

#### `POST /blogs` -- Create a Blog Post

**Request Body**

| Field     | Type   | Required | Description        |
|-----------|--------|----------|--------------------|
| `title`   | string | Yes      | Title of the post  |
| `content` | string | Yes      | Content of the post|

```bash
curl -X POST http://localhost:8080/blogs \
  -H "Content-Type: application/json" \
  -d '{"title": "Hello World", "content": "This is my first post."}'
```

```json
// 201 Created
{
  "message": "Blog created successfully",
  "blog": {
    "id": 1,
    "title": "Hello World",
    "content": "This is my first post.",
    "created_at": "2026-07-19T12:00:00Z",
    "updated_at": "2026-07-19T12:00:00Z"
  }
}
```

#### `GET /blogs` -- Get All Blog Posts

```bash
curl http://localhost:8080/blogs
```

```json
// 200 OK
{
  "message": "Blogs retrieved successfully",
  "blogs": [
    {
      "id": 1,
      "title": "Hello World",
      "content": "This is my first post.",
      "created_at": "2026-07-19T12:00:00Z",
      "updated_at": "2026-07-19T12:00:00Z"
    }
  ]
}
```

#### `GET /blogs/:id` -- Get a Blog Post by ID

```bash
curl http://localhost:8080/blogs/1
```

```json
// 200 OK
{
  "message": "Blog retrieved successfully",
  "blog": {
    "id": 1,
    "title": "Hello World",
    "content": "This is my first post.",
    "created_at": "2026-07-19T12:00:00Z",
    "updated_at": "2026-07-19T12:00:00Z"
  }
}
```

```json
// 404 Not Found
{
  "error": "Blog Not Found"
}
```

#### `PUT /blogs/:id` -- Update a Blog Post

Supports partial updates. Only the fields included in the request body will be updated. At least one field must be provided.

**Request Body**

| Field     | Type   | Required | Description          |
|-----------|--------|----------|----------------------|
| `title`   | string | No       | New title of the post|
| `content` | string | No       | New content          |

```bash
curl -X PUT http://localhost:8080/blogs/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Title"}'
```

```json
// 200 OK
{
  "message": "Blog updated successfully",
  "blog": {
    "id": 1,
    "title": "Updated Title",
    "content": "This is my first post.",
    "created_at": "2026-07-19T12:00:00Z",
    "updated_at": "2026-07-19T12:30:00Z"
  }
}
```

#### `DELETE /blogs/:id` -- Delete a Blog Post

```bash
curl -X DELETE http://localhost:8080/blogs/1
```

```json
// 200 OK
{
  "message": "Blog deleted successfully",
  "blog": {
    "id": 1,
    "title": "Hello World",
    "content": "This is my first post.",
    "created_at": "2026-07-19T12:00:00Z",
    "updated_at": "2026-07-19T12:00:00Z"
  }
}
```

---

### Error Responses

All error responses follow this format:

```json
{
  "error": "error description"
}
```

| Status Code | Description                          |
|-------------|--------------------------------------|
| `400`       | Bad request or invalid input         |
| `404`       | Resource not found                   |
| `500`       | Internal server error                |

---

## Database Migration

A helper script is provided to manage migrations using `golang-migrate`.

```bash
./scripts/migrate.sh up              # Apply all pending migrations
./scripts/migrate.sh up <N>          # Apply next N migrations
./scripts/migrate.sh down            # Rollback last migration (requires confirmation)
./scripts/migrate.sh down all        # Rollback all migrations (requires confirmation)
./scripts/migrate.sh create <name>   # Create new migration files
./scripts/migrate.sh redo            # Rollback last and re-apply
./scripts/migrate.sh status          # Show current version
./scripts/migrate.sh force <V>       # Force set version (fix dirty state)
```

---

## Roadmap

- [x] Blog CRUD (Create, Read, Update, Delete)
- [ ] User authentication (JWT)
- [ ] Categories (one-to-many with posts)
- [ ] Likes
- [ ] Comments

---

## License

This project is for educational purposes.
