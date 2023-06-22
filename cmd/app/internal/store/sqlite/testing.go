package sqlitestore

import (
	"database/sql"
	"io/fs"
	"os"
	"path"
	"sort"

	"fmt"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestDB(t *testing.T, dsn string) (string, func(...string)) {
	t.Helper()

	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	if MigrateSQLite("../../../migrations/sqlite", conn) != nil {
		conn.Close()
		fmt.Println(err)
		t.Fatal(err)
	}

	return dsn, func(tables ...string) {
		if len(tables) > 0 {
			for _, table := range tables {
				_, err := conn.Exec(fmt.Sprintf("DELETE FROM %s", table))
				if err != nil {
					fmt.Println(err)
					t.Fatal(err)
				}
			}
		}
		conn.Close()
	}
}

func MigrateSQLite(pathToMigrations string, db *sql.DB) error {

	files, err := os.ReadDir(pathToMigrations)
	if err != nil {
		return err
	}
	SortFileByName(files)
	for _, file := range files {
		if strings.Contains(file.Name(), ".up.") {
			b, err := os.ReadFile(path.Join(pathToMigrations, file.Name()))
			if err != nil {
				return err
			}
			_, err = db.Exec(string(b))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SortFileByName(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.Name() < infoJ.Name()
	})
}
