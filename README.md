# Golang Blog API

A RESTful blog API built with Go, Gin, and PostgreSQL. Designed with a clean architecture pattern separating concerns into handlers, repositories, models, and middleware.

## Tech Stack

- **Language:** Go 1.26+
- **Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL
- **Driver:** [pgx v5](https://github.com/jackc/pgx) (connection pool)
- **Authentication:** JWT (JSON Web Tokens) & bcrypt
- **Migration:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Hot Reload:** [Air](https://github.com/air-verse/air)

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                  # Application entrypoint & routes
├── internal/
│   ├── config/
│   │   └── config.go                # Environment configuration
│   ├── database/
│   │   └── postgres.go              # Database connection pool
│   ├── handlers/
│   │   ├── blog_handler.go          # Blog HTTP handlers
│   │   ├── category_handler.go      # Category HTTP handlers
│   │   └── user_handler.go          # User authentication handlers
│   ├── middleware/
│   │   └── auth_middleware.go       # JWT Authentication middleware
│   ├── models/
│   │   ├── blog.go                  # Blog model
│   │   ├── category.go              # Category model
│   │   └── user.go                  # User model
│   └── repository/
│       ├── blog_repository.go       # Blog database queries
│       ├── category_repository.go   # Category database queries
│       └── user_repository.go       # User database queries
├── migrations/                       # SQL migration files
├── scripts/
│   └── migrate.sh                   # Migration helper script
├── .air.toml                         # Air hot reload config
├── .env                              # Environment variables
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
DATABASE_URL=postgresql://postgres:password@localhost:5432/blog_api?sslmode=disable
PORT=8080
JWT_SECRET=your_super_secret_jwt_key
```

4. Create the database

```bash
psql -U postgres -c "CREATE DATABASE blog_api;"
```

5. Run migrations

```bash
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

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET`  | `/`      | No   | Server health check & DB status |

---

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/auth/register` | No | Register a new user |
| `POST` | `/auth/login`    | No | Authenticate user & receive JWT token |

#### `POST /auth/register`

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "johndoe", "email": "john@example.com", "password": "password123"}'
```

#### `POST /auth/login`

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "johndoe", "password": "password123"}'
```

---

### Categories

| Method   | Endpoint          | Auth | Description |
|----------|-------------------|------|-------------|
| `POST`   | `/categories`     | Yes  | Create a new category |
| `GET`    | `/categories`     | No   | Get all categories |
| `GET`    | `/categories/:id` | No   | Get category by ID |
| `PUT`    | `/categories/:id` | Yes  | Update category |
| `DELETE` | `/categories/:id` | Yes  | Delete category |

#### `POST /categories`

```bash
curl -X POST http://localhost:8080/categories \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Technology", "slug": "technology"}'
```

---

### Blog Posts

| Method   | Endpoint          | Auth | Description |
|----------|-------------------|------|-------------|
| `GET`    | `/blogs`          | No   | Get all public blog posts (with category names via `LEFT JOIN`) |
| `GET`    | `/blogs/:id`      | No   | Get blog post by ID |
| `POST`   | `/blogs`          | Yes  | Create a new blog post |
| `GET`    | `/blogs/user`     | Yes  | Get blog posts authored by logged-in user |
| `PUT`    | `/blogs/:id`      | Yes  | Update blog post (Owner only) |
| `DELETE` | `/blogs/:id`      | Yes  | Delete blog post (Owner only) |

#### `POST /blogs`

```bash
curl -X POST http://localhost:8080/blogs \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Golang REST API", "content": "Building APIs with Gin and PostgreSQL", "category_id": 1}'
```

```json
// 201 Created
{
  "message": "Blog created successfully",
  "blog": {
    "id": 1,
    "title": "Golang REST API",
    "content": "Building APIs with Gin and PostgreSQL",
    "created_at": "2026-07-20T21:50:00Z",
    "updated_at": "2026-07-20T21:50:00Z",
    "user_id": "77051b22-ea10-4a78-a7de-f880e5b6d60a",
    "category_id": 1
  }
}
```

#### `GET /blogs`

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
      "title": "Golang REST API",
      "content": "Building APIs with Gin and PostgreSQL",
      "created_at": "2026-07-20T21:50:00Z",
      "updated_at": "2026-07-20T21:50:00Z",
      "user_id": "77051b22-ea10-4a78-a7de-f880e5b6d60a",
      "category_id": 1,
      "category_name": "Technology"
    }
  ]
}
```

---

### Error Handling

| Status Code | Description |
|-------------|-------------|
| `400 Bad Request` | Invalid JSON input or validation failure |
| `401 Unauthorized` | Missing, invalid, or expired JWT token |
| `403 Forbidden` | User is not the owner of the resource |
| `404 Not Found` | Requested resource does not exist |
| `409 Conflict` | Duplicate resource (e.g. username/email/category name already exists) |
| `500 Internal Server Error` | Unexpected server error (masked for security) |

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
./scripts/migrate.sh fresh           # Drop all tables & re-apply all
```

---

## Roadmap

- [x] Blog CRUD (Create, Read, Update, Delete)
- [x] User Authentication & Authorization (JWT & bcrypt)
- [x] Blog Ownership Authorization (Owner-only Edit/Delete)
- [x] Categories CRUD & Foreign Key Relation (`category_id` & `LEFT JOIN`)
- [x] Masked Error Handling for Security
- [ ] Likes
- [ ] Comments with reply threads

---

## License

This project is for educational purposes.
