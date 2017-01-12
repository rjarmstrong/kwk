package account

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/rpc/src/usersRpc"
	"google.golang.org/grpc"
	"time"
)

const (
	userDbKey = "user"
)

type StdManager struct {
	client   usersRpc.UsersRpcClient
	settings config.Settings
	headers  *rpc.Headers
}

func NewStdManager(conn *grpc.ClientConn, s config.Settings, h *rpc.Headers) *StdManager {
	return &StdManager{client: usersRpc.NewUsersRpcClient(conn), settings: s, headers: h}
}

func (u *StdManager) SignIn(username string, password string) (*models.User, error) {
	if res, err := u.client.SignIn(u.headers.GetContext(), &usersRpc.SignInRequest{Username: username, Password: password}); err != nil {
		return nil, err
	} else {
		model := &models.User{}
		mapUser(res, model)
		return model, nil
	}
}

func (u *StdManager) SignUp(email string, username string, password string) (*models.User, error) {
	if res, err := u.client.SignUp(u.headers.GetContext(), &usersRpc.SignUpRequest{Username: username, Email: email, Password: password}); err != nil {
		return nil, err
	} else {
		model := &models.User{}
		mapUser(res, model)
		return model, nil
	}
}

func (u *StdManager) Get() (*models.User, error) {
	user := &models.User{}
	if err := u.settings.Get(userDbKey, user, 0); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (u *StdManager) Signout() error {
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
