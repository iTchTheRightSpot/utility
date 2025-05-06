package utils

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

type helper struct {
	Called bool
	Error  error
}

func (d *helper) fn(_ *sql.Tx) error {
	d.Called = true
	return d.Error
}

func TestTransactionFunc(t *testing.T) {
	t.Parallel()

	lg := DevLogger("UTC")

	t.Run("rollback. error starting transaction", func(t *testing.T) {
		t.Parallel()

		// given
		h := helper{}
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err.Error())
		}

		defer func(db *sql.DB) {
			if err = db.Close(); err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
		}(db)

		mock.ExpectClose()

		// method to test & assert
		if err = RunInTx(context.Background(), lg, db, h.fn); err == nil {
			t.Error("expect error")
			t.FailNow()
		}

		mess := "server error"
		if err.Error() != mess {
			t.Errorf("expect %s, given %s", mess, err.Error())
			t.FailNow()
		}

		if h.Called {
			t.Error("transaction func called but should not")
			t.FailNow()
		}
	})

	t.Run("commit tx. fn does not return error", func(t *testing.T) {
		t.Parallel()

		// given
		h := helper{}
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err.Error())
		}

		defer func(db *sql.DB) {
			if err = db.Close(); err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
		}(db)

		mock.ExpectBegin()
		mock.ExpectCommit()
		mock.ExpectClose()

		// method to test & assert
		if err = RunInTx(context.Background(), lg, db, h.fn); err != nil {
			t.Errorf("expect nil, given %s", err.Error())
			t.FailNow()
		}

		if !h.Called {
			t.Error("expect fn true, given false")
			t.FailNow()
		}
	})

	t.Run("rollback. fn does not return error but commit returns error", func(t *testing.T) {
		t.Parallel()

		// given
		h := helper{}
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err.Error())
		}

		defer func(db *sql.DB) {
			if err = db.Close(); err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
		}(db)

		mock.ExpectBegin()
		mock.ExpectClose()

		// method to test & assert
		if err = RunInTx(context.Background(), lg, db, h.fn); err == nil {
			t.Error("expect error, given nil")
			t.FailNow()
		}

		if !h.Called {
			t.Error("expect fn true, given false")
			t.FailNow()
		}

		mess := "server error"
		if err.Error() != mess {
			t.Errorf("expect %s, given %s", mess, err.Error())
			t.FailNow()
		}
	})

	t.Run("rollback. fn returns error", func(t *testing.T) {
		t.Parallel()

		// given
		str := "fn returns error"
		h := helper{Error: errors.New(str)}
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err.Error())
		}

		defer func(db *sql.DB) {
			if err = db.Close(); err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
		}(db)

		mock.ExpectBegin()
		mock.ExpectRollback()
		mock.ExpectClose()

		// method to test & assert
		if err = RunInTx(context.Background(), lg, db, h.fn); err == nil {
			t.Error("expect error, given nil")
			t.FailNow()
		}

		if !h.Called {
			t.Error("expect fn true, given false")
			t.FailNow()
		}

		if err.Error() != str {
			t.Errorf("expect %s, given %s", str, err.Error())
			t.FailNow()
		}
	})

	t.Run("rollback. fn returns error and rollback returns error", func(t *testing.T) {
		t.Parallel()

		// given
		h := helper{Error: errors.New("fn returns error")}
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err.Error())
		}

		defer func(db *sql.DB) {
			if err = db.Close(); err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
		}(db)

		mock.ExpectBegin()
		mock.ExpectClose()

		// method to test & assert
		if err = RunInTx(context.Background(), lg, db, h.fn); err == nil {
			t.Error("expect error, given nil")
			t.FailNow()
		}

		if !h.Called {
			t.Error("expect fn true, given false")
			t.FailNow()
		}

		mess := "server error"
		if err.Error() != mess {
			t.Errorf("expect %s, given %s", mess, err.Error())
			t.FailNow()
		}
	})
}
