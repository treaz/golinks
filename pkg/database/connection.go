package database

import (
	"crhuber/golinks/pkg/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connection represents a database connection
type DbConnection struct {
	Db *gorm.DB
}

// NewConnection create a new connection
func NewConnection(dbType, dbDSN string) (*DbConnection, error) {
	var db *gorm.DB
	var err error
	switch dbType {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(dbDSN), &gorm.Config{})
	default:
		db, err = gorm.Open(sqlite.Open(dbDSN), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}
	db.Logger.LogMode(logger.Info)
	return &DbConnection{Db: db}, nil
}

// Close closes the underlying db connection
func (c *DbConnection) Close() error {
	dbConn, err := c.Db.DB()
	if err != nil {
		return err
	}

	return dbConn.Close()
}

// RunMigration runs db migrations
func (c *DbConnection) RunMigration() error {
	// Drop the is_parameterized column if it exists (migration cleanup)
	if c.Db.Migrator().HasColumn(&models.Link{}, "is_parameterized") {
		c.Db.Migrator().DropColumn(&models.Link{}, "is_parameterized")
	}
	
	// Drop link_tags junction table if it exists (migration cleanup)
	if c.Db.Migrator().HasTable("link_tags") {
		c.Db.Migrator().DropTable("link_tags")
	}
	
	// Drop tags table if it exists (migration cleanup)
	if c.Db.Migrator().HasTable("tags") {
		c.Db.Migrator().DropTable("tags")
	}
	
	return c.Db.AutoMigrate(&models.Link{})
}
