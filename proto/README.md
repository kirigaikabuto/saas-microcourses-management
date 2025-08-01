# Proto Files Structure

This directory contains the Protocol Buffer definitions for the SaaS Microcourses Management system.

## Directory Structure

```
proto/
├── common/v1/          # Shared types and enums
│   └── common.proto
├── company/v1/         # Company service definitions
│   └── company.proto
├── user/v1/            # User service definitions
│   └── user.proto
└── gen/                # Generated Go code (auto-generated)
    └── proto/
        ├── common/v1/
        ├── company/v1/
        └── user/v1/
```

## Fixed Issues

### user.proto fixes:
- ✅ Fixed field naming: `User User = 1;` → `User user = 1;`
- ✅ Fixed ID type consistency: User.id is string, GetUserRequest.id now matches
- ✅ Fixed copy-paste errors: `ListCompaniesRequest/Response` → `ListUsersRequest/Response`
- ✅ Fixed service method: `ListCompanies` → `ListUsers`
- ✅ Added proper fields: email, company_id for user relationships

### company.proto improvements:
- ✅ Added updated_at timestamp
- ✅ Improved structure and consistency
- ✅ Added proper response types

### Both files:
- ✅ Separated go_package paths to avoid conflicts
- ✅ Organized into versioned directories for future scalability
- ✅ Added common types for reusability

## Adding New Proto Files

To add a new service (e.g., `course`):

1. Create directory: `proto/course/v1/`
2. Create file: `proto/course/v1/course.proto`
3. Use consistent naming:
   ```protobuf
   syntax = "proto3";
   package course.v1;
   option go_package = "github.com/kirigaikabuto/saas-microcourses-management/proto/gen/course/v1;coursev1";
   ```
4. Run: `buf generate`

## Generation Commands

```bash
# Generate all proto files
buf generate

# Generate specific service
buf generate --path proto/company/v1/company.proto

# Lint proto files
buf lint

# Format proto files
buf format -w
```

## Import Paths in Go Code

Update your Go imports to use the new generated paths:

```go
// Company service
import companyv1 "github.com/kirigaikabuto/saas-microcourses-management/proto/gen/proto/company/v1"

// User service  
import userv1 "github.com/kirigaikabuto/saas-microcourses-management/proto/gen/proto/user/v1"

// Common types
import commonv1 "github.com/kirigaikabuto/saas-microcourses-management/proto/gen/proto/common/v1"
```