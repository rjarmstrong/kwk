package models

import (
	"google.golang.org/grpc/codes"
	"encoding/json"
	"google.golang.org/grpc"
	"fmt"
)

type Code uint32

const (
	Code_Unspecified Code = 0
	// Snippets
	Code_NoTags                 Code = 3000
	Code_NewFullKeyEmpty        Code = 3010
	Code_FullKeyExistsWhenClone Code = 3020
	Code_MultipleSnippetsFound  Code = 3030

	// Users
	Code_WrongCreds      Code = 4001
	Code_UsernameTaken   Code = 4002
	Code_EmailTaken      Code = 4003
	Code_EmptyToken      Code = 4004
	Code_InvalidEmail    Code = 4005
	Code_InvalidUsername Code = 4006
	Code_InvalidPassword Code = 4107

	Code_MultiplePouches  Code = 4200
	Code_IncompleteAlias  Code = 4201
	Code_AliasMaxSegments Code = 4202
	Code_NoSnippetName    Code = 4203
	Code_PouchMaxSegments    Code = 4204
)

// ParseGrpcErr should be used at RPC service call level. i.e. the errors
// returned by the GRPC stubs need to be converted to local errors.
func ParseGrpcErr(e error) error {
	desc := grpc.ErrorDesc(e)
	m := &ClientErr{}
	m.TransportCode = grpc.Code(e)
	if err := json.Unmarshal([]byte(desc), m); err != nil {
		m.Title = desc
		return m
	}
	return m
}

func ErrOneLine(c Code, description string) error {
	return &ClientErr{TransportCode: codes.InvalidArgument, Msgs:[]Msg{{Code:c, Desc:description}}}
}

type ClientErr struct {
	TransportCode codes.Code
	Msgs          []Msg
	Title         string
}

func (e ClientErr) Error() string {
	return fmt.Sprintf("%d %s\n%v+", e.TransportCode, e.Title, e.Msgs)
}

type Msg struct {
	Code Code
	Desc string
}
