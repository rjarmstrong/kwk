package controllers

import (
	"github.com/kwk-links/kwk-cli/libs/services/users"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"bufio"
	"os"
)

type AccountController struct {
	service users.IUsers
	settings settings.ISettings
}

func NewAccountController(u users.IUsers, s settings.ISettings) *AccountController {
	return &AccountController{service:u, settings:s}
}

func (c *AccountController) Get(){
	if u, err := c.service.Get(); err != nil {
		fmt.Println(err)
		fmt.Println("You are not logged in please log in: kwk login <username> <password>")
	} else {
		fmt.Println("~~~~~~ Your Profile ~~~~~~~~~")
		fmt.Println(gui.Build(2, gui.Space) + gui.Build(11, "~") + gui.Build(2, gui.Space))
		fmt.Println(gui.Build(6, gui.Space) + gui.Build(3, gui.UBlock) + gui.Build(6, gui.Space))
		fmt.Println(gui.Build(5, gui.Space) + gui.Build(5, gui.UBlock) + gui.Build(5, gui.Space))
		fmt.Println(gui.Build(6, gui.Space) + gui.Build(3, gui.UBlock) + gui.Build(6, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(2, gui.Space) + gui.Build(11, "~") + gui.Build(2, gui.Space))

		fmt.Printf("Email:      %v\n", u.Email)
		fmt.Printf("Username:   %v\n", u.Username)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	}
}

func (c *AccountController) SignUp(email string, username string, password string){
	if u, err := c.service.SignUp(email, username, password); err != nil {
		fmt.Println(err)
	} else {
		if len(u.Token) > 50 {
			c.settings.Upsert("me", u)
			fmt.Printf("Welcome to kwk %s! You're signed in already.", u.Username)
		}
	}
}


func (c *AccountController) SignIn(username string, password string){
	if username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(gui.Colour(gui.LightBlue, "Your Kwk Username: "))
		usernameBytes, _, _ := reader.ReadLine()
		username = string(usernameBytes)
	}
	if password == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(gui.Colour(gui.LightBlue, "Your Password: "))
		passwordBytes, _, _ := reader.ReadLine()
		password = string(passwordBytes)
	}
	if u, err := c.service.SignIn(username, password); err != nil {
		fmt.Println(err)
	} else {
		if len(u.Token) > 50 {
			c.settings.Upsert("me", u)
			fmt.Printf("Welcome back %s!", u.Username)
		}
	}
}

func (c *AccountController) SignOut(){
	c.settings.Delete("me")
}