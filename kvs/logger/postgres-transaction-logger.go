package logger

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Anonymously import the driver package
)

// PostgresDbParams parameters required for Postgres Database connection.
type PostgresDbParams struct {
	DbName   string
	Host     string
}

// PostgresTransactionLogger defines the Database transaction logger.
type PostgresTransactionLogger struct {
	events chan<- Event // Write-only channel for sending events
	errors <-chan error // Read-only channel for receiving errors
	db     *sql.DB      // Our database access interface
}

// NewPostgreTransactionLogger creates a new Database transaction logger.
func NewPostgreTransactionLogger(config PostgresDbParams) (TransactionLogger, error) {

	connStr := fmt.Sprintf("host=%s dbname=%s sslmode=disable",
        config.Host, config.DbName)
        
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to created db value :%w", err)
	}

	err = db.Ping() // Test the database connection.

	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	tl := &PostgresTransactionLogger{db: db}

	exists, err := tl.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}

	if !exists {
		if err = tl.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

    fmt.Println("transaction logger created successfully")
	return tl, nil

}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	query := `SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE  table_schema = 'public'
        AND    table_name   = 'transactions'
        )`

	rows, err := l.db.Query(query)
	if err != nil {
		return false, err
	}

	var exists bool
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&exists)

		if err != nil {
			return false, err
		}
	}

	err = rows.Err()
	if err != nil {
		return false, err
	}

	return exists, nil

}

func (l *PostgresTransactionLogger) createTable() error {
	query := `create table postgres.public.transactions (
        sequence serial,
        event_type int,
        key varchar,
        value varchar
   )`

	_, err := l.db.Exec(query)
	return err
}

// WritePut writes PUT event in the log.
func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

// WriteDelete writes DELETE event in the log.
func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

// Err returns errors channel to commmunicate errors.
func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

// Run the PostgresTransactionLogger.
func (l *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions
        (event_type,key,value)
        VALUES ($1,$2,$3)`

		for e := range events {
			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}

// ReadEvents reads from events database transactions tables
// and replays the event into the store.
func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event) // An unbuffered events channel
	outError := make(chan error) // A buffered errors channel

	go func() {
		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		query := "SELECT sequence,event_type,key,value FROM transactions"
		rows, err := l.db.Query(query) // Run query: get result
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close()

		e := Event{}

		for rows.Next() {
			err = rows.Scan(
				&e.Sequence, &e.EventType, &e.Key, &e.Value)

			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
				return
			}

			outEvent <- e
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}
