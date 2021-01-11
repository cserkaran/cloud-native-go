package logger

// EventType is a constant which defines the action taken.
type EventType byte

const (
	_                     = iota // iota == 0; ignore the zero value.
    // EventDelete for action DELETE.
    EventDelete EventType = iota // iota = 1.
    // EventPut for action PUT
	EventPut                     // iota == 2; implicitly repeat.
)

// Event Record which defines an entry in the transaction log.
type Event struct {
	Sequence  uint64    // A unique record ID.
	EventType EventType // The action taken.
	Key       string    // The key affected by this transaction.
	Value     string    // The value of a PUT the transaction.
}
