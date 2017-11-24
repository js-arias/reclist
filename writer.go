// Copyright (c) 2017, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package reclist

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// A Writer writes records to a reclist enconded file.
type Writer struct {
	w *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

// Write writes a single record to w.
func (w *Writer) Write(rec *Record) error {
	keys := rec.Keys()
	if len(keys) == 0 {
		return nil
	}

	if rec.ID() == "" || rec.Type() == "" {
		return nil
	}

	if _, err := w.w.WriteString("@" + rec.Type() + "=" + rec.ID() + "\n"); err != nil {
		return errors.Wrap(err, "reclist: writer")
	}

	for _, key := range keys {
		value := rec.Get(key)
		if value == "" {
			continue
		}
		if err := w.writeField(key, value); err != nil {
			return err
		}
	}
	return nil
}

// WriteField writes a key, value pair.
func (w *Writer) writeField(key, value string) error {
	var s string

	if strings.Contains(value, "\n") {
		return w.writeQuoted(key, value)
	}

	if len(key) < 6 {
		s = "\t" + key + ":\t" + value + "\n"
	} else {
		s = "\t" + key + ": " + value + "\n"
	}
	if _, err := w.w.WriteString(s); err != nil {
		return errors.Wrap(err, "reclist: writer")
	}
	return nil
}

// WriteQuoted writes a key, value pair,
// in which value is enclosed by
// quotation marks.
func (w *Writer) writeQuoted(key, value string) error {
	var s string
	if len(key) < 6 {
		s = "\t" + key + ":\t\""
	} else {
		s = "\t" + key + ": \""
	}
	if _, err := w.w.WriteString(s); err != nil {
		return errors.Wrap(err, "reclist: writer")
	}

	first, space, line := true, false, false
	for _, r := range value {
		if r == '\n' {
			if !first {
				line = true
				space = false
			}
			continue
		}
		if unicode.IsSpace(r) {
			if !first && !line {
				space = true
			}
			continue
		}
		if line {
			line = false
			space = false
			if _, err := w.w.WriteString("\n\t\t"); err != nil {
				return errors.Wrap(err, "reclist: writer")
			}
		}
		if space {
			space = false
			if _, err := w.w.WriteRune(' '); err != nil {
				return errors.Wrap(err, "reclist: writer")
			}
		}
		first = false
		if r == '"' {
			if _, err := w.w.WriteString(`\"`); err != nil {
				return errors.Wrap(err, "reclist: writer")
			}
			continue
		}
		if r == '\\' {
			if _, err := w.w.WriteString(`\\`); err != nil {
				return errors.Wrap(err, "reclist: writer")
			}
			continue
		}
		if _, err := w.w.WriteRune(r); err != nil {
			return errors.Wrap(err, "reclist: writer")
		}
	}
	if _, err := w.w.WriteString("\"\n"); err != nil {
		return errors.Wrap(err, "reclist: writer")
	}
	return nil
}

// Flush writes any buffered data
// to the underlying io.Writer.
// To check if an error occurred during the flush,
// call Err.
func (w *Writer) Flush() {
	w.w.Flush()
}

// Err reports any errar that has occurred
// during a previous Write
// of Flush.
func (w *Writer) Err() error {
	if _, err := w.w.Write(nil); err != nil {
		return errors.Wrap(err, "reclist: writer")
	}
	return nil
}
