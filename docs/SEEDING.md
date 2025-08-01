# Database Seeding Guide

This project provides two methods for seeding the database with sample company data:

## Method 1: Go-based Seeding (Recommended)

The Go-based seeding provides more control and better error handling.

### Commands

```bash
# Set database URL
export DATABASE_URL=postgres://admin:password@localhost:5432/saas_microcourses?sslmode=disable

# Seed database with sample companies
make seed

# Clear seed data
make seed-clear
```

### Features

- **Duplicate Detection**: Automatically skips companies that already exist
- **Detailed Logging**: Shows exactly what was created/skipped
- **Safe Operations**: Won't create duplicates on multiple runs
- **Flexible**: Easy to modify seed data in `cmd/seed/main.go`
- **Clean Removal**: Can selectively remove only seed data

### Sample Output

```
Seeding companies...
  ✓ Created company 'Acme Corporation' with plan 'enterprise' (ID: d1d40a52-...)
  ✓ Company 'TechStart Inc' already exists, skipping
  
Seeding completed: 14 created, 1 skipped
Total companies in database: 17
```

## Method 2: SQL Migration-based Seeding

SQL-based seeding through dbmate migrations.

### Commands

```bash
# Apply seed migration (part of normal migration process)
make migrate-up

# Rollback seed migration (if down migration is uncommented)
make migrate-down
```

### Location

Seed data is defined in: `db/migrations/20250801191118_seed_companies.sql`

### Characteristics

- **Migration-based**: Part of normal database migration workflow
- **Version Controlled**: Tracked like any other database migration
- **Simple**: Pure SQL INSERT statements
- **Permanent**: Data persists unless explicitly rolled back

## Seed Data

Both methods create the same sample companies:

| Company Name | Subscription Plan |
|--------------|------------------|
| Acme Corporation | enterprise |
| TechStart Inc | premium |
| Digital Solutions Ltd | basic |
| Innovation Labs | premium |
| Global Enterprises | enterprise |
| StartupCo | basic |
| Enterprise Solutions | enterprise |
| Cloud Services Ltd | premium |
| Data Analytics Inc | basic |
| AI Innovations | premium |
| Future Tech | basic |
| Quantum Computing Co | enterprise |
| Mobile First Ltd | premium |
| Blockchain Ventures | enterprise |
| Green Energy Solutions | basic |

## When to Use Each Method

### Use Go-based Seeding When:
- Development environment setup
- Testing with fresh data
- Need to reset data frequently
- Want duplicate protection
- Need detailed feedback

### Use SQL Migration Seeding When:
- Production/staging initial setup
- One-time data population
- Want seeding as part of deployment
- Need version-controlled seed data

## Development Workflow

Typical development workflow:

```bash
# 1. Start fresh database
make docker-clean && make docker-up

# 2. Apply migrations (includes SQL seeds)
make migrate-up

# 3. Add additional test data with Go seeding
make seed

# 4. Clear test data when needed
make seed-clear

# 5. Reset everything
make migrate-reset && make seed
```

## Customizing Seed Data

### Modifying Go Seeds

Edit `cmd/seed/main.go`:

```go
var seedCompanies = []SeedCompany{
    {"Your Company", "premium"},
    {"Another Company", "basic"},
    // Add more companies here
}
```

### Modifying SQL Seeds

Edit `db/migrations/20250801191118_seed_companies.sql`:

```sql
INSERT INTO companies (name, subscription_plan) VALUES
    ('Your Company', 'premium'),
    ('Another Company', 'basic');
```

## Troubleshooting

### Connection Issues
Ensure DATABASE_URL is set correctly:
```bash
export DATABASE_URL=postgres://admin:password@localhost:5432/saas_microcourses?sslmode=disable
```

### Permission Issues
Make sure the database user has INSERT/DELETE permissions on the companies table.

### Duplicate Key Errors
Go-based seeding prevents duplicates. If using SQL seeding directly, ensure you don't run it multiple times or add appropriate constraints.