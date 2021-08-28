package model

import (
	"asscRegsitration/config"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DB interface {
	Stop() (err error)

	Migrate(version uint) error

	TxFromContext(ctx context.Context) (*sqlx.Tx, bool)
	TxToContext(ctx context.Context, tx *sqlx.Tx) context.Context
	BeginTransaction() (*sqlx.Tx, error)
	CommitTransaction(tx *sqlx.Tx) error
	RollbackTransaction(tx *sqlx.Tx)

	TeamsDB
	PlayersDB
}

type TxKey string

var (
	TxContextKey     = TxKey("tx")
	ErrNoTxInContext = errors.New("no transaction in context")
)

type db struct {
	cfg    *config.DBConfig
	conn   *sqlx.DB
	logger *zap.Logger
}

func NewDB(logger *zap.Logger, cfg *config.DBConfig) (res DB, err error) {
	var (
		dsl = cfg.String()
		sdb *sqlx.DB
	)
	l := logger.With(zap.String("component", "db"))

	sdb, err = sqlx.Open(cfg.Dialect, dsl)
	if err != nil {
		l.Fatal("failed to create sqlx db", zap.Error(err))
		return
	}

	sdb.SetMaxIdleConns(cfg.ConnectionPool)
	sdb.SetMaxOpenConns(cfg.ConnectionMax)
	sdb.SetConnMaxLifetime(60 * time.Minute)

	res = &db{
		logger: l,
		cfg:    cfg,
		conn:   sdb,
	}

	err = sdb.Ping()
	if err != nil {
		logger.With(zap.Error(err)).Fatal("failed ping db")
		return
	}

	// Migrate DB
	err = res.Migrate(cfg.SchemaVersion)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("failed to migrate")
	}

	l.Info("DB initialized")
	return
}

func (d *db) Migrate(version uint) error {
	driver, err := postgres.WithInstance(d.conn.DB, &postgres.Config{})
	if err != nil {
		d.logger.With(zap.Error(err)).Fatal("failed init db driver for migration")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", d.cfg.MigrationPath),
		"postgres", driver)
	if err != nil {
		d.logger.With(zap.Error(err)).Fatal("failed init migration")
		return err
	}
	err = m.Migrate(version)
	if err == migrate.ErrNoChange {
		err = nil
	}
	return err
}

func (d *db) Stop() (err error) {
	return d.conn.Close()
}

func (d *db) closeNamedRows(rows *sqlx.Rows) {
	err := rows.Close()
	if err != nil {
		d.logger.With(zap.Error(err)).Error("failed close named rows")
	}
}
func (d *db) closeStatement(stmt *sqlx.NamedStmt) {
	err := stmt.Close()
	if err != nil {
		d.logger.With(zap.Error(err)).Error("failed close statement")
	}
}

func (d *db) TxFromContext(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(TxContextKey).(*sqlx.Tx)
	return tx, ok
}

func (d *db) TxToContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, TxContextKey, tx)
}

func (d *db) BeginTransaction() (*sqlx.Tx, error) {
	return d.conn.Beginx()
}

func (d *db) CommitTransaction(tx *sqlx.Tx) error {
	return tx.Commit()
}

func (d *db) RollbackTransaction(tx *sqlx.Tx) {
	err := tx.Rollback()
	if err != nil {
		d.logger.With(zap.Error(err)).Error("failed rollback transaction")
	}
}
