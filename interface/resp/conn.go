package resp

// Connection instantiation at the protocol layer
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
