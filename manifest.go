package go_streams

// Entry is the data model that go-streams passes between
// different operators although the user never handle it directly
// when creating new streams.
type Entry struct {
	// Key is a unique identifier that represents
	// this entry, useful for sources that
	// needs to commit its progress.
	Key string

	// Value is the actual model that the user transforms
	Value interface{}

	// Filtered indicates that this row should be filtered the filter,
	// in order to decrease array mutation operations.
	Filtered bool
}

// MapFunc is a function which transforms its input
type MapFunc func(entry interface{}) interface{}

// FilterFunc is a function that takes an entry an decided
// if this entry should be filtered out
// return true to keep the record or false to filter it out.
type FilterFunc func(entry interface{}) bool

// ErrorHandler is a function that takes an error
// useful when you want to handle errors yourself.
type ErrorHandler func(err error)

// EntryChannel is a type alias, used internally
type EntryChannel chan Entry

// ErrorChannel is a type alias, used internally
type ErrorChannel chan error

// Stream is an interface that represents all the operations
// that you can use for a given stream.
type Stream interface {
	// Filter entries
	Filter(fn FilterFunc) Stream

	// Map entries
	Map(fn MapFunc) Stream

	// Sink (or dump) the stream entries to this Sink implementation
	// such as file, database, memory, etc...
	Sink(sink Sink) Stream

	// Process takes a processor implementation and an error channel
	// and will start processing the stream.
	Process(processor Processor, errs ErrorChannel)

	// Will return the handlers (filters, maps, sinks) associated with this stream.
	GetHandlers() []interface{}

	// Will return the source of the stream.
	GetSource() Source
}

// Processor is responsible of the processing strategy,
// for example: a DirectProcessor will stream each message sent from source
// to the whole pipeline and a BufferedProcessor will buffer messages before moving them
// throughout the pipeline.
type Processor interface {
	Process(stream Stream, errs ErrorChannel)
}

// ProcessorFactory is a type alias to a function which generates processors
type ProcessorFactory func() Processor

// Source is responsible for sending entries into streams.
type Source interface {
	// Start will notify the source to start sending new entries to the EntryChannel
	// Start will be called with a new go-routine.
	Start(channel EntryChannel, errorChannel ErrorChannel)

	// Stop will stop the source from sending new entries.
	// NOTICE that stopped source cannot be started again.
	Stop() error

	// Commit entry will be called by the stream sinks
	// its an optional method and if you don't need to manage idempotence you can
	// leave it to 'return nil'
	CommitEntry(key string) error

	// Name should return a unique name for the given source
	Name() string
}

// Sink is responsible for dumping entries into sinks such as: files, databases, memory, etc...
type Sink interface {
	// Single will dump a single entry to the sink
	Single(entry Entry) error

	// Batch will dump a batch of entries to the sink
	Batch(entry ...Entry) error
}

// Engine is responsible for managing one or more streams
// it allows you to start/stop groups of streams and provide
// central error handling for your streams.
type Engine interface {
	// Add new stream, NOTICE that streams with the same source cannot be added.
	Add(stream ...Stream) error

	// Sets an error handler that will be called whenever an error is reported.
	SetErrorHandler(handler ErrorHandler)

	// Will start all attached streams
	Start()

	// Will start all stop streams
	Stop()
}