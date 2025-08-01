package service

// Error messages for company service
const (
	// Validation errors
	ErrCompanyNameRequired         = "name is required"
	ErrCompanySubscriptionRequired = "subscription_plan is required"
	ErrCompanyIDRequired          = "id is required"
	ErrInvalidUUIDFormat          = "invalid UUID format"

	// Operation errors
	ErrCompanyNotFound      = "company not found"
	ErrFailedToCreateCompany = "failed to create company"
	ErrFailedToGetCompany   = "failed to get company"
	ErrFailedToUpdateCompany = "failed to update company"
	ErrFailedToDeleteCompany = "failed to delete company"
	ErrFailedToListCompanies = "failed to list companies"
	ErrFailedToCountCompanies = "failed to count companies"
)