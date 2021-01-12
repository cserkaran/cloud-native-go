package logger

import (
	"bufio"
	"fmt"
	"os"
)

// FileTransactionLogger defines the File logger.
type FileTransactionLogger struct {
	events       chan<- Event // Write only channel for sending events.
	errors       <-chan error // Read only channel for receiving errors.
	lastSequence uint64       // The last used event sequence number.
	file         *os.File     // The location of transaction log.
}

// NewFileTransactionLogger creates new FileTransactionLogger
func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log")
	}
	return &FileTransactionLogger{file: file}, nil
}

// WritePut writes PUT event in the log.
func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

// WriteDelete writes DELETE event in the log.
func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

// Err returns errors channel to commmunicate errors.
func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

// Run the FileTransaction logger.
// Reads the events written and writes them to file in separate goroutine.
// Writes errors to error channel of in case.
func (l *FileTransactionLogger) Run(){
    events := make(chan Event, 16)
    l.events = events

    errors := make(chan error,1)
    l.errors = errors

    go func ()  {
        for e := range events{
            l.lastSequence++
            _,err := fmt.Fprintf(
                l.file,
                "%d\t%d\t%s\t%s\n",
                l.lastSequence, e.EventType,e.Key,e.Value)
            if err != nil{
                errors <- err
                return
            }
        }
    }()
}

// ReadEvents reads from transaction logs  and replays the event into the store
func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
    scanner := bufio.NewScanner(l.file)  // Create a Scanner for l.file.
    outEvent := make(chan Event) // An unbuffered events channel.
    outError := make(chan error,1) // A buffered errors channel.

    go func(){
        var e Event
        
        defer close(outEvent)
        defer close(outError)

        for scanner.Scan(){
            line := scanner.Text()

            fmt.Sscanf(
                line, "%d\t%d\t%s\t%s\t",
                &e.Sequence, &e.EventType, &e.Key, &e.Value)

            // Sanity check! Are the sequence numbers in increasing order?
            if l.lastSequence >= e.Sequence{
                outError <- fmt.Errorf("transaction numbers out of sequence")
            }
            
            l.lastSequence = e.Sequence
            outEvent <- e
        }

        if err := scanner.Err(); err != nil{
            outError <- fmt.Errorf("transaction log read failure: %w", err)
        }
    }()

    return outEvent,outError
}
