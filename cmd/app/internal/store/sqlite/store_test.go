package sqlitestore_test

import (
	"os"
	"testing"
)

var (
	dsn string
)

func TestMain(m *testing.M) {

	dsn = os.Getenv("GOMTC_DSN")
	if dsn == "" {
		dsn = "sqlite_test.db"
	}

	os.Exit(m.Run())
}
