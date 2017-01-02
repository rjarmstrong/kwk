package rpc

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
	"os"
)

func GetConn(serverAddress string) *grpc.ClientConn {
	var opts []grpc.DialOption

	var trustCerts = false
	if ok := os.Getenv("TRUST_ALL_CERTS"); ok != "" && serverAddress == "localhost:8000" {
		trustCerts = true
	}
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify:trustCerts,
	})
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		panic(fmt.Sprintf("Failure: %v", err))
	}
	return conn
}

func NewHeaders(t config.Settings) *Headers {
	return &Headers{settings: t}
}

type Headers struct {
	settings config.Settings
}

func (i *Headers) GetContext() context.Context {
	u := &models.User{}
	if err := i.settings.Get(models.ProfileFullKey, u); err != nil {
		return context.Background()
	} else {
		ctx := metadata.NewContext(
			context.Background(),
			metadata.Pairs(models.TokenHeaderName, u.Token),
		)
		return ctx
	}
}

/*
func handleResponse(path string, i interface{}, r *http.Response) {
	switch {
	case r.StatusCode == http.StatusBadRequest :
		if i != nil {
			system.PrettyPrint(i)
		} else {
			fmt.Println(r.StatusCode)
		}
	case r.StatusCode == http.StatusForbidden :
		fmt.Println("Sign in please: 'kwk signin <username> <password>'")
	case r.StatusCode == http.StatusNotFound :
		fmt.Println(path + " not found.")
	case r.StatusCode != http.StatusOK :
		fmt.Println(r)
		fmt.Println(r.StatusCode)
	}
}

if e != nil {
		fmt.Print("kwk server is unavailable, please try again later or tweet us @kwklinks.")
		return
	}
*/
