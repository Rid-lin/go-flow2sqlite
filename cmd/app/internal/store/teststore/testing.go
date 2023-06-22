package teststore

import (
	"database/sql"
	"io/fs"
	"os"
	"path"
	"sort"

	"strings"
	"testing"
)

func TestDB(t *testing.T, dsn string) (string, func(...string)) {
	t.Helper()

	return dsn, func(tables ...string) {

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
