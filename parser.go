// Package ratlogparser provides parsing for the ratlog logging format
// see https://github.com/ratlog/ratlog-spec
package ratlogparser

import (
	"fmt"
	"io"
)

// EntryReader is an interface which allows the reading of
// entries
type EntryReader interface {
	Entry() (Entry, error)
}

// EntryWriter is an interface which allows the writing of
// entries
type EntryWriter interface {
	Write(Entry) error
}

// Entry is a Ratlog log entry
// see https://github.com/ratlog/ratlog-spec
type Entry struct {
	Tags    []Tag
	Message string
	Fields  Fields
}

// Tag a tag associated with an entry
type Tag string

// Fields holds a map of fields. A Field is anything that
// implements the fmt.Stringer interface. The key is a string.
// https://github.com/ratlog/ratlog-spec#fields
type Fields map[string]fmt.Stringer

// Find a field
func (f Fields) Find(key string) fmt.Stringer {
	return f[key]
}

// BasicField is a string that implements fmt.Stringer
type BasicField string

// Convert a basic field to it's string representation
func (b BasicField) String() string {
	return string(b)
}

// EntryReaderWriter with the underlying slice buffered in memory
// the underlying structure is FIFO.
type EntryReaderWriter struct {
	entries []Entry
}

// Entry returns the first entry from the stack. If the entries are nil return
// an error. This will happen if trying to read the same buffer multiple
// times.
func (b *EntryReaderWriter) Entry() (Entry, error) {
	if len(b.entries) == 0 {
		return Entry{}, io.EOF
	}
	if len(b.entries) == 1 {
		e := b.entries[0]
		b.entries = nil
		return e, nil
	}
	e, newEntries := b.entries[0], b.entries[1:]
	b.entries = newEntries
	return e, nil
}

// Entries returns all the entries as a slice
func (b *EntryReaderWriter) Entries() ([]Entry, error) {
	if b.entries == nil {
		return nil, fmt.Errorf("empty buffer")
	}
	return b.entries, nil
}

// Write writes an entry to the underlying buffer
func (b *EntryReaderWriter) Write(e Entry) error {
	if b.entries == nil {
		entries := make([]Entry, 1)
		entries[0] = e
		b.entries = entries
	}
	b.entries = append(b.entries, e)
	return nil
}

// NewEntryReader creates a new reader from a slice of entries
// copies entries so the original slice is safe to reuse
func NewEntryReaderWriter(e []Entry) *EntryReaderWriter {
	b := EntryReaderWriter{}
	be := make([]Entry, len(e))
	copy(be, e)
	b.entries = be
	return &b
}
