package logger

// TransactionLogger interface for logging transactions done
// on the map store.
type TransactionLogger interface {
	WriteDelete(key string)
    WritePut(key, value string)
    Err() <-chan error
    ReadEvents() (<-chan Event, <-chan error)
    Run()
}
