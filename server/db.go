package main

import (
	"database/sql"
	"strings"
	_ "github.com/mattn/go-sqlite3"
	_ "embed"
)

//go:embed schema.sql
var databaseschema string

type sqliteStore struct {
	db *sql.DB
	tagsCache map[string]int
}

func (s *sqliteStore) index(filename string, note *Note) {
	uuid := strings.CutSufix(filename, ".md")
	if s.isIndexed(uuid) {
		_ = s.updateIndex(uuid, note)
	}
	_ = s.createIndex(uuid, note)
}
func (s *sqliteStore) updateIndex(filenameuuid string, note *Note) error {
}
func (s *sqliteStore) createIndex(filenameuuid string, note *Note) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	// Notes table
	query := `insert into notes (filename, title) value (?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec(filenameuuid, note.Header.Title)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Tags
	
	return nil
}
func (s *sqliteStore) isIndexed(filenameuuid string) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	query := `select * from notes where filename = ?`
	row, err := tx.Query(query, filename)
	if err != nil {
		return false, err
	}
	n := struct {
		filename string
		title string
	}{}
	for row.Next() {
		row.Scan(
			&n.filename,
			&n.title,
		)
		return n.filename == filename, nil
	}
	return false, nil
}

func NewSqliteStore(dbpath string) (*sqliteStore, error) {
	if _, err := os.Stat(dbpath); errors.Is(err, os.ErrNotExist) {
		createDatabase(dbpath)
	}
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &sqliteStore{
		db: db,
		tags: make(map[string]int)
	}, nil
}

func createDatabase(dbpath string) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=rw", dbpath))
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(databaseschema)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("created new database at: %s", dbpath)
	db.Close()
}
