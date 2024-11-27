package api

import (
	"context"
	"database/sql"
	"errors"
	"maple/internal/schema"
)

func (a *MapleAPI) Transaction(ctx context.Context, fn func(tx *sql.Tx, queries *schema.Queries) error) error {
	tx, err := a.Conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	fnError := fn(tx, a.Queries.WithTx(tx))
	if fnError != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return errors.New("error rolling back transaction: " + txErr.Error())
		}
	} else {
		if txErr := tx.Commit(); txErr != nil {
			return errors.New("error commiting back transaction: " + txErr.Error())
		}
	}

	return fnError
}
