package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/KeKsBoTer/socialloot/models"
	"github.com/astaxie/beego"
)

type AuthController struct {
	beego.Controller

	User *models.User
}

// UserInfoKey is the session key for the user id
const UserInfoKey = "userinfo"

func (c *AuthController) Prepare() {
	isLogin := c.GetSession(UserInfoKey) != nil
	if isLogin {
		c.User = c.GetLogin()
	}

	// set data for html rendering
	if c.Ctx.Input.IsGet() {
		c.Data["IsLogin"] = isLogin
		c.Data["User"] = c.User
		if dst := c.GetString("dest"); len(dst) > 1 {
			c.Data["Dest"] = dst
		}

		c.Data["HeadStyles"] = []string{}
		c.Data["HeadScripts"] = []string{}
		c.Data["URL"] = c.Ctx.Input.URI()

		c.Layout = "base.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["BaseHeader"] = "components/header.tpl"
	}
}

func (c *AuthController) RedirectForm() {
	if dst := c.GetString("dest"); len(dst) > 0 {
		c.Redirect(dst, http.StatusSeeOther)
	}
}

func (c *AuthController) GetLogin() *models.User {
	if i, ok := c.GetSession(UserInfoKey).(int); ok {
		u := &models.User{
			Id: i,
		}
		u.Read()
		return u
	}
	return nil
}

func (c *AuthController) IsLogin() bool {
	return c.User != nil
}

func (c *AuthController) DelLogin() {
	c.DelSession(UserInfoKey)
}

func (c *AuthController) SetLogin(user *models.User) {
	c.SetSession(UserInfoKey, user.Id)
}

func (c *AuthController) LoginPath() string {
	return c.URLFor("LoginController.Login")
}

// NeedsAuthController redirects to Login page if user is not authenticated
type NeedsAuthController struct {
	AuthController
}

func (c *NeedsAuthController) Prepare() {
	c.AuthController.Prepare()
	if !c.IsLogin() {
		if c.Ctx.Input.IsGet() {
			c.RedirectForm()
			if !c.Ctx.Output.IsRedirect() {
				c.Ctx.Redirect(http.StatusUnauthorized, c.URLFor("IndexController.Index"))
			}
		} else if c.Ctx.Input.IsPost() {
			r := ApiResponse{
				Success: false,
				Message: "unauthorized",
				Dest:    c.URLFor("LoginController.LoginPage"),
			}
			j, _ := json.Marshal(r)
			c.CustomAbort(http.StatusUnauthorized, string(j))
		} else {
			c.Abort("401")
		}
	}
}