package models

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	DbManager struct {
		seeder *Seeder
		Db     *gorm.DB

		dbHost    string
		dbPort    string
		dbUser    string
		dbPass    string
		dbName    string
		dbSSLMode string
	}
)

func NewDbManager() (dbMgr *DbManager, err error) {
	dbMgr = &DbManager{}
	if err = dbMgr.Validate(); err != nil {
		return
	}
	if err = dbMgr.Setup(); err != nil {
		return
	}
	dbMgr.seeder = NewSeeder(dbMgr.Db)
	if err = dbMgr.PostMigrate(); err != nil {
		return
	}
	return
}

func (dbMgr *DbManager) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbMgr.dbHost, dbMgr.dbPort, dbMgr.dbUser, dbMgr.dbPass, dbMgr.dbName, dbMgr.dbSSLMode)
}

func (dbMgr *DbManager) Validate() (err error) {
	// Get database connection details from environment variables
	dbMgr.dbHost = os.Getenv("DB_HOST")
	dbMgr.dbPort = os.Getenv("DB_PORT")
	dbMgr.dbUser = os.Getenv("DB_USER")
	dbMgr.dbPass = os.Getenv("DB_PASSWORD")
	dbMgr.dbName = os.Getenv("DB_NAME")
	dbMgr.dbSSLMode = os.Getenv("DB_SSLMODE")
	if dbMgr.dbSSLMode == "" {
		dbMgr.dbSSLMode = "disable" // default for local development
	}

	// Validate required environment variables
	if dbMgr.dbHost == "" || dbMgr.dbPort == "" || dbMgr.dbUser == "" || dbMgr.dbPass == "" || dbMgr.dbName == "" {
		return fmt.Errorf("missing required database environment variables")
	}
	return
}

func (dbMgr *DbManager) AutoMigrate() (err error) {
	if err = dbMgr.Db.AutoMigrate(
		&User{},
		&Profile{},
		&Account{},
		&Project{},
		&ProjectPermission{},
		&RoleAccess{},
		&Plan{},
	); err != nil {
		return
	}
	return
}

func (dbMgr *DbManager) Setup() (err error) {
	if dbMgr.Db, err = gorm.Open(postgres.Open(dbMgr.GetDSN()), &gorm.Config{}); err != nil {
		return
	}
	if err = dbMgr.AutoMigrate(); err != nil {
		return
	}
	return
}

func (dbMgr *DbManager) PostMigrate() (err error) {
	if err = dbMgr.seeder.Seed(); err != nil {
		return
	}
	return
}
