package users

import (
	"bitbucket.com/sharingmachine/kwkweb/rpc/usersRpc"
	"google.golang.org/grpc"
	"github.com/kwk-links/kwk-cli/libs/rpc"
	"github.com/kwk-links/kwk-cli/libs/models"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
)

const (
	userDbKey = "user"
)

type Users struct {
	client usersRpc.UsersRpcClient
	settings settings.ISettings
	rpc.Headers
}

func New(conn *grpc.ClientConn, s settings.ISettings) *Users {
	return &Users{client:usersRpc.NewUsersRpcClient(conn), settings:s}
}

func (u *Users) SignIn(username string, password string) (*models.User, error) {
	 if res, err := u.client.SignIn(u.GetContext(), &usersRpc.SignInRequest{Username:username, Password:password}); err != nil {
		 return nil, err
	 } else {
		 model := &models.User{}
		 mapUser(res, model)
		 return model, nil
	 }
}

func (u *Users) SignUp(email string, username string, password string) (*models.User, error) {
	if res, err := u.client.SignUp(u.GetContext(), &usersRpc.SignUpRequest{Username:username, Email:email, Password:password}); err != nil {
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

func (u *Users) Signout() {
	//a.Settings.Delete(userDbKey)
	//fmt.Println("Logged out.")
}

func mapUser(rpc *usersRpc.UserResponse, model *models.User){
	model.Id = rpc.Id
	model.Username = rpc.Username
	model.Email = rpc.Email
	model.Token = rpc.Token
}