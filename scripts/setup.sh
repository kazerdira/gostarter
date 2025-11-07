#!/bin/bash

# Go + SQLC Starter Kit Setup Script
# This script helps you get started quickly

set -e

echo "🚀 Go + SQLC Starter Kit Setup"
echo "================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker is not installed. You'll need it to run PostgreSQL."
    echo "   Visit: https://docs.docker.com/get-docker/"
fi

# Check if make is installed
if ! command -v make &> /dev/null; then
    echo "⚠️  Make is not installed. You can still run commands manually."
fi

echo ""
echo "📦 Installing Go dependencies..."
go mod download
echo "✅ Dependencies installed"

echo ""
echo "🔧 Installing development tools..."

# Install SQLC
if ! command -v sqlc &> /dev/null; then
    echo "Installing SQLC..."
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    echo "✅ SQLC installed"
else
    echo "✅ SQLC already installed"
fi

# Install migrate
if ! command -v migrate &> /dev/null; then
    echo "Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    echo "✅ golang-migrate installed"
else
    echo "✅ golang-migrate already installed"
fi

echo ""
echo "📝 Setting up environment..."

if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ Created .env file from template"
    
    # Generate a random JWT secret
    if command -v openssl &> /dev/null; then
        JWT_SECRET=$(openssl rand -hex 32)
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i '' "s/your-secret-key-change-this-in-production-use-at-least-32-characters/$JWT_SECRET/" .env
        else
            # Linux
            sed -i "s/your-secret-key-change-this-in-production-use-at-least-32-characters/$JWT_SECRET/" .env
        fi
        echo "✅ Generated random JWT secret"
    fi
else
    echo "✅ .env file already exists"
fi

echo ""
echo "🐳 Starting PostgreSQL with Docker..."

if command -v docker-compose &> /dev/null; then
    docker-compose up -d postgres
    echo "✅ PostgreSQL is starting..."
    echo "   Waiting for PostgreSQL to be ready..."
    sleep 5
elif command -v docker &> /dev/null; then
    docker run -d \
        --name go_api_postgres \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_PASSWORD=postgres \
        -e POSTGRES_DB=go_api_db \
        -p 5432:5432 \
        postgres:15-alpine
    echo "✅ PostgreSQL is starting..."
    sleep 5
else
    echo "⚠️  Docker not available. Please start PostgreSQL manually."
fi

echo ""
echo "🔄 Generating SQLC code..."
sqlc generate
echo "✅ SQLC code generated"

echo ""
echo "================================"
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Review .env file and update if needed"
echo "2. Start the API: make run"
echo "3. Test the health endpoint: curl http://localhost:8080/health"
echo ""
echo "📚 Documentation:"
echo "   API Docs: docs/API.md"
echo "   Deployment: docs/DEPLOYMENT.md"
echo ""
echo "🎉 Happy coding!"
