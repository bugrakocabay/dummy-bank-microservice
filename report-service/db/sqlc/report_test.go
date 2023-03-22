package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSaveDailyTransactionReport(t *testing.T) {
	todayDate := time.Now().Format("02.01.2006")
	arg := SaveDailyTransactionReportParams{
		AvgTransactionAmount:   float64(10),
		TotalTransactionAmount: int32(10),
		NumTransactions:        int32(10),
		TotalCommission:        float64(10),
		Day:                    todayDate,
	}
	err := testQueries.SaveDailyTransactionReport(context.Background(), arg)
	require.NoError(t, err)

	report, err := testQueries.GetDailyTransactionReport(context.Background(), time.Now())
	require.NoError(t, err)

	fmt.Println(arg)
	fmt.Println(report)

	require.Equal(t, report.AvgTransactionAmount, arg.AvgTransactionAmount)
	require.Equal(t, report.TotalTransactionAmount, arg.TotalTransactionAmount)
	require.Equal(t, report.TotalCommission, arg.TotalCommission)
	require.Equal(t, report.NumTransactions, arg.NumTransactions)
	require.Equal(t, report.Day, arg.Day)
}
