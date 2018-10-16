package company

import "errors"

// Company Service Error descriptions
const (
	ErrorInvalidName     = "invalid company name"
	ErrorRequireName     = "missing required name"
	ErrorCompanyNotFound = "company not found"
	ErrorRepository      = "unable to query repository"
)

// Device Service Errors
var (
	ErrInvalidName     = errors.New(ErrorInvalidName)
	ErrRequireName     = errors.New(ErrorRequireName)
	ErrCompanyNotFound = errors.New(ErrorCompanyNotFound)
	ErrRepository      = errors.New(ErrorRepository)
)
