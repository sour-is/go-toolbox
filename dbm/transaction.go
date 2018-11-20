package dbm

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/jmoiron/sqlx"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/uuid"

	sq "github.com/Masterminds/squirrel"
	opentracing "github.com/opentracing/opentracing-go"
)

// Tx database transaction
type Tx struct {
	*sql.Tx
	context.Context
	DbType      string
	Placeholder sq.PlaceholderFormat
	Returns     bool
}

// NewTx create new transaction
func (db DB) NewTx(ctx context.Context, readonly bool) (tx *Tx, err error) {
	sp, nctx := opentracing.StartSpanFromContext(ctx, "NewTx")
	defer sp.Finish()

	opts := new(sql.TxOptions)
	if readonly {
		opts.Isolation = sql.LevelReadCommitted
		opts.ReadOnly = true
	}

	tx = new(Tx)
	tx.Context = nctx
	tx.Tx, err = db.Conn.BeginTx(nctx, opts)
	tx.Placeholder = db.Placeholder
	tx.DbType = db.DbType
	tx.Returns = db.Returns

	return
}

// Transaction starts a new database tranaction and executes the supplied func.
func Transaction(txFunc func(*Tx) error) (err error) {
	return stdDB.TransactionContext(context.Background(), txFunc)
}

// TransactionContext starts a new database tranaction and executes the supplied func with context.
func TransactionContext(ctx context.Context, txFunc func(*Tx) error) (err error) {
	return stdDB.TransactionContext(ctx, txFunc)
}

// Transaction starts a new database transction and executes the supplied func.
func (db DB) Transaction(fn func(*Tx) error) error {
	return db.TransactionContext(context.Background(), fn)
}

// TransactionContext starts a new database transction with context and executes the supplied func.
func (db DB) TransactionContext(ctx context.Context, txFunc func(*Tx) error) (err error) {
	tx, err := db.NewTx(ctx, false)

	if err != nil {
		log.Error(err.Error())
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())

			debug.PrintStack()
			return
		}
		err = tx.Commit()
	}()
	err = txFunc(tx)
	return err
}

// QueryContext starts a new database tranaction and executes the supplied func with context.
func QueryContext(ctx context.Context, txFunc func(*Tx) error) (err error) {
	return stdDB.QueryContext(ctx, txFunc)
}

// QueryContext starts a new database transction with context and executes the supplied func.
func (db DB) QueryContext(ctx context.Context, txFunc func(*Tx) error) (err error) {
	tx, err := db.NewTx(ctx, true)

	if err != nil {
		log.Error(err.Error())
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())

			debug.PrintStack()
			return
		}
		err = tx.Commit()
	}()
	err = txFunc(tx)
	return err
}

// Transactionx starts a new database tranaction and executes the supplied func.
func Transactionx(txFunc func(*sqlx.Tx) error) (err error) {
	return stdDB.Transactionx(txFunc)
}

// Transactionx starts a new database tranaction and executes the supplied func.
func (db DB) Transactionx(txFunc func(*sqlx.Tx) error) (err error) {
	dbx := sqlx.NewDb(db.Conn, db.DbType)

	tx, err := dbx.Beginx()
	if err != nil {
		log.Error("dbm.Transaction ", err.Error())
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())

			debug.PrintStack()
			return
		}
		err = tx.Commit()
	}()
	err = txFunc(tx)
	return err
}

// TransactionContinue returns a transaction that can be continued by suppling the
// TxID that gets passed into the txFunc.
func TransactionContinue(TxID string, txFunc func(*Tx, string) error) (err error) {
	return stdDB.TransactionContinue(TxID, txFunc)
}

// TransactionContinue returns a transaction that can be continued by suppling the
// TxID that gets passed into the txFunc.
func (db DB) TransactionContinue(TxID string, txFunc func(*Tx, string) error) (err error) {
	var tx *Tx

	if TxID == "" {

		TxID = uuid.V4()
		tx, err = db.NewTx(context.Background(), false)
		if err != nil {
			log.Error(err.Error())
			return
		}
		txSet(TxID, tx)

	} else {
		if tx = txGet(TxID); tx != nil {
			err = txFunc(tx, TxID)
			return err
		}
		return fmt.Errorf("unable to find tx %s", TxID)
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}

		if err != nil {
			tx.Rollback()
			log.Error(err.Error())

			debug.PrintStack()
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error(err.Error())
		}

		txRm(TxID)
	}()

	err = txFunc(tx, TxID)
	return err
}

var txMap = make(map[string]*Tx)
var txMutex = sync.Mutex{}

func txSet(id string, tx *Tx) {
	txMutex.Lock()
	defer txMutex.Unlock()
	txMap[id] = tx
}
func txGet(id string) *Tx {
	txMutex.Lock()
	defer txMutex.Unlock()
	return txMap[id]
}
func txRm(id string) {
	txMutex.Lock()
	defer txMutex.Unlock()
	delete(txMap, id)
}
