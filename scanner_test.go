// Copyright (c) 2017, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package reclist

import (
	"strings"
	"testing"
)

var blob = `
# Solar system objects
@star=Sun
	radius:	109.3
	mass:	333000
	gravity: 27.94
	descrip: "The Sun is the star at the center
		of the Solar System. It is a nearly
		perfect sphere of hot plasma. It is
		by far the most important source of
		energy for life on Earth."

@planet=Jupiter
	radius:	10.97
	mass:	317.83
	gravity: 2.528
	descrip: "Jupiter is the fifth planet from
		the Sun and the largest in the Solar
		System. It is a giant planet with a
		mass one-thousandth of the Sun, but
		two-and-a-half times that of all other
		planets in the Solar System combined."
	moons:	Ganymede Callisto Io Europa

@planet=Saturn
	radius:	9.140
	mass:	95.162
	gravity: 1.065
	descrip: "Saturn is the sixth planet from the
		Sun and the second-largest in the
		Solar System, after Jupiter. It is a
		gas giant with an average radius about
		nine times that of Earth. The planet's
		most famous feature is its prominent
		ring system that is composed mostly of
		ice particles, with a small amount of
		rocky debris and dust."
	moons:	Titan Rhea Iapetus Dione Tethys Enceladus Mimas Hyperion Phoebe

@planet=Neptune
	radius: 3.865
	mass:	17.147
	gravity: 1.137
	descrip: "Nepturne is the eight and farthest
		known planet from the Sun in the Solar
		System. In the Solar System, it is the
		fourth-largest planet by diameterm the
		third-most-massive planet, and the
		densest giant planet. Neptune is not
		visible to the unaided eye and is the
		only planet in the Solar System found by
		mathematical prediction rather than by
		empirical observation."
	moons:	Triton Proteus Nereid

@planet=Uranus
	radius: 3.981
	mass:	14.539
	gravity: 0.90
	descrip: "Uranus is the seventh planet from the
		Sun. It has the third-largest planetary
		radius and fourth-largest mass in the
		Solar System. Uranus's atmosphere is the
		coldest planetary atmosepher in the Solar
		System, with a minimum temperature of 49
		K. The Uranus system has a unique
		configuration among giant planets because
		its axis of rotation is tilted sideways,
		nearly into the plane of its solar
		orbit."
	moons:	Titania Oberon Umbriel Ariel Miranda

@planet=Earth
	radius: 1
	mass: 1
	gravity: 1
	descrip: "Earth is the third planet from the Sun
		and the only object in the Universe known
		to harbor life. Earth is the densest
		planet in the Solar Systen and the largest
		of the four terrestrial planets."
	moons:	Moon

@planet=Venus
	radius: 0.9499
	mass:	0.815
	gravity: 0.905
	descrip: "Venus is the second planet from the Sun.
		It has the longest rotation perion (243
		days) of any planet in the Solar System
		and rotates in the opposite direction of
		most other planets."

@planet=Mars
	radius: 0.5320
	mass:	0.107
	gravity: 0.38
	descrip: "Mars is the fourth planet from the Sun
		and the second-smallest planet in the
		Solar System after Mercury. Mars is often
		referred as the \"Red Planet\" because
		the iron oxide prevalent on its	surface
		gives it a reddish appearance that is
		disctintive among the astronomical bodies
		visible to the naked eye."

@planet=Mercury
	radius:	0.3829
	mass:	0.0553
	gravity: 0.38
	descrip: "Mercury is the smallest and innermost
		planet in the Solar System. Its orbital
		period around the Sun of 88 days is the
		shortest of all the planetes in the Solar
		System."

@moon=Ganymede
	radius: 0.4135
	mass:	0.0248
	gravity: 0.15
	parent:	Jupiter

@moon=Titan
	radius:	0.4043
	mass:	0.0225
	gravity: 0.14
	parent: Saturn

@moon=Callisto
	radius: 0.3783
	mass:	0.018
	gravity: 0.126
	parent: Jupiter

@moon=Io
	radius: 0.2859
	mass:	0.015
	gravity: 0.183
	parent: Jupiter

@moon=Moon
	radius: 0.2727
	mass:	0.0123
	gravity: 0.166
	parent: Earth

@moon=Europa
	radius: 0.1450
	mass:	0.0080
	gravity: 0.134
	parent: Jupiter

@moon=Triton
	radius: 0.2124
	mass:	0.0036
	gravity: 0.0797
	parent:	Neptune

@dwarf=Eris
	radius:	0.1825
	mass:	0.0028
	gravity: 0.0672
	family:	SDO
	moons:	Dysnomia

@dwarf=Pluto
	radius:	0.186
	mass:	0.0022
	gravity: 0.062
	family: Plutino
	moons:	Charon
`

var testData = []struct {
	name   string
	tp     string
	parent string
}{
	{"Sun", "star", ""},
	{"Jupiter", "planet", ""},
	{"Saturn", "planet", ""},
	{"Neptune", "planet", ""},
	{"Uranus", "planet", ""},
	{"Earth", "planet", ""},
	{"Venus", "planet", ""},
	{"Mars", "planet", ""},
	{"Mercury", "planet", ""},
	{"Ganymede", "moon", "Jupiter"},
	{"Titan", "moon", "Saturn"},
	{"Callisto", "moon", "Jupiter"},
	{"Io", "moon", "Jupiter"},
	{"Moon", "moon", "Earth"},
	{"Europa", "moon", "Jupiter"},
	{"Triton", "moon", "Neptune"},
	{"Eris", "dwarf", ""},
	{"Pluto", "dwarf", ""},
}

func TestReclistScan(t *testing.T) {
	s := NewScanner(strings.NewReader(blob))
	i := 0
	for s.Scan() {
		rec := s.Record()
		if rec.ID() != testData[i].name {
			t.Errorf("%s, want %s", rec.ID(), testData[i].name)
		}
		if rec.Type() != testData[i].tp {
			t.Errorf("%s type %q, want %q", rec.ID(), rec.Type(), testData[i].tp)
		}
		if rec.Type() == "planet" {
			ds := rec.Get("descrip")
			if len(ds) == 0 {
				t.Errorf("%s empty description", rec.ID())
			}
		}
		if rec.Get("parent") != testData[i].parent {
			t.Errorf("%s parent %q, want %q", rec.ID(), rec.Get("parent"), testData[i].parent)
		}
		i++
	}
	if err := s.Err(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i != len(testData) {
		t.Errorf("%d records, want %d", i, len(testData))
	}
}
