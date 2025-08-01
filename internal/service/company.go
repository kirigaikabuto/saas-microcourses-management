package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kirigaikabuto/saas-microcourses-management/internal/helper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kirigaikabuto/saas-microcourses-management/internal/db"
	companyv1 "github.com/kirigaikabuto/saas-microcourses-management/proto/gen/proto/company/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CompanyService struct {
	companyv1.UnimplementedCompanyServiceServer
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewCompanyService(pool *pgxpool.Pool) *CompanyService {
	return &CompanyService{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (s *CompanyService) CreateCompany(ctx context.Context, req *companyv1.CreateCompanyRequest) (*companyv1.CreateCompanyResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, ErrCompanyNameRequired)
	}
	if req.SubscriptionPlan == "" {
		return nil, status.Error(codes.InvalidArgument, ErrCompanySubscriptionRequired)
	}

	company, err := s.queries.CreateCompany(ctx, db.CreateCompanyParams{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Name:             pgtype.Text{String: req.Name, Valid: true},
		SubscriptionPlan: pgtype.Text{String: req.SubscriptionPlan, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, ErrFailedToCreateCompany)
	}

	return &companyv1.CreateCompanyResponse{
		Company: helper.FromDbCompanyToProto(company),
	}, nil
}

func (s *CompanyService) GetCompany(ctx context.Context, req *companyv1.GetCompanyRequest) (*companyv1.GetCompanyResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, ErrCompanyIDRequired)
	}

	// Parse UUID
	uuidParsed, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidUUIDFormat)
	}

	var uuidParam pgtype.UUID
	uuidParam.Bytes = uuidParsed
	uuidParam.Valid = true

	company, err := s.queries.GetCompany(ctx, uuidParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, ErrCompanyNotFound)
		}
		return nil, status.Error(codes.Internal, ErrFailedToGetCompany)
	}

	return &companyv1.GetCompanyResponse{
		Company: helper.FromDbCompanyToProto(company),
	}, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, req *companyv1.UpdateCompanyRequest) (*companyv1.UpdateCompanyResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, ErrCompanyIDRequired)
	}

	// Parse UUID
	uuidParsed, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidUUIDFormat)
	}

	arg := db.UpdateCompanyParams{
		ID: pgtype.UUID{
			Bytes: uuidParsed,
			Valid: true,
		},
	}

	if req.Name != nil {
		arg.Name = pgtype.Text{String: *req.Name, Valid: true}
	}

	if req.SubscriptionPlan != nil {
		arg.SubscriptionPlan = pgtype.Text{String: *req.SubscriptionPlan, Valid: true}
	}

	company, err := s.queries.UpdateCompany(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, ErrCompanyNotFound)
		}
		return nil, status.Error(codes.Internal, ErrFailedToUpdateCompany)
	}

	return &companyv1.UpdateCompanyResponse{
		Company: helper.FromDbCompanyToProto(company),
	}, nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, req *companyv1.DeleteCompanyRequest) (*companyv1.DeleteCompanyResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, ErrCompanyIDRequired)
	}

	// Parse UUID
	uuidParsed, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidUUIDFormat)
	}

	var uuidParam pgtype.UUID
	uuidParam.Bytes = uuidParsed
	uuidParam.Valid = true

	err = s.queries.DeleteCompany(ctx, uuidParam)
	if err != nil {
		return nil, status.Error(codes.Internal, ErrFailedToDeleteCompany)
	}

	return &companyv1.DeleteCompanyResponse{
		Success: true,
	}, nil
}

func (s *CompanyService) ListCompanies(ctx context.Context, req *companyv1.ListCompaniesRequest) (*companyv1.ListCompaniesResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	companies, err := s.queries.ListCompanies(ctx, db.ListCompaniesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, ErrFailedToListCompanies)
	}

	total, err := s.queries.CountCompanies(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, ErrFailedToCountCompanies)
	}

	pbCompanies := make([]*companyv1.Company, len(companies))
	for i, company := range companies {
		pbCompanies[i] = helper.FromDbCompanyToProto(company)
	}

	return &companyv1.ListCompaniesResponse{
		Companies: pbCompanies,
		Total:     int32(total),
	}, nil
}
