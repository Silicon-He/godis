package reply

// Using the same instance globally reduces memory overhead
var (
	PongReplyInstance           = NewPongReply()
	OkReplyInstance             = NewOkReply()
	NullBulkReplyInstance       = NewNullBulkReply()
	EmptyMultiBulkReplyInstance = NewEmptyMultiBulkReply()
	NoReplyInstance             = NewNoReply()
	QueuedReplyInstance         = NewQueuedReply()
)

// PongReply reply a "PONG"
type PongReply struct{}

var pongBytes = []byte("+PONG\r\n")

func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

func NewPongReply() *PongReply {
	return &PongReply{}
}

// OkReply reply an "OK"
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

func (o *OkReply) ToBytes() []byte {
	return okBytes
}

func NewOkReply() *OkReply {
	return &OkReply{}
}

// NullBulkReply reply an empty string
type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func NewNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// EmptyMultiBulkReply reply an empty multi bulk (数组)
type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

// NoReply reply a NULL
type NoReply struct{}

var noBytes = []byte("")

func (n *NoReply) ToBytes() []byte {
	return noBytes
}

func NewNoReply() *NoReply {
	return &NoReply{}
}

type QueuedReply struct{}

func (q *QueuedReply) ToBytes() []byte {
	return queuedBytes
}

func NewQueuedReply() *QueuedReply {
	return &QueuedReply{}
}

var queuedBytes = []byte("+QUEUED\r\n")
