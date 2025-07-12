package models

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gsarmaonline/goiter/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	sqliteDbPath = "gorm.db" // Default SQLite database file path
)

var (
	Models = []interface{}{
		&User{},
		&Profile{},
		&Account{},
		&Project{},
		&ProjectPermission{},
		&RoleAccess{},
		&Plan{},
	}
)

type (
	DbManager struct {
		seeder *Seeder
		cfg    *config.Config
		Db     *gorm.DB

		dbHost    string
		dbPort    string
		dbUser    string
		dbPass    string
		dbName    string
		dbSSLMode string

		models map[string]interface{}

		dbType config.DbTypeT
	}
)

func NewDbManager(cfg *config.Config) (dbMgr *DbManager, err error) {
	dbMgr = &DbManager{
		cfg: cfg,
	}
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

func (dbMgr *DbManager) ConnectSqlite() (err error) {
	dbMgr.Db, err = gorm.Open(sqlite.Open(sqliteDbPath), &gorm.Config{})
	return
}

func (dbMgr *DbManager) GetDSN() (dsn string) {
	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbMgr.dbHost, dbMgr.dbPort, dbMgr.dbUser, dbMgr.dbPass, dbMgr.dbName, dbMgr.dbSSLMode)
	log.Println(dsn)
	return
}

func (dbMgr *DbManager) Validate() (err error) {
	// Get database connection details from environment variables
	dbMgr.dbHost = dbMgr.cfg.DBHost
	dbMgr.dbPort = dbMgr.cfg.DBPort
	dbMgr.dbUser = dbMgr.cfg.DBUser
	dbMgr.dbPass = dbMgr.cfg.DBPassword
	dbMgr.dbName = dbMgr.cfg.DBName
	dbMgr.dbType = dbMgr.cfg.DBType
	dbMgr.dbSSLMode = "disable"
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
	for _, model := range Models {
		log.Println(reflect.TypeOf(model), dbMgr.Db, "migrating")
		if err = dbMgr.Db.Debug().AutoMigrate(model); err != nil {
			return
		}
	}
	return
}

func (dbMgr *DbManager) Setup() (err error) {

	if dbMgr.dbType == config.SqliteDbType {
		log.Println("Using SQLite for local development")
		// Default to SQLite for local development
		if err = dbMgr.ConnectSqlite(); err != nil {
			return
		}
	} else {
		log.Println("Using Postgres")
		if dbMgr.Db, err = gorm.Open(postgres.Open(dbMgr.GetDSN()), &gorm.Config{}); err != nil {
			return
		}
	}
	if dbMgr.cfg.Mode == config.ModeDev {
		dbMgr.DropModels()
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

func (db *DbManager) DropModels() (err error) {
	log.Println("Dropping all models")
	for _, model := range Models {
		if err = db.Db.Migrator().DropTable(model); err != nil {
			return fmt.Errorf("failed to drop table for model %T: %w", model, err)
		}
	}
	return nil
}
