// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"
	"time"
)

type Querier interface {
	GetDailyTransactionReport(ctx context.Context, createdAt time.Time) (GetDailyTransactionReportRow, error)
	SaveDailyTransactionReport(ctx context.Context, arg SaveDailyTransactionReportParams) error
}

var _ Querier = (*Queries)(nil)
