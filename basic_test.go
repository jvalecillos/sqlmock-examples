package main

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
)

func Test_recordStats(t *testing.T) {

	type args struct {
		userID    int64
		productID int64
	}

	testArgs := args{userID: 2, productID: 3}
	updateQueryRgx := `UPDATE products SET views = views \+ 1`
	insertQueryRgx := `INSERT INTO product_viewers \(user_id, product_id\) VALUES \(\?, \?\)`

	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		args        args
		wantErr     bool
	}{
		{
			name: "success case",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(updateQueryRgx).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(insertQueryRgx).
					WithArgs(testArgs.userID, testArgs.productID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			args:    testArgs,
			wantErr: false,
		},
		{
			name: "failure on begin",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("some error"))
			},
			args:    testArgs,
			wantErr: true,
		},
		{
			name: "failure on update",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(updateQueryRgx).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			args:    testArgs,
			wantErr: true,
		},
		{
			name: "failure on insert",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(updateQueryRgx).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(insertQueryRgx).
					WithArgs(testArgs.userID, testArgs.productID).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			args:    testArgs,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Running tests in parallel :)
			t.Parallel()
			// db and mock for current iteration
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			// Set mock expectations
			tt.mockClosure(mock)

			// Testing function
			if err := recordStats(db, tt.args.userID, tt.args.productID); (err != nil) != tt.wantErr {
				t.Errorf("recordStats() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Closing mock connection
			mock.ExpectClose()
			// Explicit closing instead of deferred in order to check ExpectationsWereMet
			if err = db.Close(); err != nil {
				t.Error(err)
			}

			// Checking all expectations were met
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

		})
	}
}
