package rpc

import (
	"crypto/tls"
	//"google.golang.org/grpc/credentials"
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//serverAddr  := flag.String("server_addr", "127.0.0.1:7777", "The server address in the format of host:port")

func Conn(serverAddress string) *grpc.ClientConn {
	// test
	config := &tls.Config{}
	config.InsecureSkipVerify = true
	//cert := credentials.NewTLS(config)

	// production
	//cert := credentials.NewClientTLSFromCert(nil, "")
	//grpc.WithTransportCredentials(cert)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Failure: %v", err))
	}
	return conn
}

func NewHeaders(t settings.ISettings) *Headers {
	return &Headers{settings: t}
}

type Headers struct {
	settings settings.ISettings
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
