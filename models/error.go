package models

import (
	"google.golang.org/grpc"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
)

type Code uint32

const (
	Code_Unspecified Code = 0
	Code_NotFound Code = 10
	Code_InvalidArgument Code = 20
	Code_UnexpectedEndOfTar Code = 30
	// Snippets
	Code_NoTags                 Code = 3000
	Code_NewFullKeyEmpty        Code = 3010
	Code_FullKeyExistsWhenClone Code = 3020
	Code_MultipleSnippetsFound  Code = 3030
	Code_TwoArgumentsRequiredForMove Code = 3040
	Code_SnippetNotVerified Code = 3050

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
	Code_EnvironmentNotSupported Code = 6020
)

// ParseGrpcErr should be used at RPC service call level. i.e. the errors
// returned by the GRPC stubs need to be converted to local errors.
func ParseGrpcErr(e error) error {
	desc := grpc.ErrorDesc(e)
	m := &ClientErr{}
	m.remoteCode = grpc.Code(e)
	if err := json.Unmarshal([]byte(desc), m); err != nil {
		m.Title = desc
		return m
	}
	return m
}

func ErrOneLine(c Code, desc string, args ...interface{}) error {
	return &ClientErr{Msgs: []Msg{{Code: c, Desc: fmt.Sprintf(desc, args...)}}}
}

type ClientErr struct {
	Msgs  []Msg
	Title string
	remoteCode codes.Code
}

func (c *ClientErr) Contains(code Code) bool{
	for _, v := range c.Msgs {
		if v.Code == code {
			return true
		}
	}
	return false
}

func (c ClientErr) Error() string {
	return fmt.Sprintf("%s %+v %d", c.Title, c.Msgs, c.remoteCode)
}

type Msg struct {
	Code Code
	Desc string
}
