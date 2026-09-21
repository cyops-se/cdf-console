package db

import (
	"log"
	"path"
	"server/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(ctx types.Context) {
	filename := path.Join(ctx.Wdir, "cdf-console.db")
	database, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		// Note: Cannot use logger here due to import cycle (logger imports db)
		log.Printf("Failed to connect to database: %s", err.Error())
		return
	}

	// Note: Cannot use logger here due to import cycle (logger imports db)
	log.Printf("Database connected")

	// The User model is special due to the 'password' field, and has
	// user specific routes
	database.AutoMigrate(&types.User{})

	// Generic CRUD data types
	ConfigureTypes(database, types.Log{}, types.KeyValuePair{})
	ConfigureTypes(database, types.User{}, types.Settings{})
	ConfigureTypes(database, types.Device{}, types.InterfaceStats{})
	ConfigureTypes(database, types.Network{}, types.System{}, types.Address{})
	ConfigureTypes(database, types.DataFlow{}, types.DataFlowList{})
	ConfigureTypes(database, types.NetworkTemplate{}, types.HubTemplate{}, types.EndpointTemplate{})

	DB = database
}

func ConfigureTypes(database *gorm.DB, datatypes ...interface{}) {
	for _, datatype := range datatypes {
		stmt := &gorm.Statement{DB: database}
		stmt.Parse(datatype)
		name := stmt.Schema.Table
		types.RegisterType(name, datatype)
		database.AutoMigrate(datatype)
	}
}
