package controllers

import (
	"rps/app/routes"

	"github.com/revel/revel"
)

// Responsible for managing interactions with servers
type Servers struct {
	Application
}

func (c *Servers) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c *Servers) Index() revel.Result {
	// results, err := c.Txn.Select(models.Game{},
	// 	`select * from Game where UserOne = ? or UserTwo = ?`, c.connected().UserId, c.connected().UserId)
	// if err != nil {
	// 	panic(err)
	// }

	// var Servers []*models.Game
	// for _, r := range results {
	// 	g := r.(*models.Game)
	// 	Servers = append(Servers, g)
	// }

	return c.Render()
}
