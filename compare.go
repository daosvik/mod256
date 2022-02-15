// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

// Eq compares one residue to another, returns true when equal.
func (x *Residue) Eq(y *Residue) bool {
	m := x.m.m[0] ^ y.m.m[0]
	m |= x.m.m[1] ^ y.m.m[1]
	m |= x.m.m[2] ^ y.m.m[2]
	m |= x.m.m[3] ^ y.m.m[3]

	if m != 0 {
		return false
	}

	x.reduce4()
	y.reduce4()

	r := x.r[0] ^ y.r[0]
	r |= x.r[1] ^ y.r[1]
	r |= x.r[2] ^ y.r[2]
	r |= x.r[3] ^ y.r[3]

	return r == 0
}

// Neq compares one residue to another, returns true when different.
func (x *Residue) Neq(y *Residue) bool {
	m := x.m.m[0] ^ y.m.m[0]
	m |= x.m.m[1] ^ y.m.m[1]
	m |= x.m.m[2] ^ y.m.m[2]
	m |= x.m.m[3] ^ y.m.m[3]

	if m != 0 {
		return true
	}

	x.reduce4()
	y.reduce4()

	r := x.r[0] ^ y.r[0]
	r |= x.r[1] ^ y.r[1]
	r |= x.r[2] ^ y.r[2]
	r |= x.r[3] ^ y.r[3]

	return r != 0
}
