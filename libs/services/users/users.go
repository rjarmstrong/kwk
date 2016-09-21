package users

import (
	"bitbucket.com/sharingmachine/kwkweb/rpc/usersRpc"
	"google.golang.org/grpc"
	"github.com/kwk-links/kwk-cli/libs/rpc"
	"github.com/kwk-links/kwk-cli/libs/models"
)

const (
	userDbKey = "user"
)

type Users struct {
	client usersRpc.UsersRpcClient
	rpc.Headers
}

func New(conn *grpc.ClientConn) *Users {
	return &Users{client:usersRpc.NewUsersRpcClient(conn)}
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
	//err := a.Settings.Get(userDbKey, u)
	return nil, nil
	// map
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