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
	Code_TwoArgumentsRequiredForMove Code = 3040

	// Users
	Code_WrongCreds      Code = 4010
	Code_UsernameTaken   Code = 4020
	Code_EmailTaken      Code = 4030
	Code_EmptyToken      Code = 4040
	Code_InvalidEmail    Code = 4050
	Code_InvalidUsername Code = 4060
	Code_InvalidPassword Code = 4170

	Code_MultiplePouches  Code = 4210
	Code_IncompleteAlias  Code = 4220
	Code_AliasMaxSegments Code = 4230
	Code_NoSnippetName    Code = 4240
	Code_PouchMaxSegments    Code = 4250

	//Network
	Code_ApiDown Code = 5010

	Code_InvalidConfigSection Code = 6010
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
