package errs

import (
	"fmt"
)

type ErrCode uint32

const (
	CodeUnspecified      ErrCode = 0
	CodeNotFound         ErrCode = 10
	CodeInvalidArgument  ErrCode = 20
	CodeAlreadyExists    ErrCode = 40
	CodeNotImplemented   ErrCode = 50
	CodeInternalError    ErrCode = 60
	CodeNotAvailable     ErrCode = 70
	CodeNotAuthenticated ErrCode = 80
	CodePermissionDenied ErrCode = 90
)

var (
	// High level
	NotFound         = New(CodeNotFound, "Not found")
	NotAuthenticated = New(CodeNotAuthenticated, "Not authenticated")
	PermissionDenied = New(CodePermissionDenied, "Permission denied")
	AlreadyExists    = New(CodeAlreadyExists, "Already exists")
	NotImplemented   = New(CodeNotImplemented, "Not implemented")
	Internal         = New(CodeInternalError, "Internal error")
	ApiDown          = New(CodeNotAvailable, "The kwk api is down, please try again.")
	FileNotFound     = New(CodeNotFound, "File not found")
	FileExpired      = New(CodeNotFound, "File found but expired")
	NoAuthToken      = New(CodeInvalidArgument, "No auth token supplied.")
	// kwk level
	SnipNameInvalid = New(CodeInvalidArgument,
		"Snippet names should contain: alphanumeric, '.', '-' or '_' and be between 1 and 50 characters.")
	ExtensionInvalid = New(CodeInvalidArgument,
		"Extensions should have only alphanumeric characters and be at most 12 characters.")
	PouchNameInvalid = New(CodeInvalidArgument,
		"Pouch names should contain: alphanumeric, '-' or '_' and be be between 1 and 30 characters.")
	UsernameInvalid = New(CodeInvalidArgument,
		"Username should contain: alphanumeric, '-' or '_' and be be between 3 and 15 characters.")
	SnippetNameRequired       = New(CodeInvalidArgument, "Snippet name required")
	SnippetNotVerified        = New(CodeInvalidArgument, "The checksum doesn't match the snippet. Please `kwk edit <uri>` to fix.")
	SnipNameExists            = New(CodeAlreadyExists, "That name/ext combo already exists.")
	TwoArgumentsReqForMove    = New(CodeInvalidArgument, "Two arguments are required for the move command.")
	EmbeddedRequiresExtension = New(CodeInvalidArgument, "Embedded snippets require an ext.")
	EmbeddedRequiresUsername  = New(CodeInvalidArgument, "Embedded snippets require username.")

	AliasTooManySegments  = New(CodeInvalidArgument, "A snippet alias can only consist of up to 3 segments.")
	MultipleTargetPouches = New(CodeInvalidArgument, "Cannot move or rename snippets from multiple source pouches.")
)

func New(c ErrCode, format string, args ...interface{}) error {
	return &Error{Message: fmt.Sprintf(format, args...), Code: c}
}

type Error struct {
	Message string  `json:"message"`
	Code    ErrCode `json:"code"`
}

func HasCode(err error, code ErrCode) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == code
}

func (c Error) Error() string {
	return fmt.Sprintf("%s", c.Message)
}

type Handler interface {
	Handle(err error)
}

type HandlerFunc func(err error)

func (p HandlerFunc) Handle(err error) {
	p(err)
}
