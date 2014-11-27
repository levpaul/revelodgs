package controllers

import (
	"rps/app/models"
	"rps/app/routes"

	"github.com/revel/revel"
)

// Responsible for managing admin stuff
type Admin struct {
	Application
}

func (c *Admin) isUserAdmin(user *models.User) bool {
	var result models.User

	err := c.Txn.SelectOne(&result, `select u.* from user as u join accounttype as a where u.accounttype == a.accounttypeid and u.userid == ? and a.name == ?`, user.UserId, "ADMIN") //models.AccountTypeAdmin)
	c.RenderArgs["User"] = result
	if err != nil {
		return false
	} else {
		return true
	}

}

/* Checks if a user is connected, then makes sure that the user
 *  is an admin user. */
func (c *Admin) checkUser() revel.Result {
	user := c.connected()
	if user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	if c.isUserAdmin(user) == false {
		c.Flash.Error("Unauthorised")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c *Admin) Index() revel.Result {
	return c.Render()
}
