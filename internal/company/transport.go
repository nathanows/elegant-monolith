package company

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	types "github.com/gogo/protobuf/types"
	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/nathanows/elegant-monolith/_protos/companyusers"
)

type grpcServer struct {
	save    grpctransport.Handler
	find    grpctransport.Handler
	delete  grpctransport.Handler
	findAll grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints Set, logger log.Logger) pb.CompanySvcServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		save:    newGPRCServer(endpoints.SaveEndpoint, options...),
		find:    newGPRCServer(endpoints.FindEndpoint, options...),
		delete:  newGPRCServer(endpoints.DeleteEndpoint, options...),
		findAll: newGPRCServer(endpoints.FindAllEndpoint, options...),
	}
}

func (s *grpcServer) Save(ctx oldcontext.Context, req *pb.SaveCompanyRequest) (*pb.Company, error) {
	_, rep, err := s.save.ServeGRPC(ctx, req)
	if err != nil {
		encodedErr := encodeError(err)
		return nil, encodedErr
	}
	return rep.(*pb.Company), nil
}

func (s *grpcServer) Find(ctx oldcontext.Context, req *pb.FindCompanyRequest) (*pb.Company, error) {
	_, rep, err := s.find.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.Company), nil
}

func (s *grpcServer) Delete(ctx oldcontext.Context, req *pb.DeleteCompanyRequest) (*types.Empty, error) {
	_, rep, err := s.delete.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*types.Empty), nil
}

func (s *grpcServer) FindAll(ctx oldcontext.Context, req *pb.FindAllCompaniesRequest) (*pb.FindAllCompaniesResponse, error) {
	_, rep, err := s.findAll.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FindAllCompaniesResponse), nil
}

func encodeError(err error) error {
	switch err {
	case ErrRepository:
		return status.Error(codes.Internal, err.Error())
	case ErrRequireName, ErrInvalidName, ErrCompanyNotFound:
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}

}

func newGPRCServer(endpoint endpoint.Endpoint, options ...grpctransport.ServerOption) *grpctransport.Server {
	return grpctransport.NewServer(
		endpoint,
		decodeGRPCRequest,
		encodeGRPCResponse,
		options...,
	)
}

func decodeGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
}

func encodeGRPCResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	return grpcResp, nil
}
