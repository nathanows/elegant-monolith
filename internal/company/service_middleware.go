package company

import (
	"context"

	"github.com/go-kit/kit/log"
	pb "github.com/nathanows/elegant-monolith/_protos/companyusers"
)

// ServiceMiddleware describes a service middleware
type ServiceMiddleware func(Service) Service

// ServiceLoggingMiddleware takes a logger as a dependency and returns a service middleware
func ServiceLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return serviceLoggingMiddleware{logger, next}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw serviceLoggingMiddleware) Save(ctx context.Context, company *pb.Company) (returned *pb.Company, err error) {
	defer func() {
		if err == nil {
			mw.logger.Log("method", "Save", "id", returned.ID)
		} else {
			mw.logger.Log("method", "Save", "err", err.Error())
		}
	}()
	return mw.next.Save(ctx, company)
}

func (mw serviceLoggingMiddleware) Find(ctx context.Context, id int64) (*pb.Company, error) {
	defer func() {
		mw.logger.Log("method", "Find", "id", id)
	}()
	return mw.next.Find(ctx, id)
}

func (mw serviceLoggingMiddleware) Delete(ctx context.Context, id int64) error {
	defer func() {
		mw.logger.Log("method", "Delete", "id", id)
	}()
	return mw.next.Delete(ctx, id)
}

func (mw serviceLoggingMiddleware) FindAll(ctx context.Context) (returned []*pb.Company, err error) {
	defer func() {
		mw.logger.Log("method", "FindAll", "results_returned", len(returned))
	}()
	return mw.next.FindAll(ctx)
}
