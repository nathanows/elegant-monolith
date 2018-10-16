package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	pb "github.com/nathanows/elegant-monolith/_protos/companyusers"
)

// NewHTTPServer mounts all of the service endpoints into an http.Handler.
func NewHTTPServer(endpoints Set, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeHTTPError),
	}

	r.Methods("POST").Path("/save").Handler(httptransport.NewServer(
		endpoints.SaveEndpoint,
		decodeHTTPSaveRequest,
		encodeHTTPGenericResponse,
		options...,
	))

	return r
}

func decodeHTTPSaveRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req *pb.SaveCompanyRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeHTTPError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
