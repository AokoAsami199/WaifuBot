// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addOneToClaimCountStmt, err = db.PrepareContext(ctx, addOneToClaimCount); err != nil {
		return nil, fmt.Errorf("error preparing query AddOneToClaimCount: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.deleteCharStmt, err = db.PrepareContext(ctx, deleteChar); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteChar: %w", err)
	}
	if q.getCharStmt, err = db.PrepareContext(ctx, getChar); err != nil {
		return nil, fmt.Errorf("error preparing query GetChar: %w", err)
	}
	if q.getDateStmt, err = db.PrepareContext(ctx, getDate); err != nil {
		return nil, fmt.Errorf("error preparing query GetDate: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.getUserListStmt, err = db.PrepareContext(ctx, getUserList); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserList: %w", err)
	}
	if q.getUserProfileStmt, err = db.PrepareContext(ctx, getUserProfile); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserProfile: %w", err)
	}
	if q.giveCharStmt, err = db.PrepareContext(ctx, giveChar); err != nil {
		return nil, fmt.Errorf("error preparing query GiveChar: %w", err)
	}
	if q.insertCharStmt, err = db.PrepareContext(ctx, insertChar); err != nil {
		return nil, fmt.Errorf("error preparing query InsertChar: %w", err)
	}
	if q.setFavoriteStmt, err = db.PrepareContext(ctx, setFavorite); err != nil {
		return nil, fmt.Errorf("error preparing query SetFavorite: %w", err)
	}
	if q.setQuoteStmt, err = db.PrepareContext(ctx, setQuote); err != nil {
		return nil, fmt.Errorf("error preparing query SetQuote: %w", err)
	}
	if q.updateUserDateStmt, err = db.PrepareContext(ctx, updateUserDate); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserDate: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addOneToClaimCountStmt != nil {
		if cerr := q.addOneToClaimCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addOneToClaimCountStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.deleteCharStmt != nil {
		if cerr := q.deleteCharStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteCharStmt: %w", cerr)
		}
	}
	if q.getCharStmt != nil {
		if cerr := q.getCharStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCharStmt: %w", cerr)
		}
	}
	if q.getDateStmt != nil {
		if cerr := q.getDateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDateStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.getUserListStmt != nil {
		if cerr := q.getUserListStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserListStmt: %w", cerr)
		}
	}
	if q.getUserProfileStmt != nil {
		if cerr := q.getUserProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserProfileStmt: %w", cerr)
		}
	}
	if q.giveCharStmt != nil {
		if cerr := q.giveCharStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing giveCharStmt: %w", cerr)
		}
	}
	if q.insertCharStmt != nil {
		if cerr := q.insertCharStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertCharStmt: %w", cerr)
		}
	}
	if q.setFavoriteStmt != nil {
		if cerr := q.setFavoriteStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setFavoriteStmt: %w", cerr)
		}
	}
	if q.setQuoteStmt != nil {
		if cerr := q.setQuoteStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setQuoteStmt: %w", cerr)
		}
	}
	if q.updateUserDateStmt != nil {
		if cerr := q.updateUserDateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserDateStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                     DBTX
	tx                     *sql.Tx
	addOneToClaimCountStmt *sql.Stmt
	createUserStmt         *sql.Stmt
	deleteCharStmt         *sql.Stmt
	getCharStmt            *sql.Stmt
	getDateStmt            *sql.Stmt
	getUserStmt            *sql.Stmt
	getUserListStmt        *sql.Stmt
	getUserProfileStmt     *sql.Stmt
	giveCharStmt           *sql.Stmt
	insertCharStmt         *sql.Stmt
	setFavoriteStmt        *sql.Stmt
	setQuoteStmt           *sql.Stmt
	updateUserDateStmt     *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                     tx,
		tx:                     tx,
		addOneToClaimCountStmt: q.addOneToClaimCountStmt,
		createUserStmt:         q.createUserStmt,
		deleteCharStmt:         q.deleteCharStmt,
		getCharStmt:            q.getCharStmt,
		getDateStmt:            q.getDateStmt,
		getUserStmt:            q.getUserStmt,
		getUserListStmt:        q.getUserListStmt,
		getUserProfileStmt:     q.getUserProfileStmt,
		giveCharStmt:           q.giveCharStmt,
		insertCharStmt:         q.insertCharStmt,
		setFavoriteStmt:        q.setFavoriteStmt,
		setQuoteStmt:           q.setQuoteStmt,
		updateUserDateStmt:     q.updateUserDateStmt,
	}
}
