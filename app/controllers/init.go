package controllers

import "github.com/revel/revel"

func init() {
	revel.OnAppStart(InitDB)

	// Database Interceptors
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)

	// User Interceptors
	revel.InterceptMethod(Games.checkUser, revel.BEFORE)
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)

	// Admin Interceptors
}
