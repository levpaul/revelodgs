package controllers

import (
	"errors"
	"rps/app/models"
	"rps/app/routes"

	"github.com/revel/revel"
)

// Responsible for managing interactions with servers in a REST like way
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

func (c *Servers) List() revel.Result {
	results, err := c.Txn.Select(models.ServerResult{},
		"SELECT s.ServerId, g.ShortDesc, s.LaunchTime, s.ExpiryTime, s.State, s.Options FROM Game as 'g' JOIN Server as 's' WHERE s.UserId = ? AND g.GameId = s.GameId", c.connected().UserId)
	if err != nil {
		panic(err)
	}

	var servers []*models.ServerResult
	for _, r := range results {
		s := r.(*models.ServerResult)
		servers = append(servers, s)
	}

	return c.RenderJson(servers)
}

func (c *Servers) ListGames() revel.Result {
	results, err := c.Txn.Select(models.Game{}, `select * from Game`)
	if err != nil {
		panic(err)
	}

	var Games []*models.Game
	for _, r := range results {
		g := r.(*models.Game)
		Games = append(Games, g)
	}

	return c.RenderJson(Games)
}

func (c *Servers) Show(id int) revel.Result {
	var server models.Server
	err := c.Txn.SelectOne(&server,
		`select * from Server where serverid = ? and userid = ?`, id, c.connected().UserId)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// TODO: Tidy this up a little bit. Probably want to just bounce back to userIndex with a flash error "Unauthorised"
			return c.RenderError(errors.New("Unauthorised"))
		}
		panic(err)
	}

	return c.RenderJson(server)
}

func (c *Servers) Delete(id int) revel.Result {
	return c.Render()
}

func (c *Servers) New() revel.Result {
	return c.Render()
}

func (c *Servers) Create() revel.Result {
	return c.Render()
}
