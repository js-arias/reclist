// Copyright (c) 2017, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package reclist

import (
	"bytes"
	"strings"
	"testing"
)

func TestRecListWrite(t *testing.T) {
	s := NewScanner(strings.NewReader(blob))
	out := &bytes.Buffer{}
	w := NewWriter(out)

	var recs []*Record

	for s.Scan() {
		rec := s.Record()
		if err := w.Write(rec); err != err {
			t.Errorf("unexpected error: %v", err)
		}
		recs = append(recs, rec)
	}
	if err := s.Err(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	w.Flush()
	if err := w.Err(); err != nil {
		t.Errorf("unexpected error: &v", err)
	}

	s = NewScanner(strings.NewReader(out.String()))
	i := 0
	for s.Scan() {
		rec := s.Record()
		for _, key := range recs[i].Keys() {
			if rec.Get(key) != recs[i].Get(key) {
				t.Errorf("%s %s, key %s = %q, want %q", rec.Type(), rec.ID(), key, rec.Get(key), recs[i].Get(key))
			}
		}
		for _, key := range rec.Keys() {
			if rec.Get(key) != recs[i].Get(key) {
				t.Errorf("%s %s, key %s = %q, want %q", rec.Type(), rec.ID(), key, rec.Get(key), recs[i].Get(key))
			}
		}
		i++
	}
	if err := s.Err(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i != len(recs) {
		t.Errorf("%d records, want %d", i, len(recs))
	}
}
