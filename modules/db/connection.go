package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

const (
	DefaultDbKey = "default"
)

type DbService interface {
	AddConfig(profName string, config *DbProfile)
	RemoveConfig(profName string)
	Default() *gorm.DB
	Get(profName string) *gorm.DB
	AutoMigrate() bool
}

// DbProfile = db configuration
type DbProfile struct {
	Connection string
	Host       string
	Port       string
	Database   string
	Username   string
	Password   string
	Locale     string
	DbUrl      string
	Logging    bool
}

func (module *DbModule) addDefaultConfig() {
	module.autoMigrate = module.config.Getenv("DB_AUTOMIGRATION", "") == "true"
	module.AddConfig(DefaultDbKey, &DbProfile{
		Connection: module.config.Getenv("DB_CONNECTION", ""),
		Host:       module.config.Getenv("DB_HOST", ""),
		Port:       module.config.Getenv("DB_PORT", ""),
		Database:   module.config.Getenv("DB_DATABASE", ""),
		Username:   module.config.Getenv("DB_USERNAME", ""),
		Password:   module.config.Getenv("DB_PASSWORD", ""),
		Locale:     module.config.Getenv("DB_LOCALE", ""),
		DbUrl:      module.config.Getenv("DATABASE_URL", ""),
		Logging:    module.config.Getenv("DB_LOGGING", "") == "true",
	})
}

// impl `DbService` start

// AddConfig = add configuration
func (module *DbModule) AddConfig(profName string, config *DbProfile) {
	var err error
	gormConfig := gorm.Config{}

	if !config.Logging {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	if config.Connection == "mysql" {
		suffix := "?loc="
		if config.Locale == "" {
			suffix += "UTC"
		} else {
			suffix += url.QueryEscape(config.Locale)
		}
		suffix += "&parseTime=true&multiStatements=true"
		module.db[profName], err = gorm.Open(
			mysql.Open(config.Username+":"+config.Password+"@tcp("+config.Host+":"+config.Port+")/"+config.Database+suffix),
			&gormConfig,
		)
	} else if config.Connection == "postgres" {
		var connInfo string
		if config.DbUrl != "" {
			connInfo = config.DbUrl
		} else {
			connInfo = fmt.Sprintf("host='%s' port=%s user='%s' "+
				"password='%s' dbname='%s' sslmode=disable",
				config.Host, config.Port, config.Username, config.Password, config.Database)
		}
		module.db[profName], err = gorm.Open(
			postgres.Open(connInfo),
			&gormConfig,
		)
	}

	var sqlDB *sql.DB
	if err == nil {
		sqlDB, err = module.db[profName].DB()
	}

	if err != nil {
		log.Fatalf("DB profile `%s` connect error: %v", profName, err)
	}

	err = sqlDB.Ping()

	if err != nil {
		log.Fatalf("DB profile `%s` ping error: %v", profName, err)
	}

	if module.autoMigrate && config.Connection == "postgres" {
		module.db[profName].Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	}

	if err := module.db[profName].Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}

// RemoveConfig = remove configuration
func (module *DbModule) RemoveConfig(profName string) {
	delete(module.db, profName)
}

// Default : get default DB profile
func (module *DbModule) Default() *gorm.DB {
	return module.db[DefaultDbKey]
}

// Get : get DB profile
func (module *DbModule) Get(profName string) *gorm.DB {
	return module.db[profName]
}

func (module *DbModule) AutoMigrate() bool {
	return module.autoMigrate
}

// impl `DbService` end
