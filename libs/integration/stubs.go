package integration

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bytes"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/app"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"google.golang.org/grpc"
	"bufio"
	"os/exec"
	"strings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/search"
	"fmt"
)

func createApp(conn *grpc.ClientConn, writer *bytes.Buffer, r *bufio.Reader) *app.KwkApp{
	s := system.New()
	t := settings.New(s, "settings")
	h := rpc.NewHeaders(t)
	u := users.New(conn, t, h)
	a := aliases.New(conn, t, h)
	w := gui.NewTemplateWriter(writer)
	d := gui.NewDialogues(w, r)
	o := openers.New(s, a, w)
	ch := search.New(conn, t, h)
	return app.NewKwkApp(a, s, t, o, u, d, w, ch)
}

const (
	sqlContainer="cass"
	testHost="localhost:8000"
	email="test@kwk.co"
	username="testuser"
	password="TestPassword1"
	notLoggedIn="You are not logged in please log in: kwk login <username> <password>\n"
)

func signin(reader *bytes.Buffer, kwk *app.KwkApp) {
	reader.WriteString(username + "\n")
	reader.WriteString(password + "\n")
	kwk.Run("signin")
}

func signup(reader *bytes.Buffer, kwk *app.KwkApp) {
	reader.WriteString(email + "\n")
	reader.WriteString(username + "\n")
	reader.WriteString(password + "\n")
	kwk.Run("signup")
}

func getApp(reader *bytes.Buffer, writer *bytes.Buffer) *app.KwkApp{
	conn := rpc.Conn(testHost);
	r := bufio.NewReader(reader)
	return createApp(conn, writer, r)
}

func cleanup(){
		cmd := exec.Command("/bin/sh", "-c", "docker exec -i " + sqlContainer + " cqlsh -e 'use kwk; TRUNCATE snips; TRUNCATE users_by_email; TRUNCATE users;'")
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	fmt.Println("clean up")
}

func lastLine(input string) string {
	lines := strings.Split(input, "\n")
	l := len(lines)
	if lines[l-1] == ""{
		return lines[l-2]
	}
	return lines[l-1]
}

func line(input string, index int) string {
	lines := strings.Split(input, "\n")
	return lines[index]
}