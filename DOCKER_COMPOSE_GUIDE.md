# Docker Compose Quick Start

## Run Everything with Docker Compose

Docker Compose will build and run all services together.

### Prerequisites
- Docker installed
- Docker Compose installed

### Quick Start (One Command)

```bash
# 1. Copy environment file (already has defaults)
cp .env.example .env

# 2. Build and start all services
docker-compose up -d --build
```

### What Gets Started

| Service | Port | Description |
|---------|------|-------------|
| PostgreSQL | 5433 | Database |
| Redis | 6380 | Job queue |
| Backend | 8081 | API Server |
| Seeder | - | (Runs once, then exits) |

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f postgres
docker-compose logs -f redis
```

### Seed the Database

The seeder runs as a service but needs to be triggered manually:

```bash
# Run the seeder (creates 5 test accounts)
docker-compose run --rm seeder

# Or use the tools profile
docker-compose --profile tools up seeder
```

### Access the Application

- **API**: http://localhost:8081
- **Health Check**: http://localhost:8081/api/health
- **Swagger Docs**: http://localhost:8081/swagger/index.html
- **PostgreSQL**: localhost:5433
- **Redis**: localhost:6380

### Test Accounts

After seeding, you can login with:

```bash
# Admin
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@admin.com","password":"admin123"}'

# Student
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"student1@innogen.com","password":"student123"}'

# Teacher
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"teacher@innogen.com","password":"teacher123"}'
```

### Stop Everything

```bash
docker-compose down
```

### Clean Everything (including volumes)

```bash
docker-compose down -v
docker-compose rm -f
```

## Docker Compose Configuration

### Services

#### 1. Backend Service
- **Image**: Built from Dockerfile
- **Port**: 8081:8080
- **Dependencies**: PostgreSQL (healthy), Redis (healthy)
- **Health Check**: HTTP /api/health endpoint
- **Auto-restart**: Unless stopped

#### 2. PostgreSQL Service
- **Image**: postgres:15
- **Port**: 5433:5432
- **Database**: innogendb
- **User**: innogen
- **Password**: maiphuongdangyeu
- **Health Check**: pg_isready

#### 3. Redis Service
- **Image**: redis:7
- **Port**: 6380:6379
- **Health Check**: redis-cli ping

#### 4. Seeder Service (Profile: tools)
- **Image**: Built from Dockerfile
- **Command**: ./innogen-seed
- **Purpose**: Creates test accounts and sample data
- **Runs**: Once, then exits

### Environment Variables

All services use the `.env` file via `env_file`.

#### Required Variables:
- `JWT_SECRET`: Access token secret (32+ chars)
- `JWT_REFRESH_SECRET`: Refresh token secret (32+ chars)
- `ADMIN_PASSWORD`: Admin user password
- `POSTGRES_*`: Database configuration
- `REDIS_*`: Redis configuration

### Networks

All services run in `innogen-network` bridge network.

### Volumes

- `postgres_data`: PostgreSQL data persistence
- `redis_data`: Redis data persistence

## Development Workflow

### 1. First Time Setup

```bash
# Start all services
docker-compose up -d --build

# Wait for services to be healthy
docker-compose ps

# Run the seeder
docker-compose run --rm seeder
```

### 2. Development

```bash
# View logs
docker-compose logs -f backend

# Restart backend only
docker-compose restart backend

# Rebuild backend after code changes
docker-compose build backend
docker-compose up -d backend
```

### 3. Making Changes

#### Backend Code Changes
```bash
# Rebuild and restart backend
docker-compose build backend
docker-compose up -d backend
```

#### Database Changes
```bash
# GORM auto-migrate runs on startup
# If needed, restart backend
docker-compose restart backend
```

#### Seeding Changes
```bash
# Re-run seeder (it's safe to run multiple times)
docker-compose run --rm seeder
```

## Troubleshooting

### Port Already in Use

```bash
# Find what's using the port
lsof -i :8081

# Or change ports in docker-compose.yml
# backend: "8082:8080"
```

### Database Connection Failed

```bash
# Check PostgreSQL is healthy
docker-compose ps
docker-compose logs postgres

# Wait a bit longer for PostgreSQL to initialize
sleep 5
```

### Seeder Fails

```bash
# Check backend is running
docker-compose ps

# Check database is healthy
docker-compose exec postgres pg_isready -U innogen

# Try running seeder again
docker-compose run --rm seeder
```

### Clear Everything and Start Fresh

```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d --build
docker-compose run --rm seeder
```

## Production Considerations

### 1. Security
- Change default passwords
- Use strong JWT secrets
- Use environment-specific .env files
- Enable HTTPS
- Restrict database access

### 2. Performance
- Use production PostgreSQL instance
- Add Redis persistence
- Configure connection pooling
- Add Redis for caching

### 3. Monitoring
- Add log aggregation
- Health check monitoring
- Metrics collection
- Alerting

## Commands Reference

```bash
# Start all services
docker-compose up -d

# Build all services
docker-compose build

# Build specific service
docker-compose build backend

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Restart specific service
docker-compose restart backend

# Run seeder
docker-compose run --rm seeder

# Check service status
docker-compose ps

# Enter container
docker-compose exec backend sh
docker-compose exec postgres psql -U innogen -d innogendb

# Scale service
docker-compose up -d --scale backend=2
```