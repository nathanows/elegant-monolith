package company

import "errors"

// Company Service Error descriptions
const (
	ErrorInvalidName     = "invalid company name"
	ErrorRequireName     = "missing required name"
	ErrorUniqueName      = "company with name already exists"
	ErrorCompanyNotFound = "company not found"
	ErrorRepository      = "unable to query repository"
)

// Device Service Errors
var (
	ErrInvalidName     = errors.New(ErrorInvalidName)
	ErrRequireName     = errors.New(ErrorRequireName)
	ErrUniqueName      = errors.New(ErrorUniqueName)
	ErrCompanyNotFound = errors.New(ErrorCompanyNotFound)
	ErrRepository      = errors.New(ErrorRepository)
)
