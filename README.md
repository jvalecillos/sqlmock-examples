# Examples of go-sqlmock

**DATA-DOG/go-sqlmock** is a mock library implementing [sql/driver](https://godoc.org/database/sql/driver). Which has one and only
purpose - to simulate any **sql** driver behavior in tests, without needing a real database connection. It helps to
maintain correct **TDD** workflow.

There are some basic examples of using **DATA-DOG/go-sqlmock** in the project repository. Nevertheless while using the library on different projects I've found myself trying to discover how to write certain tests, thus I developed several tricks that I want to share in this repo.

## Using sqlmock for table driven tests

The basic example have a positive path and a negative path and one test for each. What does it happen if you want to effectively cover more scenarios?, well you would need to create a different mock and a different dummy connection every time. That's not optimal since the idea of table-driven-tests is to make unit tests structured, and to make easy to add different cases.

You could solve this issue creating a dummy connection and a different mock for each test case. Then, each test could have a closure that set the expecations and behaviour for each case:

```go
func Test_example(t *testing.T) {

	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "my test case",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO test \(key, value\) VALUES \(1,1\)`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// db and mock for current iteration
			db, mock, _ := sqlmock.New()
			defer db.Close()
			// set mock expectations
			tt.mockClosure(mock)
						
			// ...

			// Checking all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
```

Better example [basic_test.go](basic_test.go)