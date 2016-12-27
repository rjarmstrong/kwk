package users

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/rpc/src/usersRpc"
	"google.golang.org/grpc"
	"time"
)

const (
	userDbKey = "user"
)

type Users struct {
	client   usersRpc.UsersRpcClient
	settings settings.ISettings
	headers  *rpc.Headers
}

func New(conn *grpc.ClientConn, s settings.ISettings, h *rpc.Headers) *Users {
	return &Users{client: usersRpc.NewUsersRpcClient(conn), settings: s, headers: h}
}

func (u *Users) SignIn(username string, password string) (*models.User, error) {
	if res, err := u.client.SignIn(u.headers.GetContext(), &usersRpc.SignInRequest{Username: username, Password: password}); err != nil {
		return nil, err
	} else {
		model := &models.User{}
		mapUser(res, model)
		return model, nil
	}
}

func (u *Users) SignUp(email string, username string, password string) (*models.User, error) {
	if res, err := u.client.SignUp(u.headers.GetContext(), &usersRpc.SignUpRequest{Username: username, Email: email, Password: password}); err != nil {
		return nil, err
	} else {
		model := &models.User{}
		mapUser(res, model)
		return model, nil
	}
}

func (u *Users) Get() (*models.User, error) {
	user := &models.User{}
	if err := u.settings.Get(userDbKey, user); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (u *Users) Signout() error {
	// Implement service call which would be more informational/analytical in nature
	return nil
}

func mapUser(rpc *usersRpc.UserResponse, model *models.User) {
	model.Id = rpc.Id
	model.Username = rpc.Username
	model.Email = rpc.Email
	model.Token = rpc.Token
	model.SnipCount = rpc.SnipCount
	model.RunCount = rpc.RunCount
	model.ClonedCount = rpc.ClonedCount
	model.Created = time.Unix(rpc.Created/1000, 0)
}
