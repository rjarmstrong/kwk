package account

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/rpc/src/usersRpc"
	"google.golang.org/grpc"
	"time"
)

type StdManager struct {
	client   usersRpc.UsersRpcClient
	settings config.Persister
	headers  *rpc.Headers
}

func NewStdManager(conn *grpc.ClientConn, s config.Persister, h *rpc.Headers) *StdManager {
	return &StdManager{client: usersRpc.NewUsersRpcClient(conn), settings: s, headers: h}
}

func (u *StdManager) SignIn(username string, password string) (*models.User, error) {
	if res, err := u.client.SignIn(u.headers.Context(), &usersRpc.SignInRequest{Username: username, Password: password}); err != nil {
		return nil, err
	} else {
		model := &models.User{}
		mapUser(res, model)
		return model, nil
	}
}

func (u *StdManager) SignUp(email string, username string, password string) (*models.User, error) {
	if res, err := u.client.SignUp(u.headers.Context(), &usersRpc.SignUpRequest{Username: username, Email: email, Password: password}); err != nil {
		return nil, err
	} else {
		model := &models.User{}
		mapUser(res, model)
		return model, nil
	}
}

var user *models.User

func (u *StdManager) Get() (*models.User, error) {
	if user != nil && user.Id != "" {
		return user, nil
	}
	user = &models.User{}
	if err := u.settings.Get(models.ProfileFullKey, user, 0); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (u *StdManager) HasValidCredentials() bool {
	if user, err := u.Get(); err != nil {
		return false
	} else if user != nil {
		// TODO: Check jwt expiry
		return true
	}
	return false
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
