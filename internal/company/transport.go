package company

import (
	"context"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	types "github.com/gogo/protobuf/types"
	oldcontext "golang.org/x/net/context"

	pb "github.com/nathanows/elegant-monolith/_protos"
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
		save: grpctransport.NewServer(
			endpoints.SaveEndpoint,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
		find: grpctransport.NewServer(
			endpoints.FindEndpoint,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
		delete: grpctransport.NewServer(
			endpoints.DeleteEndpoint,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
		findAll: grpctransport.NewServer(
			endpoints.FindAllEndpoint,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
	}
}

func (s *grpcServer) Save(ctx oldcontext.Context, req *pb.SaveCompanyRequest) (*pb.Company, error) {
	_, rep, err := s.save.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
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

func decodeGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
}

func encodeGRPCResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	return grpcResp, nil
}
