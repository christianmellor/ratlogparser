package ratlogparser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// SimpleParser which returns the Entry results a EntryReaderWriter
// with the underlying data stored in a slice
type SimpleParser struct{}

// Parse an io.Reader into a EntryReaderWriter
func (p SimpleParser) Parse(r io.Reader, w EntryWriter) error {
	scanner := bufio.NewScanner(r)
	buffer := bytes.NewBuffer(nil)
	for scanner.Scan() {
		buffer.Write(scanner.Bytes())
		e, err := p.parseLine(buffer)
		if err != nil {
			return err
		}
		err = w.Write(e)
		if err != nil {
			return err
		}
		buffer.Reset()
	}
	return nil
}
func (p SimpleParser) parseLine(line io.RuneReader) (Entry, error) {
	entry := Entry{}
	entry.Fields = make(Fields)
	inTags := false
	inMessage := false
	messageSet := false
	inFields := false
	inEscape := false
	currentKey := ""
	i := 0
	currentStr := bytes.Buffer{}
	for {
		r, _, err := line.ReadRune()
		i++
		if err == io.EOF {
			// If the message hasn't been set this is the
			// message
			if !messageSet {
				entry.Message = strings.TrimSpace(currentStr.String())
			}
			if inFields {
				entry.Fields[strings.TrimSpace(currentKey)] = BasicField(bytes.TrimSpace(currentStr.Bytes()))
			}
			// io.EOF is ok
			return entry, nil
		}

		if err != nil {
			return Entry{}, err
		}

		// If the first character is open bracket
		// ignore
		if i == 1 && r == '[' {
			// Set in tags flag
			inTags = true
			// Skip to end
			continue
		}
		if i == 1 && r != '[' {
			// No tags
			inMessage = true
		}

		if !inEscape && r == '\\' {
			inEscape = true
			continue
		}
		if inEscape && (r == '\\' || r == '[' || r == '|') {
			currentStr.WriteRune(r)
			inEscape = false
			continue
		}

		if inTags && !inEscape && r == ']' {
			// End of tags sequence
			inTags = false
			entry.Tags = append(entry.Tags, Tag(strings.TrimSpace(currentStr.String())))
			inMessage = true
			currentStr.Reset()
			continue

		}
		if inTags && r == '|' {
			entry.Tags = append(entry.Tags, Tag(currentStr.String()))
			currentStr.Reset()
			continue
		}

		if inMessage && !inEscape && r == '|' {
			entry.Message = strings.TrimSpace(currentStr.String())
			currentStr.Reset()
			inMessage = false
			messageSet = true
			inFields = true
			continue
		}

		if inFields && r == '|' {
			entry.Fields[strings.TrimSpace(currentKey)] = BasicField(strings.TrimSpace(currentStr.String()))
			currentStr.Reset()
			continue
		}

		if inFields && inEscape && r == ':' {
			inEscape = false
			currentStr.WriteRune(r)
			continue
		}

		if inFields && !inEscape && r == ':' {
			// In the fields, look for unescaped colon
			// currentStr is now the key
			currentKey = currentStr.String()
			currentStr.Reset()
			continue
		}

		// Write the character to the current string
		currentStr.WriteRune(r)

	}

}
