// Copyright (c) 2017, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package reclist reads
// and writes records in a simple list format.
//
// The relist format is a human readable,
// utf-8 encoded text-file,
// with data records.
//
// In reclist format each record has a type and ID,
// the record type is indicated by a string starting
// with the at sign (@) and separated from the ID
// with an equal sign (=).
// Records are composed by key value pairs,
// each pair on a single line,
// with the key separated from its value
// by the colon sign (:).
// Type,
// and keys
// are case insensitive,
// and without spaces.
// Blank lines
// and lines starting with the sharp sign (#)
// will be ignored.
// Multiple line values
// should be enclosed by quotation marks (“),
// if in the multi-line value a quotation is used,
// it can be escaped
// using the slash (\) before the quotation mark.
// Leading spaces will be ignored.
//
// An example of a reclist is:
//
//	# Solar system objects
//	@star=Sun
//		radius:	109.3
//		mass:	333000
//		gravity: 27.94
//		descrip: "The Sun is the star at the center
//			of the Solar System. It is a nearly
//			perfect sphere of hot plasma. It is
//			by far the most important source of
//			energy for life on Earth."
//
//	@planet=Jupiter
//		radius:	10.97
//		mass:	317.83
//		gravity: 2.528
//		descrip: "Jupiter is the fifth planet from
//			the Sun and the largest in the Solar
//			System. It is a giant planet with a
//			mass one-thousandth of the Sun, but
//			two-and-a-half times that of all other
//			planets in the Solar System combined."
//		moons:	Ganymede Callisto Io Europa
//
//	@planet=Mars
//		radius: 0.5320
//		mass:	0.107
//		gravity: 0.38
//		descrip: "Mars is the fourth planet from the Sun
//			and the second-smallest planet in the
//			Solar System after Mercury. Mars is often
//			referred as the \"Red Planet\" because
//			the iron oxide prevalent on its	surface
//			gives it a reddish appearance that is
//			disctintive among the astronomical bodies
//			visible to the naked eye."
//
//	@moon=Titan
//		radius:	0.4043
//		mass:	0.0225
//		gravity: 0.14
//		parent: Saturn
//
//	@dwarf=Eris
//		radius:	0.1825
//		mass:	0.0028
//		gravity: 0.0672
//		family:	SDO
//
// Note that indentation
// and blank lines are optional
// and used for reading ease.
//
// This format is inspired from
// the record-jar format
// <http://www.catb.org/esr/writings/taoup/html/ch05s02.html#id2906931>
// described by
// E. Raymond in "The Art of Unix Programming" (2003) Addison-Wesley,
// by the list format
// <http://www.strozzi.it/cgi-bin/CSA/tw7/I/en_US/NoSQL/Table%20structure>
// of several flat text database systems
// such as C. Strozzi "NoSQL" (2007),
// and the BibTeX bibliography format
// <https://en.wikipedia.org/wiki/BibTeX>
// of O. Patashnik and L. Lamport (1985).
package reclist

import (
	"sort"
	"strings"
)

// Record holds the information of a reclist record.
type Record struct {
	id   string
	typ  string
	data map[string]string
}

// NewRecord returns a new record with
// the given type and ID.
func NewRecord(typ, ID string) *Record {
	typ = strings.ToLower(strings.Join(strings.Fields(typ), "-"))
	if typ == "" {
		return nil
	}
	ID = strings.Join(strings.Fields(ID), " ")
	if ID == "" {
		return nil
	}
	return &Record{
		id:   ID,
		typ:  typ,
		data: make(map[string]string),
	}
}

// ID returns the value of the ID.
func (r *Record) ID() string {
	return r.id
}

// Type returs the value of the record type.
func (r *Record) Type() string {
	return r.typ
}

// Keys returns the keys of the given record.
func (r *Record) Keys() []string {
	if len(r.data) == 0 {
		return nil
	}
	ks := make([]string, 0, len(r.data))
	for key := range r.data {
		ks = append(ks, key)
	}
	sort.Strings(ks)
	return ks
}

// Get returns the value associated with the given key.
func (r *Record) Get(key string) string {
	key = strings.ToLower(strings.Join(strings.Fields(key), "-"))
	return r.data[key]
}

// Set sets a value for a given key,
// if the key does not exists,
// it will be added to the record.
func (r *Record) Set(key, value string) {
	key = strings.ToLower(strings.Join(strings.Fields(key), "-"))
	if key == "" {
		return
	}
	value = strings.TrimSpace(value)
	if value == "" {
		delete(r.data, key)
		return
	}
	r.data[key] = value
}
