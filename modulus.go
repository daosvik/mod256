// mod256: Arithmetic modulo 193-256 bit moduli 
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

type Modulus struct {
	m	[4]uint64	// modulus
	mu	[5]uint64	// reciprocal
	mmu0	[4]uint64	// (m*mu + 0*m)
	mmu1	[4]uint64	// (m*mu + 1*m) % 2^256
}

// Set value from little-endian array of uint64
func (z *Modulus) FromUint64(m [4]uint64) (*Modulus) {

	if m[3] == 0 {
		panic("Modulus < 2^192")
	}

	// Store the modulus itself

	z.m = m

	// Compute reciprocal of m

	z.mu = reciprocal(m)

	// Compute mmu0, mmu1

	mmu0(z)
	mmu1(z)

	return z
}

func (z *Modulus) ToUint64() ([4]uint64) {
	return z.m
}
