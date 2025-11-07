# Deployment Guide

This guide covers deploying your Go + SQLC API to various platforms.

## Table of Contents
- [Railway](#railway)
- [Render](#render)
- [Fly.io](#flyio)
- [DigitalOcean](#digitalocean)
- [AWS ECS](#aws-ecs)

---

## Railway

Railway offers one of the easiest deployment experiences with built-in PostgreSQL.

### Steps

1. **Install Railway CLI**
   ```bash
   npm install -g @railway/cli
   railway login
   ```

2. **Initialize Project**
   ```bash
   railway init
   railway link
   ```

3. **Add PostgreSQL**
   ```bash
   railway add postgresql
   ```

4. **Set Environment Variables**
   ```bash
   railway variables set JWT_SECRET=$(openssl rand -hex 32)
   railway variables set ENV=production
   railway variables set CORS_ALLOWED_ORIGINS=https://yourfrontend.com
   ```

5. **Deploy**
   ```bash
   railway up
   ```

Railway automatically detects the Dockerfile and deploys your application!

### Getting DATABASE_URL

Railway automatically injects `DATABASE_URL` for PostgreSQL. No manual configuration needed.

---

## Render

Render provides free tier with automatic deploys from GitHub.

### Steps

1. **Push to GitHub**
   ```bash
   git push origin main
   ```

2. **Create New Web Service**
   - Go to [render.com](https://render.com)
   - Click "New +" → "Web Service"
   - Connect your GitHub repository

3. **Configure Service**
   - **Name**: your-api-name
   - **Environment**: Docker
   - **Build Command**: (auto-detected from Dockerfile)
   - **Start Command**: (auto-detected)

4. **Add PostgreSQL**
   - Click "New +" → "PostgreSQL"
   - Copy the Internal Database URL

5. **Set Environment Variables**
   ```
   DATABASE_URL=<your-postgres-url>
   JWT_SECRET=<generate-with-openssl-rand-hex-32>
   ENV=production
   CORS_ALLOWED_ORIGINS=*
   ```

6. **Deploy**
   - Click "Create Web Service"
   - Render automatically deploys on every push to main

---

## Fly.io

Fly.io offers edge deployment with global distribution.

### Steps

1. **Install Fly CLI**
   ```bash
   curl -L https://fly.io/install.sh | sh
   fly auth signup
   ```

2. **Launch Application**
   ```bash
   fly launch
   ```
   Answer the prompts:
   - App name: your-api-name
   - Organization: your-org
   - Region: choose closest to users
   - PostgreSQL: Yes (recommended)

3. **Set Secrets**
   ```bash
   fly secrets set JWT_SECRET=$(openssl rand -hex 32)
   fly secrets set ENV=production
   ```

4. **Deploy**
   ```bash
   fly deploy
   ```

### Scaling
```bash
fly scale count 2  # Run 2 instances
fly scale vm shared-cpu-1x  # Change machine size
```

---

## DigitalOcean

Deploy to a VPS with full control.

### Steps

1. **Create Droplet**
   - Go to DigitalOcean dashboard
   - Create Droplet (Ubuntu 22.04)
   - Choose size ($6/month minimum)

2. **SSH into Server**
   ```bash
   ssh root@your-droplet-ip
   ```

3. **Install Dependencies**
   ```bash
   # Update system
   apt update && apt upgrade -y

   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sh get-docker.sh

   # Install Docker Compose
   apt install docker-compose -y

   # Install PostgreSQL
   apt install postgresql postgresql-contrib -y
   ```

4. **Setup PostgreSQL**
   ```bash
   sudo -u postgres psql
   CREATE DATABASE go_api_db;
   CREATE USER api_user WITH ENCRYPTED PASSWORD 'your-password';
   GRANT ALL PRIVILEGES ON DATABASE go_api_db TO api_user;
   \q
   ```

5. **Clone and Deploy**
   ```bash
   git clone https://github.com/yourusername/go-sqlc-starter.git
   cd go-sqlc-starter

   # Create .env file
   cp .env.example .env
   nano .env  # Edit with your values

   # Build and run
   docker-compose up -d
   ```

6. **Setup Nginx (Optional)**
   ```bash
   apt install nginx -y

   # Create Nginx config
   cat > /etc/nginx/sites-available/api << 'EOF'
   server {
       listen 80;
       server_name your-domain.com;

       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   EOF

   ln -s /etc/nginx/sites-available/api /etc/nginx/sites-enabled/
   nginx -t
   systemctl restart nginx
   ```

7. **Setup SSL with Let's Encrypt**
   ```bash
   apt install certbot python3-certbot-nginx -y
   certbot --nginx -d your-domain.com
   ```

---

## AWS ECS

Enterprise-grade deployment with AWS ECS and RDS.

### Prerequisites
- AWS CLI installed and configured
- Docker installed locally

### Steps

1. **Create ECR Repository**
   ```bash
   aws ecr create-repository --repository-name go-api
   ```

2. **Build and Push Image**
   ```bash
   # Login to ECR
   aws ecr get-login-password --region us-east-1 | \
     docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

   # Build
   docker build -t go-api .

   # Tag
   docker tag go-api:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/go-api:latest

   # Push
   docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/go-api:latest
   ```

3. **Create RDS PostgreSQL Instance**
   - Go to AWS RDS Console
   - Create Database → PostgreSQL
   - Choose dev/test or production template
   - Note the endpoint and credentials

4. **Create ECS Cluster**
   ```bash
   aws ecs create-cluster --cluster-name go-api-cluster
   ```

5. **Create Task Definition**
   Create `task-definition.json`:
   ```json
   {
     "family": "go-api-task",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "256",
     "memory": "512",
     "containerDefinitions": [
       {
         "name": "go-api",
         "image": "<account-id>.dkr.ecr.us-east-1.amazonaws.com/go-api:latest",
         "portMappings": [
           {
             "containerPort": 8080,
             "protocol": "tcp"
           }
         ],
         "environment": [
           {"name": "PORT", "value": "8080"},
           {"name": "ENV", "value": "production"}
         ],
         "secrets": [
           {
             "name": "DATABASE_URL",
             "valueFrom": "arn:aws:secretsmanager:region:account:secret:db-url"
           },
           {
             "name": "JWT_SECRET",
             "valueFrom": "arn:aws:secretsmanager:region:account:secret:jwt-secret"
           }
         ]
       }
     ]
   }
   ```

   Register task:
   ```bash
   aws ecs register-task-definition --cli-input-json file://task-definition.json
   ```

6. **Create Service**
   ```bash
   aws ecs create-service \
     --cluster go-api-cluster \
     --service-name go-api-service \
     --task-definition go-api-task \
     --desired-count 2 \
     --launch-type FARGATE \
     --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx],securityGroups=[sg-xxx],assignPublicIp=ENABLED}"
   ```

---

## Environment Variables Checklist

Ensure these are set in production:

- [ ] `DATABASE_URL` - PostgreSQL connection string
- [ ] `JWT_SECRET` - Strong random secret (min 32 characters)
- [ ] `ENV=production`
- [ ] `CORS_ALLOWED_ORIGINS` - Specific origins, not `*`
- [ ] `PORT` - Usually 8080
- [ ] `JWT_ACCESS_EXPIRY=15m`
- [ ] `JWT_REFRESH_EXPIRY=168h`

## Security Checklist

- [ ] Use strong JWT secret (32+ random characters)
- [ ] Enable HTTPS/TLS
- [ ] Set specific CORS origins in production
- [ ] Use environment variables, never hardcode secrets
- [ ] Enable database SSL connections
- [ ] Set up monitoring and logging
- [ ] Configure firewall rules
- [ ] Regular security updates

## Monitoring

Consider adding:
- Application logs (CloudWatch, Datadog, LogRocket)
- Uptime monitoring (UptimeRobot, Pingdom)
- Error tracking (Sentry)
- Performance monitoring (New Relic)

---

Need help? Check the main README or open an issue!
