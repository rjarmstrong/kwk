package rpc

import (
	"crypto/tls"
	"google.golang.org/grpc/credentials"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"golang.org/x/net/context"
)

//serverAddr  := flag.String("server_addr", "127.0.0.1:7777", "The server address in the format of host:port")

func Conn(serverAddress string) *grpc.ClientConn {
	// test
	config := &tls.Config{}
	config.InsecureSkipVerify = true
	cert := credentials.NewTLS(config)

	// production
	//cert := credentials.NewClientTLSFromCert(nil, "")

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(cert))
	if err != nil {
		panic(fmt.Sprintf("Failure: %v", err))
	}
	return conn
}

type Headers struct {

}

func (i *Headers) GetContext() context.Context {
	//if a.Settings.Get(userDbKey, u); u == nil {
	//	fmt.Println("You are not logged in.")
	//	return
	//}

	ctx := metadata.NewContext(
		context.Background(),
		metadata.Pairs("token", "ZingZongZang"),
	)
	return ctx
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
