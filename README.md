# 🚀 Production-Ready Go + SQLC API Starter Kit

A battle-tested, production-ready REST API starter built with Go, SQLC, PostgreSQL, and modern best practices. Save 20+ hours of boilerplate setup and start building features immediately.

## ⚡ What You Get

This isn't a tutorial project—this is **production-grade code** that handles everything you need:

- ✅ **Complete REST API** with Gin framework
- ✅ **Type-safe database queries** with SQLC (no ORMs!)
- ✅ **JWT Authentication** with refresh tokens
- ✅ **PostgreSQL** with automatic migrations
- ✅ **Docker & Docker Compose** setup
- ✅ **Comprehensive testing** examples
- ✅ **Rate limiting** and middleware
- ✅ **Request validation** with proper error handling
- ✅ **Graceful shutdown** handling
- ✅ **Environment-based configuration**
- ✅ **Structured logging** with zerolog
- ✅ **Health check endpoints**
- ✅ **Database transactions** examples
- ✅ **CI/CD ready** (GitHub Actions included)
- ✅ **Deployment guides** for Railway, Render, Fly.io

## 🎯 Who Is This For?

- Backend developers tired of setting up the same boilerplate
- Teams starting new Go microservices
- Developers wanting to learn Go best practices
- Anyone building a REST API who values type safety

## 📁 Project Structure

```
go-sqlc-starter/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── router.go            # Route definitions
│   │   ├── middleware/          # Auth, CORS, rate limiting
│   │   └── handlers/            # HTTP handlers
│   ├── db/
│   │   ├── sqlc/                # Generated SQLC code
│   │   ├── migrations/          # SQL migrations
│   │   └── queries/             # SQLC queries
│   ├── auth/
│   │   ├── jwt.go               # JWT generation/validation
│   │   └── password.go          # Password hashing
│   ├── models/                  # Request/response models
│   ├── config/                  # Configuration management
│   └── validator/               # Request validation
├── tests/
│   ├── integration/             # Integration tests
│   └── unit/                    # Unit tests
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── docs/
│   ├── DEPLOYMENT.md            # Deployment guides
│   ├── TESTING.md               # Testing guide
│   └── API.md                   # API documentation
├── scripts/
│   ├── setup.sh                 # Initial setup script
│   └── migrate.sh               # Migration helper
├── .github/
│   └── workflows/
│       └── ci.yml               # CI/CD pipeline
├── sqlc.yaml                    # SQLC configuration
├── Makefile                     # Common commands
└── README.md
```

## 🏁 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make (optional but recommended)

### Setup in 2 Minutes

```bash
# Clone the repository
git clone <your-repo>
cd go-sqlc-starter

# Copy environment variables
cp .env.example .env

# Start PostgreSQL
docker-compose up -d postgres

# Run migrations
make migrate-up

# Start the server
make run
```

The API will be running at `http://localhost:8080` 🎉

## 🔧 Key Features Explained

### 1. Type-Safe Database with SQLC

No more `db.Query()` and manual scanning. Write SQL, get type-safe Go code:

```sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name)
VALUES ($1, $2, $3)
RETURNING *;
```

Generates:
```go
user, err := queries.GetUser(ctx, userID)
// user is a fully typed struct!
```

### 2. JWT Authentication with Refresh Tokens

Complete auth system with:
- Access tokens (15 min expiry)
- Refresh tokens (7 days expiry)
- Token rotation on refresh
- Secure password hashing with bcrypt

### 3. Database Migrations

Simple, version-controlled migrations:
```bash
make migrate-up      # Apply all pending migrations
make migrate-down    # Rollback last migration
make migrate-create  # Create new migration
```

### 4. Comprehensive Testing

Both unit and integration tests included:
```bash
make test            # Run all tests
make test-coverage   # With coverage report
make test-integration # Integration tests only
```

### 5. Production-Ready Middleware

- CORS handling
- Rate limiting (configurable per route)
- Request logging
- Panic recovery
- JWT validation
- Request ID tracking

## 📚 API Endpoints

### Authentication
```
POST   /api/v1/auth/register    # Register new user
POST   /api/v1/auth/login       # Login
POST   /api/v1/auth/refresh     # Refresh access token
POST   /api/v1/auth/logout      # Logout
```

### Users (Protected)
```
GET    /api/v1/users/me         # Get current user
PUT    /api/v1/users/me         # Update current user
DELETE /api/v1/users/me         # Delete account
GET    /api/v1/users/:id        # Get user by ID (admin)
```

### Health
```
GET    /health                  # Health check
GET    /ready                   # Readiness check
```

## 🐳 Docker Deployment

### Development
```bash
docker-compose up
```

### Production
```bash
docker build -f docker/Dockerfile -t api:latest .
docker run -p 8080:8080 --env-file .env api:latest
```

## ☁️ Cloud Deployment

Detailed guides included for:
- **Railway**: One-click deploy with PostgreSQL
- **Render**: Free tier compatible
- **Fly.io**: Edge deployment
- **DigitalOcean**: VPS setup
- **AWS ECS**: Container orchestration

See `docs/DEPLOYMENT.md` for step-by-step instructions.

## 🔐 Security Features

- ✅ Password hashing with bcrypt (cost 12)
- ✅ JWT with HMAC-SHA256
- ✅ SQL injection prevention (parameterized queries)
- ✅ CORS configuration
- ✅ Rate limiting
- ✅ Secure headers middleware
- ✅ Input validation
- ✅ No credentials in logs

## ⚙️ Configuration

All configuration via environment variables (12-factor app):

```env
# Server
PORT=8080
ENV=development

# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/db

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

## 🧪 Testing

```bash
# Unit tests
make test-unit

# Integration tests (requires PostgreSQL)
make test-integration

# Coverage report
make test-coverage

# Benchmark tests
make benchmark
```

## 📖 Documentation

- `docs/API.md` - Complete API documentation with examples
- `docs/DEPLOYMENT.md` - Step-by-step deployment guides
- `docs/TESTING.md` - Testing strategies and examples
- `docs/ARCHITECTURE.md` - Design decisions and patterns

## 🛠️ Development Commands

```bash
make run              # Run the server
make build            # Build binary
make test             # Run tests
make migrate-up       # Apply migrations
make migrate-down     # Rollback migrations
make sqlc-generate    # Regenerate SQLC code
make docker-build     # Build Docker image
make lint             # Run linters
make fmt              # Format code
```

## 🎓 Learning Resources

The code includes extensive comments explaining:
- Why certain patterns are used
- Common pitfalls to avoid
- Performance considerations
- Security best practices

## 🤝 What's NOT Included (By Design)

This starter focuses on the backend API. You'll need to add:
- Frontend (React, Vue, etc.)
- Email service integration
- File upload handling
- WebSocket support
- Caching layer (Redis)

These are intentionally left out to keep the starter focused and not opinionated about your specific needs.

## 📝 License

This is a **commercial product**. By purchasing, you get:
- Unlimited personal and commercial use
- Access to all updates
- Email support for 30 days
- Source code (no attribution required)

NOT allowed:
- Reselling this starter kit
- Creating competing starter kits using this code

## 🚀 Get Started Now

This starter kit saves you 20+ hours of setup time. Instead of configuring SQLC, setting up migrations, implementing auth, and writing middleware—just start building your actual features.

**Time is money. Start shipping today.**

---

Built with ❤️ by a developer who was tired of copy-pasting the same setup code.
