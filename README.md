# SaaS Microcourses Management

A microservices-based SaaS application for managing companies and microcourses, built with Go, gRPC, PostgreSQL, and dbmate for database migrations.

## Features

- **Company Management**: Full CRUD operations with UUID-based identifiers
- **gRPC API**: Protocol buffer-based API with proper versioning
- **Database Migrations**: Managed with dbmate for version control
- **Docker Compose**: Complete development environment setup
- **PostgreSQL**: Production-ready database with UUID support
- **Type-safe queries**: Using SQLC for code generation

## Company Model

- `id` (string/UUID) - Primary key as UUID
- `name` (string) - Company name
- `subscription_plan` (string) - Subscription plan type
- `created_at` (timestamp) - Creation timestamp
- `updated_at` (timestamp) - Last update timestamp

## Project Structure

```
.
├── proto/               # Protobuf definitions
│   └── company.proto
├── pb/                  # Generated protobuf Go code
│   ├── company.pb.go
│   └── company_grpc.pb.go
├── db/
│   ├── migrations/      # SQL schema migrations
│   └── queries/         # SQL queries for sqlc
├── internal/
│   ├── db/              # Generated sqlc code
│   └── service/         # gRPC service implementation
├── cmd/server/          # Server entry point
└── Makefile             # Build commands
```

## Setup

### Prerequisites

- Go 1.23+ (for local development)
- Docker & Docker Compose (for containerized deployment) 
- [Buf CLI](https://docs.buf.build/installation) (recommended for protobuf management)
- sqlc (installed automatically via go install)

### Quick Start with Docker

The easiest way to run the application is using Docker Compose:

```bash
# Start PostgreSQL and the application
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

This will:
- Start PostgreSQL database on port 5432
- Run database migrations automatically
- Start the gRPC server on port 8080

### Local Development Setup

#### Database Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE saas_microcourses;
```

2. Run the migration:
```sql
\i db/migrations/001_create_companies_table.sql
```

#### Running the Server

1. Set environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/saas_microcourses?sslmode=disable"
export PORT="8080"  # optional, defaults to 8080
```

2. Run the server:
```bash
go run cmd/server/main.go
```

Or use the Makefile:
```bash
make run
```

## Database Migrations with dbmate

This project uses [dbmate](https://github.com/amacneil/dbmate) for database migrations, providing better version control and rollback capabilities.

### dbmate Installation

Install dbmate locally for migration management:

```bash
# Linux/macOS
curl -fsSL -o dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
chmod +x dbmate && mv dbmate ~/.local/bin/

# Set DATABASE_URL for local development
export DATABASE_URL=postgres://admin:password@localhost:5432/saas_microcourses?sslmode=disable
```

### Migration Commands

```bash
# Check migration status
dbmate status

# Create new migration
dbmate new migration_name

# Apply pending migrations
dbmate up

# Rollback last migration
dbmate rollback

# Reset database (rollback all, then apply all)
dbmate down && dbmate up
```

### Migration Format

dbmate migrations use up/down blocks:

```sql
-- migrate:up
CREATE TABLE example (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- migrate:down
DROP TABLE IF EXISTS example;
```

### Container Integration

The application automatically runs pending migrations on startup via the `entrypoint.sh` script. Migrations are applied before the gRPC server starts.

## API Endpoints

The gRPC service provides the following methods:

- `CreateCompany(CreateCompanyRequest) -> CreateCompanyResponse`
- `GetCompany(GetCompanyRequest) -> GetCompanyResponse`  
- `UpdateCompany(UpdateCompanyRequest) -> UpdateCompanyResponse`
- `DeleteCompany(DeleteCompanyRequest) -> DeleteCompanyResponse`
- `ListCompanies(ListCompaniesRequest) -> ListCompaniesResponse`

## Development

### Docker Commands

```bash
# Build the application image
docker build -t saas-company-service .

# Run with custom database
docker run -e DATABASE_URL="your-db-url" -p 8080:8080 saas-company-service

# View application logs
docker-compose logs -f app

# View database logs
docker-compose logs -f postgres

# Restart services
docker-compose restart

# Clean up (remove containers and volumes)
docker-compose down -v
```

### Generate Code

**Using Buf (Recommended):**
```bash
# Generate protobuf code with buf
make buf-generate

# Lint protobuf files
make buf-lint

# Format protobuf files
make buf-format

# Check for breaking changes
make buf-breaking
```

**Legacy protoc command:**
```bash
make proto
```

**Generate database code:**
```bash
make sqlc
```

### Testing with grpcurl

Install grpcurl for testing:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Test the service (server must be running with reflection enabled):
```bash
# List services
grpcurl -plaintext localhost:8080 list

# Create a company
grpcurl -plaintext -d '{"name":"Test Company","subscription_plan":"premium"}' localhost:8080 company.v1.CompanyService/CreateCompany

# Get a company
grpcurl -plaintext -d '{"id":1}' localhost:8080 company.v1.CompanyService/GetCompany

# List companies
grpcurl -plaintext -d '{"page":1,"limit":10}' localhost:8080 company.v1.CompanyService/ListCompanies
```

## Buf Integration

This project uses [Buf](https://buf.build) for modern protobuf management, which provides:

- **Code Generation**: Automated Go code generation with proper versioning
- **Linting**: Strict protobuf style and best practices enforcement  
- **Breaking Change Detection**: API compatibility verification
- **Dependency Management**: Centralized protobuf module dependencies
- **HTTP Annotations**: REST API generation support (ready for gRPC-Gateway)

### Buf Configuration Files

- `buf.yaml`: Main configuration with linting rules and dependencies
- `buf.gen.yaml`: Code generation plugin configuration
- `buf.lock`: Dependency lock file (auto-generated)

### Buf Commands

```bash
# Install buf (if not already installed)
go install github.com/bufbuild/buf/cmd/buf@latest

# Generate code
buf generate

# Lint proto files
buf lint

# Format proto files  
buf format -w

# Check breaking changes against main branch
buf breaking --against '.git#branch=main'
```

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DATABASE_URL` | PostgreSQL connection string | - | Yes |
| `PORT` | Server port | `8080` | No |

### Docker Compose Configuration

The `docker-compose.yml` includes:
- **PostgreSQL 15**: Database with automatic migration
- **Application**: gRPC server with health checks
- **Networking**: Isolated network for services
- **Volumes**: Persistent database storage

## Production Deployment

For production deployment:

1. Update database credentials in `docker-compose.yml`
2. Configure environment-specific settings
3. Set up proper logging and monitoring
4. Enable TLS for gRPC connections
5. Configure backup strategies for PostgreSQL data# saas-microcourses-management
