# Error Handling Guide

This project uses centralized error message constants to ensure consistency across all services.

## Error Constants Location

All error messages are defined in: `internal/service/errors.go`

## Error Categories

### Validation Errors (InvalidArgument)
- `ErrCompanyNameRequired` - "name is required"
- `ErrCompanySubscriptionRequired` - "subscription_plan is required" 
- `ErrCompanyIDRequired` - "id is required"
- `ErrInvalidUUIDFormat` - "invalid UUID format"

### Operation Errors (Internal)
- `ErrFailedToCreateCompany` - "failed to create company"
- `ErrFailedToGetCompany` - "failed to get company"
- `ErrFailedToUpdateCompany` - "failed to update company"
- `ErrFailedToDeleteCompany` - "failed to delete company"
- `ErrFailedToListCompanies` - "failed to list companies"
- `ErrFailedToCountCompanies` - "failed to count companies"

### Business Logic Errors (NotFound)
- `ErrCompanyNotFound` - "company not found"

## Usage in Services

```go
// Instead of hardcoded strings
return nil, status.Error(codes.InvalidArgument, "name is required")

// Use constants
return nil, status.Error(codes.InvalidArgument, ErrCompanyNameRequired)
```

## Benefits

1. **Consistency**: Same error messages across all endpoints
2. **Maintainability**: Change error message in one place
3. **Internationalization Ready**: Easy to add i18n support later
4. **Type Safety**: Compile-time checking of error message usage
5. **Documentation**: Clear overview of all possible errors

## gRPC Error Codes Used

- `codes.InvalidArgument` - Invalid input data (validation errors)
- `codes.NotFound` - Resource not found
- `codes.Internal` - Internal server errors (database errors, etc.)

## Client Error Handling Examples

### grpcurl
```bash
# InvalidArgument example
grpcurl -plaintext -d '{"name": ""}' localhost:50051 company.v1.CompanyService/CreateCompany
# ERROR: Code: InvalidArgument, Message: name is required

# NotFound example  
grpcurl -plaintext -d '{"id": "00000000-0000-0000-0000-000000000000"}' localhost:50051 company.v1.CompanyService/GetCompany
# ERROR: Code: NotFound, Message: company not found
```

### Go Client
```go
resp, err := client.GetCompany(ctx, &companyv1.GetCompanyRequest{Id: "invalid"})
if err != nil {
    if status.Code(err) == codes.NotFound {
        // Handle company not found
    } else if status.Code(err) == codes.InvalidArgument {
        // Handle validation error
    }
}
```

## Adding New Error Messages

1. Add constant to `internal/service/errors.go`:
```go
const (
    // ... existing errors
    ErrNewValidationError = "new validation error message"
)
```

2. Use in service methods:
```go
if invalidCondition {
    return nil, status.Error(codes.InvalidArgument, ErrNewValidationError)
}
```

## Future Enhancements

- **Structured Errors**: Add error codes for programmatic handling
- **Internationalization**: Support multiple languages
- **Error Details**: Add additional context using gRPC error details
- **Logging**: Consistent error logging with correlation IDs