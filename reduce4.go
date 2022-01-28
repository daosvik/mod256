// mod256: Arithmetic modulo 193-256 bit moduli 
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
//	"fmt"
)

// reduce4 computes the least non-negative residue of z
// and stores it back in z
func (z *Residue) reduce4() (*Residue) {

	// NB: Most variable names in the comments match the pseudocode for
	// 	Barrett reduction in the Handbook of Applied Cryptography.

	var x0, x1, x2, x3, x4, r0, r1, r2, r3, r4, q3, t0, t1, c uint64

	mu := z.m.mu
	m  := z.m.m

	// q1 = x/2^192
	// q2 = q1 * mu; q3 = q2 / 2^320

	q3, _ = Mul64(z.r[3], mu[4])

	// r1 = x mod 2^320 = x

	x0 = z.r[0]
	x1 = z.r[1]
	x2 = z.r[2]
	x3 = z.r[3]
	x4 = 0

	// r2 = q3 * m mod 2^320

	r2, r1 = Mul64(q3, m[1])
	r4, r3 = Mul64(q3, m[3])

	t1, r0 = Mul64(q3, m[0]); r1, c = Add64(r1, t1, 0)
	t1, t0 = Mul64(q3, m[2]); r2, c = Add64(r2, t0, c); r3, c = Add64(r3, t1, c); r4, _ = Add64(r4, 0, c)

	// r = r1 - r2 = x - r2

	// Note: x < 2^256
	//    => q3 <= x/m
	//    => q3*m <= x
	//    => r2 <= x
	//    => r >= 0

	var b uint64

	r0, b = Sub64(x0, r0, 0)
	r1, b = Sub64(x1, r1, b)
	r2, b = Sub64(x2, r2, b)
	r3, b = Sub64(x3, r3, b)
	r4, _ = Sub64(x4, r4, b)

	for {
		// if r>=m then r-=m

		x0, b = Sub64(r0, m[0], 0)
		x1, b = Sub64(r1, m[1], b)
		x2, b = Sub64(r2, m[2], b)
		x3, b = Sub64(r3, m[3], b)
		x4, b = Sub64(r4,    0, b)

		if b != 0 {
			break
		}

		// commit if no borrow (r1 >= r2 + m)

		r4, r3, r2, r1, r0 = x4, x3, x2, x1, x0
	}

	z.r[3], z.r[2], z.r[1], z.r[0] = r3, r2, r1, r0

	return z
}
