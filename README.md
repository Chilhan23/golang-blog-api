# Golang Blog REST API

A high-performance, production-ready RESTful Blog API built with Go, Gin Framework, and PostgreSQL. Designed following the Standard Go Project Layout and Clean Architecture principles (Handlers, Repositories, Models, and Middleware).

## Features

- **Authentication & Authorization:** Secure user registration and login with bcrypt password hashing and JWT (JSON Web Tokens).
- **Role-Based Access Control (RBAC):** Admin and User roles. Sensitive management endpoints (e.g. Categories) are strictly restricted to Admin users via `AdminMiddleware`.
- **Resource Ownership Verification:** Authors can edit and delete only their own articles and comments.
- **Blog Posts Management:** Full CRUD operations for articles with `LEFT JOIN` queries for categories and authors.
- **Categories System:** Category management for organizing articles.
- **Likes & Favorites System:** Atomic toggle like mechanism (`POST /blogs/:id/like`) preventing duplicate likes via database constraints.
- **Comments System:** Real-time discussion system (`POST`, `GET`, and `DELETE` comments with author validation).
- **Masked Security Error Handling:** Database/SQL errors are masked from client responses to prevent internal information disclosure while detailed logs are recorded server-side.
- **Database Migrations:** Automated database schema versioning managed via `golang-migrate` and custom helper scripts.

---

## Tech Stack

- **Language:** Go 1.26+
- **Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL 14+
- **Driver & Pool:** [pgx v5](https://github.com/jackc/pgx) (`pgxpool`)
- **Authentication:** JWT (golang-jwt/jwt v5) & bcrypt
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)

---

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                  # Application entrypoint & route registration
├── internal/
│   ├── config/
│   │   └── config.go                # Environment variable configuration
│   ├── database/
│   │   └── postgres.go              # PostgreSQL connection pool initialization
│   ├── handlers/
│   │   ├── blog_handler.go          # Blog HTTP handlers
│   │   ├── category_handler.go      # Category HTTP handlers
│   │   ├── comment_handler.go       # Comment HTTP handlers
│   │   ├── like_handler.go          # Like HTTP handlers
│   │   └── user_handler.go          # Authentication HTTP handlers
│   ├── middleware/
│   │   ├── admin_middleware.go      # Admin role authorization middleware
│   │   ├── auth_middleware.go       # JWT Bearer token middleware
│   │   └── cors_middleware.go       # Cross-Origin Resource Sharing middleware
│   ├── models/
│   │   ├── blog.go                  # Blog model
│   │   ├── category.go              # Category model
│   │   ├── comment.go               # Comment model
│   │   ├── like.go                  # Like model
│   │   └── user.go                  # User model
│   └── repository/
│       ├── blog_repository.go       # Blog database queries & JOINs
│       ├── category_repository.go   # Category database queries
│       ├── comment_repository.go    # Comment database queries
│       ├── like_repository.go       # Like toggle & count queries
│       └── user_repository.go       # User database queries
├── migrations/                      # SQL migration files
├── scripts/
│   └── migrate.sh                   # Helper script for migration operations
├── .air.toml                        # Hot reload config
├── .env                             # Environment variables
├── go.mod
└── go.sum
```

---

## Environment Variables

Create a `.env` file in the root directory:

```env
DATABASE_URL=postgresql://postgres:password@localhost:5432/blog_api?sslmode=disable
PORT=8080
JWT_SECRET=your_production_secret_key
```

---

## Database Migrations

Run database migrations using the included helper script:

```bash
# Apply all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down

# Check migration status
./scripts/migrate.sh status
```

---

## API Documentation

### Health Check

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET`  | `/`      | None | Server health check and database status |

---

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/auth/register` | None | Register a new user (default role: `user`) |
| `POST` | `/auth/login`    | None | Authenticate user and receive JWT token |

#### `POST /auth/register`
```json
// Request Body
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}

// Response (201 Created)
{
  "id": "e82e7f45-fef3-4c81-9524-0bdee7726d71",
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2026-07-21T11:00:00Z"
}
```

#### `POST /auth/login`
```json
// Request Body
{
  "username": "johndoe",
  "password": "password123"
}

// Response (200 OK)
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### Categories

| Method   | Endpoint          | Auth | Role Access | Description |
|----------|-------------------|------|-------------|-------------|
| `GET`    | `/categories`     | None | Public      | Get all categories |
| `GET`    | `/categories/:id` | None | Public      | Get category by ID |
| `POST`   | `/categories`     | JWT  | Admin Only  | Create a new category |
| `PUT`    | `/categories/:id` | JWT  | Admin Only  | Update a category |
| `DELETE` | `/categories/:id` | JWT  | Admin Only  | Delete a category |

---

### Blog Posts

| Method   | Endpoint          | Auth | Description |
|----------|-------------------|------|-------------|
| `GET`    | `/blogs`          | None | Get all blogs with author, category, and total likes |
| `GET`    | `/blogs/:id`      | None | Get single blog by ID |
| `GET`    | `/blogs/user`     | JWT  | Get blogs authored by currently logged-in user |
| `POST`   | `/blogs`          | JWT  | Create a new blog post |
| `PUT`    | `/blogs/:id`      | JWT  | Update blog post (Owner only) |
| `DELETE` | `/blogs/:id`      | JWT  | Delete blog post (Owner only) |

#### `GET /blogs` Response Example:
```json
{
  "message": "Blogs retrieved successfully",
  "blogs": [
    {
      "id": 1,
      "title": "System Architecture Patterns",
      "content": "Exploring distributed systems and database pooling in Go.",
      "created_at": "2026-07-21T10:00:00Z",
      "updated_at": "2026-07-21T10:00:00Z",
      "user_id": "ae70c3bc-891c-4bd0-8fae-1af577867973",
      "author_name": "johndoe",
      "category_id": 1,
      "category_name": "Engineering",
      "total_likes": 12
    }
  ]
}
```

---

### Likes

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/blogs/:id/like` | JWT | Toggle like/unlike on a blog post |

#### `POST /blogs/:id/like` Response Example:
```json
{
  "message": "Blog liked successfully",
  "is_liked": true,
  "total_likes": 13
}
```

---

### Comments

| Method   | Endpoint             | Auth | Authorization | Description |
|----------|----------------------|------|---------------|-------------|
| `GET`    | `/blogs/:id/comments`| None | Public        | Get all comments for a blog post |
| `POST`   | `/blogs/:id/comments`| JWT  | User/Admin    | Post a comment on a blog post |
| `DELETE` | `/comments/:id`      | JWT  | Owner or Admin| Delete a comment |

#### `GET /blogs/:id/comments` Response Example:
```json
{
  "message": "Comments retreived successfully",
  "comments": [
    {
      "id": 1,
      "blog_id": 1,
      "user_id": "e82e7f45-fef3-4c81-9524-0bdee7726d71",
      "user_name": "rayhan",
      "content": "Great article on Go connection pooling!",
      "created_at": "2026-07-21T11:20:00Z"
    }
  ]
}
```

---

## Roadmap Status

- [x] Blog Posts CRUD & Ownership Verification
- [x] User Authentication & Authorization (JWT & bcrypt)
- [x] Categories Management & Database Relations (`LEFT JOIN`)
- [x] Role-Based Access Control (Admin-only Category Management)
- [x] Atomic Likes System (Toggle Like & Count)
- [x] Comments System (Create, Read with Author JOIN, Delete Authorization)
- [x] Masked Security Error Handling
- [x] CORS Middleware Support for Frontend Clients

---

## License

This project is open-source software under the MIT License.
