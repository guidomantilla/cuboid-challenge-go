package db

import (
	"cuboid-challenge/app/config"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CONN global gorm db connection.
var CONN *gorm.DB

// Connect set the gorm db connection.
func Connect() *gorm.DB {
	if CONN != nil {
		return CONN
	}

	var err error
	if CONN, err = gorm.Open(driver(), &gorm.Config{}); err != nil {
		panic(fmt.Errorf("failed to open the database connection. %w", err))
	}

	registerDBValidationsHooks(CONN)

	return CONN
}

func driver() gorm.Dialector {
	d := config.ENV.DBDriver
	if d == "sqlite" {
		return sqlite.Open(config.ENV.DBName)
	}

	panic(driverError{d}.Error())
}

type driverError struct{ name string }

func (d driverError) Error() string {
	return fmt.Sprintf("DB driver '%s' not supported", d.name)
}
