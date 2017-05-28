package app

import (
	"github.com/kwk-super-snippets/types"
	"os"
	"testing"
	"bytes"
	"github.com/stretchr/testify/assert"
	"os/exec"
)

var app *KwkCLI

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func Test_Grpc(t *testing.T) {
	defer cleanMysql(t)
	defer cleanRoach(t)
	cc, err := GetConn("localhost:8000", true)
	assert.Equal(t, nil, err)
	uc := types.NewUsersClient(cc)

	t.Log("SIGN-UP")
	sur := &types.SignUpRequest{Username: "test1", Email: "test1@kwk.co", Password: "Password1"}
	res, err := uc.SignUp(Ctx(), sur)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.AccessToken)
	assert.Equal(t, sur.Username, res.User.Username)

	t.Log("SIGN-IN")
	sires, err := uc.SignIn(Ctx(), &types.SignInRequest{Username: sur.Username, Password: sur.Password})
	assert.Nil(t, err)
	assert.Equal(t, sur.Email, sires.User.Email)
	assert.Equal(t, sur.Username, sires.User.Username)
	assert.NotEmpty(t, sires.AccessToken)

	principal.AccessToken = sires.AccessToken

	t.Log("CREATE SNIPPET")
	sc := types.NewSnippetsClient(cc)
	al := types.NewAlias(sur.Username, "testpouch", "hello", "js")
	cr := &types.CreateRequest{Alias: al, Content: "console.log('hello')"}
	snres, err := sc.Create(Ctx(), cr)
	assert.Nil(t, err)
	assert.Equal(t, al.URI(), snres.Snippet.Alias.URI())
	assert.Equal(t, cr.Content, snres.Snippet.Content)

	t.Log("GET SNIPPET")
	gres, err := sc.Get(Ctx(), &types.GetRequest{Alias: al})
	assert.Equal(t, nil, err)
	assert.Equal(t, al.Name, gres.Items[0].Name())

}

func cleanRoach(t *testing.T) {
	cmd := exec.Command("docker", "exec", "roach", "./cockroach.sh", "sql", "--database", "kwk_test", "--insecure", "-e", "truncate snippets; truncate pouches;")
	er := &bytes.Buffer{}
	cmd.Stderr = er
	_, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
}

func cleanMysql(t *testing.T) {
	cmd := exec.Command("docker", "exec", "mysql", "mysql", "-uroot", "-prootPassword", "-D", "kwk_users_test", "-e", "delete from users; delete from password_attempts;")
	er := &bytes.Buffer{}
	cmd.Stderr = er
	_, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(string(o))
	//t.Log(er.String())
}
