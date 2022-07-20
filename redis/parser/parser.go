package parser

import (
	"bufio"
	"bytes"
	"errors"
	"godis/interface/resp"
	"godis/lib/logger"
	"godis/redis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
)

/*
simple string:  +xxx\r\n  e.g. "+OK\r\n"
error string:  -xxx\r\n  e.g. "-ERR Invalid Synatx\r\n"
integer:   :n\r\n   e.g. ":1\r\n"
bulk string(binary safe): $(length)\r\nxxxx\r\n   e.g. "$3\r\nSET\r\n"
multi-bulk string(array): startwith *
e.g.  *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
means: ['SET', 'key', 'value']

$-1: nil
*/

type Payload struct {
	Data resp.Reply
	Err  error
}

// ParseStream stream interface usually offer to client/server
func ParseStream(reader io.Reader) chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

// ParseOne parse a byte slice
func ParseOne(data []byte) (reply resp.Reply, err error) {
	ch := make(chan *Payload)
	reader := bytes.NewReader(data)
	go parse0(reader, ch)
	payload := <-ch
	if payload == nil {
		return nil, errors.New("no reply")
	}
	return payload.Data, nil
}

// parse0 core of parser
func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err, string(debug.Stack()))
		}
	}()
	// init reader state
	var state readState
	var err error
	var line []byte
	var ioErr bool
	bufReader := bufio.NewReader(reader)
	for {
		line, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr {
				// encounter io err ,stop read
				ch <- &Payload{Err: err}
				close(ch)
				return
			}
			// protocol err , reset read state
			ch <- &Payload{Err: err}
			state = readState{}
			continue
		}

		// start parse a line
		if !state.readingMultiLine {
			switch line[0] {
			case '*':
				// multi-bulk string (array) protocol
				err = parseMultiBulkHeader(line, &state)
				if err != nil {
					ch <- &Payload{Err: err}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{Data: reply.EmptyMultiBulkReplyInstance}
					state = readState{}
				}
			case '$':
				// bulk string protocol
				err = parseBulkHeader(line, &state)
				if err != nil {
					ch <- &Payload{Err: err}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 {
					ch <- &Payload{Data: reply.NullBulkReplyInstance}
				}
				state = readState{}
				continue
			}
		} else {
			// todo impl this
			parseSingleLine()
		}
	}

}

type readState struct {
	// multi-bulk/bulk string set true
	readingMultiLine bool
	// count multi-bulk words nums
	expectedArgsCount int
	// '+'/'-'/'*'/'$'/':'
	msgType byte
	// exact words
	args [][]byte
	// next word's length
	bulkLen     int64
	readingRepl bool
}

func (s *readState) isFinished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

func readLine(bufReader *bufio.Reader, state *readState) (line []byte, ioErr bool, err error) {
	if state.bulkLen == 0 {
		// read simple string
		line, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(line) == 0 || line[len(line)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(line))
		}
	} else {
		// read a bulk string has CRLF
		// but there is no CRLF in RDB and AOF
		bulkLen := state.bulkLen + 2
		if state.readingRepl {
			bulkLen -= 2
		}
		line = make([]byte, state.bulkLen)
		_, err = io.ReadFull(bufReader, line)
		if err != nil {
			return nil, true, err
		}
		state.bulkLen = 0
	}
	return line, false, nil
}

// parseMultiBulkHeader parse multi-bulk string's header
func parseMultiBulkHeader(line []byte, state *readState) (err error) {
	var expectedLine uint64
	// e.g.  "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	// the line is the string head: *3\r\n , so the array length is [1:-2]
	expectedLine, err = strconv.ParseUint(string(line[1:len(line)-2]), 10, 64)
	if err != nil {
		return err
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		// first line of multi-bulk protocol
		state.msgType = line[0]
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
		return nil
	} else {
		return errors.New("protocol error: " + string(line))
	}
}

// parseBulkHeader parse bulk string's header
func parseBulkHeader(line []byte, state *readState) (err error) {
	// e.g. "$3\r\nSET\r\n"
	state.bulkLen, err = strconv.ParseInt(string(line[1:len(line)-2]), 10, 64)
	if err != nil {
		return err
	}
	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.args = make([][]byte, 0, 1)
		state.msgType = line[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		return nil
	} else {
		return errors.New("protocol error: " + string(line))
	}
}

// todo impl this
func parseSingleLine() {

}