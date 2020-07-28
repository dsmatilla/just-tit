package models

import (
    "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"os"
)

const (
	dbAlias       = "default"
	dbUse		  = "true"
)

func init() {
	if dbUse == "true" {
		initializeOrm()
	}
}

func initializeOrm() {
	dbHost := os.Getenv("MYSQLHOST")
	dbName := os.Getenv("MYSQLNAME")
	dbUser := os.Getenv("MYSQLUSER")
	dbPass := os.Getenv("MYSQLPASS")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(Score))
	
	if err := orm.RegisterDataBase(dbAlias, "mysql", dbUser + ":" + dbPass + "@tcp(" + dbHost +")/" + dbName + "?charset=utf8"); err != nil {
		log.Println(err)
	}

	orm.DefaultTimeLoc = time.UTC

	orm.Debug = true
	force := false   // Drop table and re-create.
	verbose := true // Print log
	// generate Tables
	if err := orm.RunSyncdb(dbAlias, force, verbose); err != nil {
		log.Println(err)
	}
}