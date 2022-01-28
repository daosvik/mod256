// mod256: Arithmetic modulo 193-256 bit moduli 
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

// A residue is stored as any 256-bit unsigned integer in the residue
// class, represented by a 4-element little-endian array of uint64

type Residue struct {
	m *Modulus
	r [4]uint64
}

// Set value from little-endian array of uint64
func (z *Residue) FromUint64(m *Modulus, x [4]uint64) (*Residue) {

	if m.m[3] == 0 {
		panic("Modulus < 2^192")
	}

	z.m = m
	z.r = x
	return z
}


func (z *Residue) ToUint64() ([4]uint64) {
	z.reduce4() // Reduce to canonical residue
	return z.r
}

func (z *Residue) Cpy(x *Residue) (*Residue) {
	z.m = x.m
	z.r = x.r
	return z
}

func shiftleft256(x [4]uint64, s uint) (z [4]uint64) {
	l := s % 64	// left shift
	r := 64 - l	// right shift

	z[0] = (x[0] << l)
	z[1] = (x[1] << l) | (x[0] >> r)
	z[2] = (x[2] << l) | (x[1] >> r)
	z[3] = (x[3] << l) | (x[2] >> r)

	return z
}
