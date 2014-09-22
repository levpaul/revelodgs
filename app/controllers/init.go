package controllers

import "github.com/revel/revel"

func init() {
	revel.OnAppStart(InitDB)

	// Database Interceptors
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)

	// User Interceptors
	revel.InterceptMethod((*Servers).checkUser, revel.BEFORE)

	// Admin Interceptors
	revel.InterceptMethod((*Admin).checkUser, revel.BEFORE)

	// General Interceptors
	revel.InterceptMethod((*Application).AddUser, revel.BEFORE)
}
