package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type MySQL interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Begin() (*sql.Tx, error)
	Ping() error
	Close() error
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	Stats() sql.DBStats
}

type Option struct {
	Username string
	Password string
	HostPort string
	DBName   string

	// Default is unlimited
	// set to lest than equal to 0 means make it unlimited or default setting,
	// more open connection means less time taken to perform query
	MaxOpenConn int

	// MaxIdleConn default is 2
	// set to lest than equal to 0 means not allow any idle connection,
	// more idle connection in the pool will improve performance,
	// since no need to establish connection from scratch)
	// by set idle connection to 0, a new connection has to be created from scratch for each operation
	// ! should be <= MaxOpenConn
	MaxIdleConn int

	// MaxLifetime set max length of time that a connection can be reused for.
	// Setting to 0 means that there is no maximum lifetime and
	// the connection is reused forever (which is the default behavior)
	// the shorter lifetime result in more memory useage
	// since it will kill the connection and recreate it
	MaxLifetime time.Duration
}

// New create new MySQL instance
func New(opt Option) MySQL {
	connectionstr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", opt.Username, opt.Password, opt.HostPort, opt.DBName)
	logrus.Infof("connectionstr: %s", connectionstr)

	db, err := sql.Open("mysql", connectionstr)
	if err != nil {
		panic(err.Error())
	}

	if opt.MaxOpenConn != 0 {
		db.SetMaxOpenConns(opt.MaxOpenConn)
	}

	if opt.MaxIdleConn != 0 {
		db.SetMaxIdleConns(opt.MaxIdleConn)
	}

	if err = db.Ping(); err != nil {
		logrus.Fatal(errors.Wrap(err, "ping mysql"))
	}

	logrus.Infof("%-7s %s", "MySQL", "✅")

	return db
}

func DefaultOption() Option {
	return Option{
		Username:    "root",
		Password:    "password",
		HostPort:    "localhost:3306",
		DBName:      "",
		MaxOpenConn: 0,
		MaxIdleConn: 2,
		MaxLifetime: 0,
	}
}
