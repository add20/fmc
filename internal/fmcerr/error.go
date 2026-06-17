package fmcerr

import "fmt"

type ErrorCode string

const (
	ErrDuplicateSlug    ErrorCode = "DUPLICATE_SLUG"
	ErrFrontMatterParse ErrorCode = "FRONTMATTER_PARSE"
	ErrConfigLoad       ErrorCode = "CONFIG_LOAD"
	ErrReadFile         ErrorCode = "READ_FILE"
	ErrWriteFile        ErrorCode = "WRITE_FILE"
)

type FMCError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *FMCError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *FMCError) Unwrap() error {
	return e.Cause
}
