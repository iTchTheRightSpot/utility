package utils

import (
	"context"
	"database/sql"
	"time"
)

type RequestContextKey string

type RequestBody struct {
	Id     string `json:"request_id,omitempty"`
	Ip     string `json:"ip_address,omitempty"`
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
}

const RequestKey RequestContextKey = "REQUEST_KEY"

func Timezone(timezone string) (*time.Location, error) {
	if timezone == "" {
		timezone = "UTC"
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}
	return location, nil
}

func RunInTx(ctx context.Context, l ILogger, db *sql.DB, fn func(*sql.Tx) error) error {
	l.Log(ctx, "STARTING TRANSACTION")

	tx, err := db.Begin()
	if err != nil {
		l.Critical(ctx, "FAILED TO BEGIN TRANSACTION: "+err.Error())
		return &ServerError{}
	}

	err = fn(tx)
	if err == nil {
		l.Log(ctx, "COMMITTING TRANSACTION")

		if commitErr := tx.Commit(); commitErr != nil {
			l.Critical(ctx, "FAILED TO COMMIT TRANSACTION: "+commitErr.Error())
			return &ServerError{}
		}

		l.Log(ctx, "TRANSACTION COMMITTED SUCCESSFULLY")
		return nil
	}

	l.Error(ctx, "TRANSACTION FUNCTION RETURNED ERROR: "+err.Error())
	l.Log(ctx, "ROLLING BACK TRANSACTION")

	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		l.Critical(ctx, "FAILED TO ROLLBACK TRANSACTION: "+rollbackErr.Error())
		return &ServerError{}
	}

	l.Log(ctx, "TRANSACTION ROLLED BACK SUCCESSFULLY")
	return err
}