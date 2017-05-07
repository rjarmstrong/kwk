package gokwk

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/rpc/src/usersRpc"
	"google.golang.org/grpc"
	"time"
	"bitbucket.com/sharingmachine/kwkcli/src/persist"
	"bitbucket.com/sharingmachine/types"
)

type UsersGrpc struct {
	client   usersRpc.UsersRpcClient
	settings persist.Persister
	headers  Headers
}

func NewUsers(conn *grpc.ClientConn, s persist.Persister, client types.AppInfo) Users {
	return &UsersGrpc{client: usersRpc.NewUsersRpcClient(conn), settings: s, headers: Headers{version:client.Version}}
}

func (u *UsersGrpc) SignIn(username string, password string) (*models.User, error) {
	if res, err := u.client.SignIn(u.headers.Context(), &usersRpc.SignInRequest{Username: username, Password: password}); err != nil {
		return nil, err
	} else {
		model := &models.User{}
		mapUser(res, model)
		return model, nil
	}
}

func (u *UsersGrpc) SignUp(email string, username string, password string, inviteCode string) (*models.User, error) {
	if res, err := u.client.SignUp(u.headers.Context(), &usersRpc.SignUpRequest{Username: username, Email: email, Password: password, InviteCode: inviteCode}); err != nil {
		return nil, err
	} else {
		model := &models.User{}
		mapUser(res, model)
		return model, nil
	}
}

func (u *UsersGrpc) Get() (*models.User, error) {
	if models.Principal != nil && models.Principal.Id != "" {
		return models.Principal, nil
	}
	models.Principal = &models.User{}
	if err := u.settings.Get(models.ProfileFullKey, models.Principal, 0); err != nil {
		return nil, err
	} else {
		return models.Principal, nil
	}
}

func (u *UsersGrpc) ResetPassword(email string) (bool, error) {
	req := &usersRpc.ResetRequest{Email: email}
	_, err := u.client.ResetPassword(u.headers.Context(), req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UsersGrpc) ChangePassword(p models.ChangePasswordParams) (bool, error) {
	req := &usersRpc.ChangeRequest{
		Email:            p.Email,            //Required if no email
		Username:         p.Username,         //Required if no email
		ExistingPassword: p.ExistingPassword, // Required if no ResetToken
		NewPassword:      p.NewPassword,      // Required
		ResetToken:       p.ResetToken,       // Required if no Existing password
	}
	_, err := u.client.ChangePassword(u.headers.Context(), req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UsersGrpc) HasValidCredentials() bool {
	usr, err := u.Get()
	if err != nil {
		return false
	}
	if usr != nil {
		// TODO: Check jwt expiry
		return true
	}
	return false
}

func (u *UsersGrpc) Signout() error {
	err := u.settings.DeleteAll()
	if err != nil {
		return err
	}
	// Implement service call which would be more informational/analytical in nature
	return nil
}

func mapUser(rpc *usersRpc.User, model *models.User) {
	model.Id = rpc.Id
	model.Username = rpc.Username
	model.Email = rpc.Email
	model.Token = rpc.Token
	model.SnipCount = rpc.SnipCount
	model.RunCount = rpc.RunCount
	model.ClonedCount = rpc.ClonedCount
	model.Created = time.Unix(rpc.Created/1000, 0)
}