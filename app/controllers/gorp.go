package controllers

import (
	"database/sql"
	"rps/app/models"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	r "github.com/revel/revel"
	"github.com/revel/revel/modules/db/app"
)

var (
	Dbm *gorp.DbMap
)

// Db aspect of controllers, other ctrls should inherit
type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

	// User Table
	t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
	t.ColMap("Password").Transient = true
	t.ColMap("Username").SetUnique(true).SetNotNull(true).SetMaxSize(20)
	t.ColMap("Name").SetNotNull(true).SetMaxSize(40)
	t.ColMap("Email").SetNotNull(true).SetMaxSize(150)
	t.ColMap("AccountType").SetNotNull(true).SetMaxSize(20)
	t.ColMap("HashedPassword").SetNotNull(true)

	// Game Table
	t = Dbm.AddTable(models.Game{}).SetKeys(true, "GameId")
	t.ColMap("Name").SetMaxSize(50).SetNotNull(true)
	t.ColMap("AmiId").SetMaxSize(15).SetNotNull(true)
	t.ColMap("Description").SetMaxSize(500)
	t.ColMap("Type").SetMaxSize(20).SetNotNull(true)

	// Servers Table
	t = Dbm.AddTable(models.Server{})
	t.ColMap("UserId").SetNotNull(true)
	t.ColMap("GameId").SetNotNull(true)
	t.ColMap("InstanceId").SetMaxSize(12).SetNotNull(true)
	t.ColMap("LaunchTime").SetNotNull(true)
	t.ColMap("AmiId").SetMaxSize(12).SetNotNull(true)
	t.ColMap("Options").SetMaxSize(300)

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.CreateTables()
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
