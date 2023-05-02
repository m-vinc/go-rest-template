package ent

import (
	"context"
	"fmt"
)

func NewTx(ctx context.Context, c *Client) (*Tx, context.Context, error) {
	var err error
	tx := TxFromContext(ctx)
	if tx == nil {
		tx, err = c.Tx(ctx)
		if err != nil {
			return nil, nil, err
		}
		return tx, NewTxContext(ctx, tx), nil
	}

	return tx, ctx, nil
}

func Commit(ctx context.Context, tx *Tx) error {
	ctxTx := TxFromContext(ctx)
	if ctxTx != nil {
		return nil
	}

	return tx.Commit()
}

func Rollback(ctx context.Context, tx *Tx, err error) error {
	ctxTx := TxFromContext(ctx)
	if ctxTx != nil {
		return err
	}

	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
