package errs

import (
	"fmt"
)

type ErrCode uint32

const (
	CodeUnspecified        ErrCode = 0
	CodeNotFound           ErrCode = 10
	CodeInvalidArgument    ErrCode = 20
	CodeUnexpectedEndOfTar ErrCode = 30
	CodeAlreadyExists      ErrCode = 40
	CodeNotImplemented     ErrCode = 50
	CodeInternalError      ErrCode = 60
	CodeNotAvailable       ErrCode = 70
	CodeNotAuthenticated   ErrCode = 80
	CodePermissionDenied   ErrCode = 90

	// Snippets
	CodeSnippetNameRequired          ErrCode = 3020
	CodeMultipleSnippetsFound        ErrCode = 3030
	CodeTwoArgumentsRequiredForMove  ErrCode = 3040
	CodeSnippetNotVerified           ErrCode = 3050
	CodeSnippetVulnerable            ErrCode = 3060
	CodeSnippetTooLong               ErrCode = 3070
	CodeInvalidSnipName              ErrCode = 3080
	CodeInvalidExtension             ErrCode = 3090
	CodeDescription                  ErrCode = 3100
	CodePreviewInvalid               ErrCode = 3110
	CodeSnippetNameExists            ErrCode = 3120
	CodeSnippetNotYetConsistent      ErrCode = 3130
	CodeBothNameAndExtensionRequired ErrCode = 3150
	CodeNoSnipNamesProvided          ErrCode = 3160
	CodeNoSnipNamesFound             ErrCode = 3170
	CodeListCategoryNotFound         ErrCode = 3180

	CodeEmbeddedSnippetCantCallItself   ErrCode = 3190
	CodeEmbeddedSnippetsRequireExt      ErrCode = 3200
	CodeEmbeddedSnippetsRequireUsername ErrCode = 3210
	CodeNoTags                          ErrCode = 3250
	CodeClonedFromAliasVersionMissing   ErrCode = 3260

	// Pouches
	CodeInvalidPouchName      ErrCode = 3500
	CodePouchExists           ErrCode = 3510
	CodePouchNotYetConsistent ErrCode = 3520
	CodePouchNotExists        ErrCode = 3530
	CodePouchFull             ErrCode = 3540
	CodeCantRenameRootPouch   ErrCode = 3550
	CodeCantDeletedRootPouch  ErrCode = 3560

	// Users
	CodeWrongCreds        ErrCode = 4010
	CodeUsernameTaken     ErrCode = 4020
	CodeEmailTaken        ErrCode = 4030
	CodeEmptyToken        ErrCode = 4040
	CodeInvalidEmail      ErrCode = 4050
	CodeInvalidUsername   ErrCode = 4060
	CodeInvalidPassword   ErrCode = 4170
	CodeInvalidInviteCode ErrCode = 4180

	CodeMultiplePouches  ErrCode = 4210
	CodeIncompleteAlias  ErrCode = 4220
	CodeURITooManySegs   ErrCode = 4230
	CodeNoSnippetName    ErrCode = 4240
	CodePouchMaxSegments ErrCode = 4250

	CodeInvalidConfigSection    ErrCode = 6010
	CodeEnvironmentNotSupported ErrCode = 6020

	//Runners
	CodeRunnerExit ErrCode = 700
	CodeErrTooDeep ErrCode = 710

	//Files
	CodeFileNotFound ErrCode = 800

	CodePrinterNotFound ErrCode = 900
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
	FileNotFound     = New(CodeFileNotFound, "File not found")
	FileExpired      = New(CodeFileNotFound, "File found but expired")
	NoAuthToken      = New(CodeEmptyToken, "No auth token supplied.")
	// kwk level
	SnipNameInvalid = New(CodeInvalidSnipName,
		"Snippet names should contain: alphanumeric, '.', '-' or '_' and be between 3 and 50 characters.")
	ExtensionInvalid = New(CodeInvalidSnipName,
		"Extensions should have only alphanumeric characters and be at most 12 characters.")
	PouchNameInvalid = New(CodeInvalidPouchName,
		"Pouch names should contain: alphanumeric, '-' or '_' and be be between 1 and 30 characters.")
	UsernameInvalid = New(CodeInvalidUsername,
		"Username should contain: alphanumeric, '-' or '_' and be be between 3 and 15 characters.")
	SnippetNameRequired       = New(CodeSnippetNameRequired, "Snippet name required")
	SnippetNotVerified        = New(CodeSnippetNotVerified, "The checksum doesn't match the snippet.")
	SnipNameExists            = New(CodeAlreadyExists, "That name/ext combo already exists.")
	TwoArgumentsReqForMove    = New(CodeTwoArgumentsRequiredForMove, "Two arguments are required for the move command.")
	EmbeddedRequiresExtension = New(CodeEmbeddedSnippetsRequireExt, "Embedded snippets require an ext.")
	EmbeddedRequiresUsername  = New(CodeEmbeddedSnippetsRequireUsername, "Embedded snippets require username.")

	AliasTooManySegments  = New(CodeURITooManySegs, "A snippet alias can only consist of up to 3 segments.")
	MultipleTargetPouches = New(CodeMultiplePouches, "Cannot move or rename snippets from multiple source pouches.")
)

func New(c ErrCode, format string, args ...interface{}) error {
	return &Error{Message: fmt.Sprintf(format, args...), Code: c}
}

type Error struct {
	Message string  `json:"message"`
	Code    ErrCode `json:"code"`
}

func HasCode(err error, code ErrCode) bool {
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
