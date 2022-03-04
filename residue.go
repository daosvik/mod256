// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

// Residue contains a representative of a residue class, and the pointer to its modulus.
// The residue is stored as any 256-bit unsigned integer in the residue
// class, represented by a 4-element little-endian array of uint64.
type Residue struct {
	m *Modulus
	r [4]uint64
}

// FromUint64 sets the residue value from a little-endian array of uint64.
func (z *Residue) FromUint64(m *Modulus, x [4]uint64) *Residue {

	if m.m[3] == 0 {
		panic("Modulus < 2^192")
	}

	z.m = m
	z.r = x
	return z
}

// ToUint64 returns an array with the canonical representative of the residue class.
func (z *Residue) ToUint64() [4]uint64 {
	z.reduce4() // Reduce to canonical residue
	return z.r
}

// Copy copies one residue to another.
// Both the residue value and the modulus pointer are copied.
func (z *Residue) Copy(x *Residue) *Residue {
	z.m = x.m
	z.r = x.r
	return z
}

// shiftleft256 shifts the 256-bit value in a little-endian array left by 0-63 bits.
func shiftleft256(x [4]uint64, s uint) (z [4]uint64) {
	l := s % 64	// left shift
	r := 64 - l	// right shift

	z[0] = (x[0] << l)
	z[1] = (x[1] << l) | (x[0] >> r)
	z[2] = (x[2] << l) | (x[1] >> r)
	z[3] = (x[3] << l) | (x[2] >> r)

	return z
}
