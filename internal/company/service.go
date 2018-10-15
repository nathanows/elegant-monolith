package company

import (
	"context"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/types"
	pb "github.com/nathanows/elegant-monolith/_protos"
)

// Service interface defines the core Company service functionality
type Service interface {
	Save(ctx context.Context, company *pb.Company) (*pb.Company, error)
	Find(ctx context.Context, id int64) (*pb.Company, error)
	FindAll(ctx context.Context) ([]*pb.Company, error)
	Delete(ctx context.Context, id int64) error
}

// NewService returns an initialized Service wired up with all middleware
func NewService(logger log.Logger, repository Repository) Service {
	var svc Service
	{
		svc = NewBasicService(repository)
		svc = ServiceLoggingMiddleware(logger)(svc)
	}
	return svc
}

// NewBasicService returns an initialized Service without middleware
func NewBasicService(repository Repository) Service {
	return basicService{
		repository: repository,
	}
}

type basicService struct {
	repository Repository
}

func (s basicService) Save(ctx context.Context, companyToSave *pb.Company) (*pb.Company, error) {
	companyDTO := toDTO(companyToSave)

	ok, err := validate(companyDTO)
	if !ok {
		return nil, err
	}

	saved, err := s.repository.save(companyDTO)
	if err != nil {
		return nil, err
	}

	return saved.toProto(), nil
}

func (s basicService) Find(ctx context.Context, id int64) (*pb.Company, error) {
	return &pb.Company{ID: 1}, nil
}

func (s basicService) Delete(ctx context.Context, id int64) error {
	return nil
}

func (s basicService) FindAll(ctx context.Context) ([]*pb.Company, error) {
	return nil, nil
}

func validate(company *companyDTO) (bool, error) {
	return govalidator.ValidateStruct(company)
}

func init() {
	govalidator.TagMap["notduck"] = govalidator.Validator(verifyNotDuck)
}

func verifyNotDuck(str string) bool {
	return strings.ToLower(str) != "duck"
}

type companyDTO struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name" valid:"notduck~No ducks allowed,required"` // example custom validator
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (company *companyDTO) toProto() *pb.Company {
	return &pb.Company{
		ID:        company.ID,
		Name:      company.Name,
		CreatedAt: genPbTimestamp(company.CreatedAt),
		UpdatedAt: genPbTimestamp(company.UpdatedAt),
	}
}

func toDTO(company *pb.Company) *companyDTO {
	return &companyDTO{
		ID:        company.ID,
		Name:      company.Name,
		CreatedAt: genDTOTimestamp(company.CreatedAt),
		UpdatedAt: genDTOTimestamp(company.UpdatedAt),
	}
}

func genDTOTimestamp(pbTime *types.Timestamp) time.Time {
	ts, err := types.TimestampFromProto(pbTime)
	if err != nil {
		return time.Unix(0, 0).UTC()
	}
	return ts
}

func genPbTimestamp(time time.Time) *types.Timestamp {
	ts, err := types.TimestampProto(time)
	if err != nil {
		return &types.Timestamp{Seconds: 0, Nanos: 0}
	}
	return ts
}
