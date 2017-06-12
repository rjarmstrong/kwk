package handlers

import (
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsers_SignUp(t *testing.T) {
	dlg.returnsFor["FormField"] = response{val: &out.DialogResponse{Value: "some string"}}
	usersClient.returnsFor["SignUp"] = response{val: &types.SignUpResponse{User: &types.User{}}}
	err := users.SignUp()
	assert.Nil(t, err)
	funcName := writer.PopCalled("EWrite")
	assert.Equal(t, "UserSignedUp", funcName)
}

func TestUsers_SignIn(t *testing.T) {
	dlg.returnsFor["FormField"] = response{val: &out.DialogResponse{Value: "some string"}}
	usersClient.returnsFor["SignIn"] = response{val: &types.SignInResponse{
		User:        &types.User{},
		AccessToken: "asdfsf",
		Root:        johnnyRoot.val.(*types.RootResponse),
	}}
	err := users.SignIn("username1")
	assert.Nil(t, err)
	assert.Equal(t, "johnny", rootPrintCalled.Username)
}

func TestUsers_SignOut(t *testing.T) {
	err := users.SignOut()
	assert.Nil(t, err)
	funcName := writer.PopCalled("EWrite")
	assert.Equal(t, "UserSignedOut", funcName)
}

func TestUsers_ChangePassword(t *testing.T) {
	dlg.returnsFor["FormField"] = response{val: &out.DialogResponse{Value: "some string"}}
	err := users.ChangePassword()
	assert.Nil(t, err)
	funcName := writer.PopCalled("EWrite")
	assert.Equal(t, "UserPasswordChanged", funcName)
}

func TestUsers_ForgotPassword(t *testing.T) {
	dlg.returnsFor["FormField"] = response{val: &out.DialogResponse{Value: "some string"}}
	err := users.ForgotPassword()
	assert.Nil(t, err)
	funcName := writer.PopCalled("EWrite")
	assert.Equal(t, "UserPasswordResetSent", funcName)
}

func TestUsers_Profile(t *testing.T) {
	t.Log("AUTHENTICATED")
	err := users.Profile()
	assert.Nil(t, err)
	funcName := writer.PopCalled("EWrite")
	assert.Equal(t, "UserProfile", funcName)

	t.Log("NOT AUTHENTICATED")
	pr.AccessToken = ""
	err = users.Profile()
	assert.Equal(t, errs.NotAuthenticated, err)
}
