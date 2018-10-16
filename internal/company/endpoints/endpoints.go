package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/types"

	pb "github.com/nathanows/elegant-monolith/_protos/companyusers"
	"github.com/nathanows/elegant-monolith/internal/company/service"
)

// Set collects all of the endpoints that compose a user service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	SaveEndpoint    endpoint.Endpoint
	FindEndpoint    endpoint.Endpoint
	DeleteEndpoint  endpoint.Endpoint
	FindAllEndpoint endpoint.Endpoint
}

// NewEndpointSet returns a constructed Set for use to instantiate server
func NewEndpointSet(svc service.Service, logger log.Logger) Set {
	var saveEndpoint endpoint.Endpoint
	{
		saveEndpoint = MakeSaveEndpoint(svc)
		saveEndpoint = LoggingMiddleware(log.With(logger, "method", "Save"))(saveEndpoint)
	}
	var findEndpoint endpoint.Endpoint
	{
		findEndpoint = MakeFindEndpoint(svc)
		findEndpoint = LoggingMiddleware(log.With(logger, "method", "Find"))(findEndpoint)
	}
	var deleteEndpoint endpoint.Endpoint
	{
		deleteEndpoint = MakeDeleteEndpoint(svc)
		deleteEndpoint = LoggingMiddleware(log.With(logger, "method", "Delete"))(deleteEndpoint)
	}
	var findAllEndpoint endpoint.Endpoint
	{
		findAllEndpoint = MakeFindAllEndpoint(svc)
		findAllEndpoint = LoggingMiddleware(log.With(logger, "method", "FindAll"))(findAllEndpoint)
	}
	return Set{
		SaveEndpoint:    saveEndpoint,
		FindEndpoint:    findEndpoint,
		DeleteEndpoint:  deleteEndpoint,
		FindAllEndpoint: findAllEndpoint,
	}
}

// MakeSaveEndpoint constructs a Save endpoint wrapping the service.
func MakeSaveEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.SaveCompanyRequest)
		company, err := s.Save(ctx, req.Company)
		if err != nil {
			return nil, err
		}
		return company, nil
	}
}

// MakeFindEndpoint constructs a Find endpoint wrapping the service.
func MakeFindEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.FindCompanyRequest)
		company, err := s.Find(ctx, req.ID)
		return company, nil
	}
}

// MakeDeleteEndpoint constructs a Delete endpoint wrapping the service.
func MakeDeleteEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.DeleteCompanyRequest)
		err = s.Delete(ctx, req.ID)
		return &types.Empty{}, nil
	}
}

// MakeFindAllEndpoint constructs a FindAll endpoint wrapping the service.
func MakeFindAllEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		companies, err := s.FindAll(ctx)
		return &pb.FindAllCompaniesResponse{Companies: companies}, nil
	}
}
