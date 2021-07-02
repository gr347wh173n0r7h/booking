package database

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/booking/config"
)

var (
	// ErrorDNE represents a Does Not Exit error
	ErrorDNE = pg.ErrNoRows
)

// Database defines interface for database interaction
type Database interface {
	Conn() *pg.DB
	Ping(ctx context.Context) error
	CreateSchema(models []interface{}) error
}

type postgres struct {
	conn *pg.DB
}

// NewPGSQLClient returns a postgres implementation of Database
func NewPGSQLClient(ctx context.Context, c *config.Config) (Database, error) {
	opt, err := pg.ParseURL(c.DBURL)
	if err != nil {
		return nil, err
	}

	conn := pg.Connect(opt)

	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &postgres{
		conn: pg.Connect(opt),
	}, nil
}

func (db *postgres) Conn() *pg.DB {
	return db.conn
}

func (db *postgres) Ping(ctx context.Context) error {
	return db.conn.Ping(ctx)
}

func (db *postgres) CreateSchema(models []interface{}) error {
	for _, model := range models {
		err := db.conn.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
