package helper

import (
	"github.com/google/uuid"
	"github.com/kirigaikabuto/saas-microcourses-management/internal/db"
	companyv1 "github.com/kirigaikabuto/saas-microcourses-management/proto/gen/proto/company/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromDbCompanyToProto(company db.Company) *companyv1.Company {
	res := &companyv1.Company{
		Id:        uuid.UUID(company.ID.Bytes).String(),
		CreatedAt: timestamppb.New(company.CreatedAt.Time),
		UpdatedAt: timestamppb.New(company.UpdatedAt.Time),
	}

	if company.Name.Valid {
		res.Name = company.Name.String
	}

	if company.SubscriptionPlan.Valid {
		res.SubscriptionPlan = company.SubscriptionPlan.String
	}

	return res
}
