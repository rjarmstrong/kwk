package integration

import (
	"github.com/smartystreets/goconvey/web/server/system"
	"bitbucket.com/sharingmachine/kwkcli/src/snippets"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/src/app"
	"bitbucket.com/sharingmachine/kwkcli/src/rpc"
	"bitbucket.com/sharingmachine/kwkcli/src/setup"
	"google.golang.org/grpc"
	"strings"
	"os/exec"
	"bufio"
	"bytes"
	"fmt"
)

func createApp(conn *grpc.ClientConn, writer *bytes.Buffer, r *bufio.Reader) *app.KwkApp {
	s := system.New()
	t := config.NewJsonSettings(s, "settings")
	h := rpc.NewHeaders(t)
	u := user.NewStdManager(conn, t, h)
	a := snippets.New(conn, t, h, su)
	w := tmpl.NewWriter(writer)
	d := dlg.New(w, r)
	o := openers.New(s, a, w)
	return app.NewApp(a, s, t, o, u, d, w, ch)
}

const (
	sqlContainer = "cass"
	testHost     = "localhost:8000"
	email        = "test@kwk.co"
	username     = "testuser"
	password     = "TestPassword1"
	notLoggedIn  = "You are not logged in please log in: kwk login <username> <password>\n"
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

func getApp(reader *bytes.Buffer, writer *bytes.Buffer) *app.KwkApp {
	conn := rpc.GetConn(testHost)
	r := bufio.NewReader(reader)
	return createApp(conn, writer, r)
}

func cleanup() {
	cmd := exec.Command("/bin/sh", "-c", "docker exec -i "+sqlContainer+" cqlsh -e 'use kwk; TRUNCATE snips; TRUNCATE users_by_email; TRUNCATE users;'")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Println("clean up")
}

func lastLine(input string) string {
	lines := strings.Split(input, "\n")
	l := len(lines)
	if lines[l-1] == "" {
		return lines[l-2]
	}
	return lines[l-1]
}

func line(input string, index int) string {
	lines := strings.Split(input, "\n")
	return lines[index]
}
