package stream

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/amirrezaask/plumber"
	"io"
	"log"
)

type SQLifier func(i interface{}) (string, []interface{})
type SQLOutputState struct{}
type SQLOutput struct {
	db        *sql.DB
	sqilifier SQLifier
	writeChan chan interface{}
}

func NewSQLOutput(db *sql.DB, f SQLifier) (plumber.Output, error) {
	s := &SQLOutput{db: db, sqilifier: f, writeChan: make(chan interface{})}
	go func() {
		for m := range s.writeChan {
			fmt.Println(m)
			stmt, values := f(m)
			_, err := db.Exec(stmt, values...)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}()
	return s, nil
}

func (S *SQLOutput) LoadState(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&SQLOutput{})
}

func (S *SQLOutput) State() ([]byte, error) {
	return nil, nil
}

func (S *SQLOutput) Name() string {
	return "sql-output"
}

func (S *SQLOutput) Output() (chan interface{}, error) {
	return S.writeChan, nil
}
