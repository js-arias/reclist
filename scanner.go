// Copyright (c) 2017, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package reclist

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// A Scanner reads records
// from a reclist encoded file.
type Scanner struct {
	closed bool
	err    error
	r      *bufio.Reader
	line   int
	rec    *Record
	b      *bytes.Buffer
}

// NewScanner returns a new scanner
// that reads from r.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
		b: &bytes.Buffer{},
	}
}

// Scan prepares the next record for reading
// with the Record method.
// It returns true on success,
// or false if there is no next record
// or an error happened while preparing it.
// Err should be consulted to distinguish
// between the two cases.
//
// Every call to Record,
// must be preceded by a call to Scan.
func (s *Scanner) Scan() bool {
	if s.closed {
		return false
	}
	for {
		rec, err := s.parseRecord()
		if err != nil {
			s.closed = true
			if errors.Cause(err) == io.EOF {
				return false
			}
			s.err = err
			return false
		}
		if len(rec.data) == 0 {
			continue
		}
		s.rec = rec
		return true
	}
}

// ParseRecord reads and parses a single
// reclist record from s.
func (s *Scanner) parseRecord() (*Record, error) {
	var rec *Record

	// read the "header" of the record
	for {
		s.line++
		line, err := s.r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, err
			}
			return nil, errors.Wrapf(err, "reclist: scanner: line %d", s.line)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == '#' {
			continue
		}
		if line[0] != '@' {
			continue
		}
		line = line[1:]

		i := strings.Index(line, "=")
		if i < 1 {
			continue
		}
		rec = NewRecord(line[:i], line[i+1:])
		if rec != nil {
			break
		}
	}

	// read the record data.
	for {
		key, delim, err := s.parseKey()
		if err != nil {
			return rec, nil
		}
		if delim == '@' {
			return rec, nil
		}
		if delim == '\n' {
			continue
		}

		value, err := s.parseValue()
		if err != nil {
			return rec, nil
		}
		rec.Set(key, value)
	}
}

// ParseKey parses the next key in the record.
func (s *Scanner) parseKey() (key string, delim rune, err error) {
	for {
		r1, err := s.readRune()
		if err != nil {
			return "", 0, err
		}
		if r1 == '\n' {
			s.line++
			continue
		}
		if unicode.IsSpace(r1) {
			continue
		}
		if r1 == '#' {
			s.skip('\n')
			s.line++
			continue
		}
		s.r.UnreadRune()
		if r1 == '@' {
			return "", '@', nil
		}
		break
	}

	s.b.Reset()
	space := false
	for {
		r1, err := s.readRune()
		if err != nil {
			return "", 0, err
		}
		if r1 == '\n' {
			s.line++
			return "", '\n', nil
		}
		if unicode.IsSpace(r1) {
			space = true
			continue
		}
		if r1 == ':' {
			return s.b.String(), ':', nil
		}
		if space {
			s.b.WriteRune('-')
			space = false
		}
		s.b.WriteRune(r1)
	}
}

// ParseValue parses the next value in the record.
func (s *Scanner) parseValue() (value string, err error) {
	for {
		r1, err := s.readRune()
		if err != nil {
			return "", err
		}
		if r1 == '\n' {
			s.line++
			return "", nil
		}
		if unicode.IsSpace(r1) {
			continue
		}
		if r1 == '"' {
			return s.quoteValue()
		}
		s.r.UnreadRune()
		break
	}
	return s.r.ReadString('\n')
}

// QuoteValue parses a quoted,
// possible multiline value.
func (s *Scanner) quoteValue() (value string, err error) {
	s.b.Reset()
	first, space, line := true, false, false
	for {
		r1, err := s.readRune()
		if err != nil {
			v := s.b.String()
			if len(v) > 0 {
				return v, nil
			}
			return "", err
		}
		if r1 == '\n' {
			if !first {
				line = true
				space = false
			}
			s.line++
			continue
		}
		if unicode.IsSpace(r1) {
			if !first && !line {
				space = true
			}
			continue
		}
		if r1 == '"' {
			return s.b.String(), nil
		}
		if r1 == '\\' {
			r1, err = s.readRune()
			if err != nil {
				continue
			}
		}
		if line {
			line = false
			space = false
			s.b.WriteRune('\n')
		}
		if space {
			space = false
			s.b.WriteRune(' ')
		}
		first = false
		s.b.WriteRune(r1)
	}
}

// Record returns the last read record.
func (s *Scanner) Record() *Record {
	if s.rec == nil {
		panic("Record called without Scan")
	}
	rec := s.rec
	s.rec = nil
	return rec
}

// Err returns the error,
// if anym
// that was encountered during iteration.
func (s *Scanner) Err() error {
	return s.err
}

// readRune reads one run from s,
// folding \r\n to \n.
func (s *Scanner) readRune() (rune, error) {
	r1, _, err := s.r.ReadRune()

	// Handle \r\n.
	if r1 == '\r' {
		r1, _, err = s.r.ReadRune()
		if err == nil {
			if r1 != '\n' {
				s.r.UnreadRune()
				r1 = '\r'
			}
		}
	}
	return r1, err
}

// skip reads runes up to and including
// the rune delim
// or until error.
func (s *Scanner) skip(delim rune) error {
	for {
		r1, err := s.readRune()
		if err != nil {
			return err
		}
		if r1 == delim {
			return nil
		}
	}
}
