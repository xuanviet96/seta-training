# 🚀 Seta Training - Go REST API

A modern Go REST API service for post management with full-text search capabilities.

## 🏗️ **Current Implementation Status**

✅ **REST API for Post Management** - Fully implemented and working!

This project currently implements a complete REST API for managing posts with:
- **Database Integration**: PostgreSQL with GORM
- **Caching**: Redis for performance optimization
- **Search**: Elasticsearch for full-text search
- **API Framework**: Gin web framework
- **Background Processing**: Async indexing
- **Health Monitoring**: Service health checks

## 🎯 Project Scope

This is a Go REST API training project that demonstrates:
- Building REST APIs with Gin framework
- PostgreSQL integration with GORM
- Redis caching for performance
- Elasticsearch for full-text search
- Clean architecture patterns
- Docker containerization

---

## 🚀 **Quick Start**

### Prerequisites

- **Go 1.22+**
- **Docker & Docker Compose**
- **Git**

### 🔧 **Setup & Installation**

```bash
# 1. Clone the repository
git clone <repository-url>
cd seta-training

# 2. Start dependencies (PostgreSQL, Redis, Elasticsearch)
docker compose up -d postgres redis elasticsearch

# 3. Wait for services to be healthy (about 30 seconds)
docker compose ps

# 4. Run database migrations
Get-Content migrations/0001_init_sql | docker exec -i seta-training-postgres-1 psql -U postgres -d blog

# 5. Build and run the server
go mod tidy
go build -o bin/server.exe ./cmd/server
./bin/server.exe
```

**Or use the convenient script:**
```powershell
./start_server.ps1
```

### 🧪 **Testing the API**

Run the automated test script:
```powershell
./test_api.ps1
```

Or test manually:
```powershell
# Health check
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET

# Create a post
$postData = '{"title":"My Post","content":"Content here","tags":["tag1","tag2"]}'
Invoke-RestMethod -Uri "http://localhost:8080/v1/posts" -Method POST -Body $postData -ContentType "application/json"
```

---

## 📋 **Currently Implemented API Endpoints**

### 🏥 **Health & Monitoring**

| Method | Path       | Description                               |
|--------|-----------|-------------------------------------------|
| GET    | `/health` | Check service health (DB, Redis, ES)     |

### 📝 **Post Management**

| Method | Path                          | Description                    |
|--------|-------------------------------|--------------------------------|
| POST   | `/v1/posts`                   | Create a new post             |
| GET    | `/v1/posts/:id`               | Get post by ID                |
| PUT    | `/v1/posts/:id`               | Update post                   |
| GET    | `/v1/posts/search-by-tag`     | Search posts by tag           |
| GET    | `/v1/posts/search`            | Full-text search with ES      |

### 📋 **Request/Response Examples**

#### Create Post
```json
POST /v1/posts
{
  "title": "My First Post",
  "content": "This is the content of my post",
  "tags": ["golang", "api", "tutorial"]
}
```

#### Response
```json
{
  "id": 1,
  "title": "My First Post",
  "content": "This is the content of my post",
  "tags": ["golang", "api", "tutorial"],
  "created_at": "2025-09-17T18:27:30.252Z"
}
```

#### Search by Tag
```
GET /v1/posts/search-by-tag?tag=golang
```

#### Full-text Search
```
GET /v1/posts/search?q=tutorial
```

---

## 🏗️ **Architecture & Tech Stack**

### **Backend Technologies**
- **Language**: Go 1.22
- **Web Framework**: Gin
- **Database**: PostgreSQL 16 with GORM
- **Cache**: Redis 7
- **Search**: Elasticsearch 8.13
- **Configuration**: Viper + .env files
- **Logging**: Zap

### **Project Structure**
```
├── cmd/server/          # Application entry point
├── internal/
│   ├── cache/          # Redis cache logic
│   ├── config/         # Configuration management
│   ├── database/       # Database connection
│   ├── domain/
│   │   ├── models/     # Data models
│   │   ├── repository/ # Data access layer
│   │   └── services/   # Business logic
│   ├── http/
│   │   ├── handlers/   # HTTP handlers
│   │   ├── middleware/ # HTTP middleware
│   │   └── router.go   # Route definitions
│   ├── logger/         # Logging setup
│   └── search/         # Elasticsearch integration
├── migrations/         # Database migrations
├── pkg/               # Shared packages
├── docker-compose.yaml # Development services
├── start_server.ps1   # Server startup script
└── test_api.ps1       # API testing script
```

### **Database Schema**
```sql
-- Posts table
CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  title VARCHAR NOT NULL,
  content TEXT NOT NULL,
  tags TEXT[] NOT NULL DEFAULT '{}',
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Activity logs for audit trail
CREATE TABLE activity_logs (
  id SERIAL PRIMARY KEY,
  action VARCHAR NOT NULL,
  post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  logged_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- GIN index for fast tag searches
CREATE INDEX idx_posts_tags_gin ON posts USING GIN (tags);
```

---

## 🔧 **Configuration**

The application uses a `.env` file for configuration:

```env
# Application Configuration
APP_PORT=8080
APP_ENV=development

# Database Configuration  
DATABASE_URL=postgres://postgres:postgres@localhost:5432/blog?sslmode=disable

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_DB=0
REDIS_TTL_SECONDS=300

# Elasticsearch Configuration
ES_ADDR=http://localhost:9200
ES_INDEX=posts
```

---

## 🛠️ **Development Commands**

```powershell
# Install dependencies
go mod tidy

# Build the application
go build -o bin/server.exe ./cmd/server

# Run the application
go run ./cmd/server

# Run tests
go test -v ./...

# Start services with Docker
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

---

## � **Learning Outcomes**

This project demonstrates key Go development concepts:

- **REST API Development**: Building HTTP endpoints with Gin framework
- **Database Integration**: Using GORM for PostgreSQL operations
- **Caching Strategies**: Implementing Redis for performance optimization
- **Search Implementation**: Integrating Elasticsearch for full-text search
- **Clean Architecture**: Separating concerns with proper layering
- **Background Processing**: Async operations for search indexing
- **Error Handling**: Proper HTTP error responses and validation
- **Testing**: API endpoint testing and validation
- **Containerization**: Docker setup for development environment

---

## 🤝 **Contributing**

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 **License**

This project is part of SETA training program.

---

## 🆘 **Troubleshooting**

### Common Issues

1. **Port already in use**: Make sure port 8080 is not used by another service
2. **Docker services not starting**: Run `docker compose down` then `docker compose up -d`
3. **Database connection failed**: Check if PostgreSQL container is healthy with `docker compose ps`
4. **Elasticsearch not responding**: ES takes longer to start, wait ~30 seconds after `docker compose up`

### Useful Commands

```bash
# Check service status
docker compose ps

# View logs
docker compose logs postgres
docker compose logs redis  
docker compose logs elasticsearch

# Restart specific service
docker compose restart postgres

# Reset everything
docker compose down -v
docker compose up -d
```

