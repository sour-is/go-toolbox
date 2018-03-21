package dbm

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"runtime/debug"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/uuid"
	"sync"

	sq "github.com/Masterminds/squirrel"
)

type Tx struct {
	*sql.Tx
	DbType      string
	Placeholder sq.PlaceholderFormat
	Returns     bool
}

func NewTx(db *sql.DB, dbType string, placeholder sq.PlaceholderFormat, returns bool) (tx *Tx, err error) {
	tx = new(Tx)
	tx.Tx, err = db.Begin()
	tx.Placeholder = placeholder
	tx.DbType = dbType
	tx.Returns = returns

	return
}

// Transaction starts a new database transction and executes the supplied func.
func Transaction(txFunc func(*Tx) error) (err error) {
	tx, err := NewTx(db, dbType, placeholder, returns)

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
func Transactionx(txFunc func(*sqlx.Tx) error) (err error) {
	dbx := sqlx.NewDb(db, dbType)

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

// TransactionContinue returns a transaction that can be continued by suppling the
// TxID that gets passed into the txFunc.
func TransactionContinue(TxID string, txFunc func(*Tx, string) error) (err error) {
	var tx *Tx

	if TxID == "" {

		TxID = uuid.V4()
		tx, err = NewTx(db, dbType, placeholder, returns)
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