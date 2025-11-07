# 🚀 Quick Start Guide

Get your API running in **5 minutes**!

## Prerequisites

- Go 1.21+ installed ([download](https://golang.org/dl/))
- Docker installed ([download](https://docs.docker.com/get-docker/))
- Git installed

## Step 1: Clone and Setup (2 minutes)

```bash
# Clone the repository
git clone <your-repo-url>
cd go-sqlc-starter

# Run setup script (installs dependencies, generates code, starts database)
chmod +x scripts/setup.sh
./scripts/setup.sh
```

That's it! The script handles everything:
- ✅ Installs dependencies
- ✅ Generates SQLC code
- ✅ Starts PostgreSQL in Docker
- ✅ Creates .env file

## Step 2: Run the API (30 seconds)

```bash
make run
```

You should see:
```
Starting application...
Database connection established
Database migrations completed
Server starting on :8080
```

## Step 3: Test It (1 minute)

### Check Health

```bash
curl http://localhost:8080/health
```

Response:
```json
{"status":"healthy","env":"development"}
```

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

You'll get back an access token!

### Make an Authenticated Request

```bash
# Save your token from the register response
TOKEN="your-access-token-here"

curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

## 🎉 Success!

Your API is running! Here's what to do next:

### Customize for Your Project

1. **Change the module name** in `go.mod`
   ```go
   module github.com/yourusername/your-project-name
   ```

2. **Update database schema** in `internal/db/migrations/`
   ```bash
   make migrate-create name=add_your_table
   # Edit the new migration files
   make migrate-up
   ```

3. **Add your queries** in `internal/db/queries/`
   ```sql
   -- name: GetYourData :many
   SELECT * FROM your_table;
   ```

4. **Regenerate SQLC code**
   ```bash
   make sqlc-generate
   ```

5. **Add your handlers** in `internal/api/handlers/`

## Common Commands

```bash
make run              # Start the server
make test             # Run tests
make migrate-up       # Apply migrations
make migrate-down     # Rollback last migration
make sqlc-generate    # Regenerate SQLC code
make docker-up        # Start all services
make docker-down      # Stop all services
```

## Environment Configuration

Edit `.env` to configure:

```env
PORT=8080
DATABASE_URL=postgresql://...
JWT_SECRET=your-secret-here
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h
```

## Project Structure

```
.
├── cmd/api/          # Application entry point
├── internal/
│   ├── api/          # HTTP handlers & routes
│   ├── auth/         # Authentication logic
│   ├── config/       # Configuration
│   └── db/           # Database & SQLC
├── docs/             # Documentation
└── scripts/          # Helper scripts
```

## Next Steps

1. **Read the docs**
   - [API Documentation](docs/API.md) - All endpoints with examples
   - [Deployment Guide](docs/DEPLOYMENT.md) - Deploy to production

2. **Customize for your needs**
   - Add your business logic
   - Create your database schema
   - Build your endpoints

3. **Deploy**
   - Follow the deployment guide
   - Railway, Render, Fly.io all work great

## Need Help?

- Check out the full README.md
- Look at the docs/ folder
- Email support: [your-email@example.com]

## Troubleshooting

### "Connection refused" error

PostgreSQL isn't running. Start it:
```bash
docker-compose up -d postgres
```

### "SQLC not found"

Install SQLC:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Port 8080 already in use

Change the port in `.env`:
```env
PORT=3000
```

---

**You're all set! Happy coding! 🎉**
