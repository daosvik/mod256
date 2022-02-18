// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	"errors"
)

// Modulus contains a modulus `m` as well as derived values that help speed up computations.
// The allowed range for `m` is `2^192` to `2^256-1`.
type Modulus struct {
	m    [4]uint64 // modulus
	mu   [5]uint64 // reciprocal
	mmu0 [4]uint64 // (m*mu + 0*m)
	mmu1 [4]uint64 // (m*mu + 1*m) % 2^256
}

// NewModulusFromUint64 creates a new modulus object from a little-endian array of uint64.
func NewModulusFromUint64(m [4]uint64) (z *Modulus, err error) {

	if m[3] == 0 {
		return nil, errors.New("Modulus < 2^192")
	}


	// Store the modulus itself
	z = &Modulus{m: m}

	// Compute reciprocal of m

	z.mu = reciprocal(m)

	// Compute mmu0, mmu1

	mmu0(z)
	mmu1(z)

	return
}

// ToUint64 returns an array with the modulus.
func (z *Modulus) ToUint64() [4]uint64 {
	return z.m
}
