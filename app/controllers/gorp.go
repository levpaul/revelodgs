package controllers

import (
	"database/sql"
	"regexp"
	"rps/app/models"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
	"github.com/revel/revel/modules/db/app"

	"time"
)

var (
	Dbm *gorp.DbMap
)

// Db aspect of controllers, other ctrls should inherit
type GorpController struct {
	*revel.Controller
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
	t = Dbm.AddTable(models.Game{}).SetKeys(false, "Name")
	t.ColMap("Name").SetMaxSize(10)
	t.ColMap("AmiId").SetMaxSize(15).SetNotNull(true)
	t.ColMap("ShortDesc").SetMaxSize(50).SetNotNull(true)
	t.ColMap("LongDesc").SetMaxSize(500)
	t.ColMap("Type").SetMaxSize(20).SetNotNull(true)

	// Servers Table
	t = Dbm.AddTable(models.Server{}).SetKeys(true, "ServerId")
	t.ColMap("UserId").SetNotNull(true)
	t.ColMap("GameId").SetNotNull(true)
	t.ColMap("InstanceId").SetMaxSize(12).SetNotNull(true)
	t.ColMap("LaunchTime").SetNotNull(true)
	t.ColMap("ExpiryTime").SetNotNull(true)
	t.ColMap("State").SetMaxSize(15).SetNotNull(true)
	t.ColMap("AmiId").SetMaxSize(12).SetNotNull(true)
	t.ColMap("Options").SetMaxSize(300)

	Dbm.TraceOn("[gorp]", revel.INFO)
	Dbm.CreateTablesIfNotExists()

	// ===================================================
	// POPULATE WITH TEST DATA

	// Helper functions
	checkErr := func(e error) {
		uniqueConstraint := regexp.MustCompile("^UNIQUE constraint failed")
		if e != nil {
			if !uniqueConstraint.MatchString(e.Error()) {
				panic(e)
			}
		}
	}

	insertGame := func(g models.Game) { checkErr(Dbm.Insert(&g)) }
	insertUser := func(u models.User) { checkErr(Dbm.Insert(&u)) }
	insertServer := func(s models.Server) {
		// This table has no PK so make sure we're not adding a dupe here (ignore time since that's dynamically generated)
		var serverResult models.Server
		err := Dbm.SelectOne(&serverResult, "select * from server where userid=? and gameid=? and instanceid=? and state=? and amiid=? and options=?", s.UserId, s.GameId, s.InstanceId, s.State, s.AmiId, s.Options)
		if err != nil {
			checkErr(Dbm.Insert(&s))
		}
	}

	// Add user data
	hp, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	user := models.User{Username: "admin", Name: "Admin Dude", Email: "admin@you.com", AccountType: models.UserAccountTypeAdmin, HashedPassword: hp}
	insertUser(user)
	hp, _ = bcrypt.GenerateFromPassword([]byte("tester"), bcrypt.DefaultCost)
	user = models.User{Username: "tester", Name: "Test Man", Email: "test@me.com", AccountType: models.UserAccountTypeUser, HashedPassword: hp}
	insertUser(user)

	// Add game data
	game := models.Game{Name: "csgo", AmiId: "ami-2342341", ShortDesc: "CounterStrike: Global Offensive", LongDesc: "A shooting game where you make money and buy better guns and more grenades. You can also make yourself look pretty via dressups.", Type: models.GameTypeSteam}
	insertGame(game)
	game = models.Game{Name: "css", AmiId: "ami-9asdfa1", ShortDesc: "CounterStrike: Source", LongDesc: "Oldschool upgrade of original CounterStrike. None of that crap with shields though.", Type: models.GameTypeSteam}
	insertGame(game)
	game = models.Game{Name: "tf2", AmiId: "ami-9123aaa", ShortDesc: "Team Fortress 2", LongDesc: "A hat collecting game which is very popular. Potential for pay to win!", Type: models.GameTypeSteam}
	insertGame(game)
	game = models.Game{Name: "mc", AmiId: "ami-minem333", ShortDesc: "MineCraft", LongDesc: "3D blocks powered by DirectX, runs on Java.", Type: models.GameTypeMinecraft}
	insertGame(game)

	// Add servers data
	server := models.Server{UserId: 1, GameId: 2, InstanceId: "i-3423423", LaunchTime: time.Now().Add(time.Hour * -1), ExpiryTime: time.Now().Add(time.Hour * 99), State: models.ServerStateRunning, AmiId: "ami-1234234g", Options: "maxplayers=99"}
	insertServer(server)
	server = models.Server{UserId: 1, GameId: 3, InstanceId: "i-bsadfss", LaunchTime: time.Now().Add(time.Hour * -108), ExpiryTime: time.Now().Add(time.Hour * -48), State: models.ServerStateStopped, AmiId: "ami-1asd34g", Options: "maxplayers=99"}
	insertServer(server)
	server = models.Server{UserId: 1, GameId: 4, InstanceId: "i-asdfaaaa", LaunchTime: time.Now().Add(time.Hour * -2), ExpiryTime: time.Now().Add(time.Hour * 9), State: models.ServerStateRunning, AmiId: "ami-1asdd34g", Options: "serversize=large"}
	insertServer(server)
	server = models.Server{UserId: 2, GameId: 1, InstanceId: "i-3423423", LaunchTime: time.Now().Add(time.Hour * -3), ExpiryTime: time.Now().Add(time.Hour * 912344), State: models.ServerStateRunning, AmiId: "ami-absdfdd", Options: "maxplayers=16"}
	insertServer(server)
}

func (c *GorpController) Begin() revel.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
