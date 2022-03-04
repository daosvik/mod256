// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

// Mul computes the product of two residues.
func (z *Residue) Mul(x *Residue) *Residue {
	if z == x {
		return z.Square()
	}

	if z.m != x.m {
		if z.m.m != x.m.m {
			panic("Incompatible moduli")
		}
	}

	// Multiplication

	var c, t0, t1, q0, q1, q2, q3, q4, q5, q6, q7 uint64

	q2, q1 = Mul64(z.r[0], x.r[1])
	q4, q3 = Mul64(z.r[0], x.r[3])

	t1, q0 = Mul64(z.r[0], x.r[0]); q1, c = Add64(q1, t1, 0)
	t1, t0 = Mul64(z.r[0], x.r[2]); q2, c = Add64(q2, t0, c); q3, c = Add64(q3, t1, c); q4, _ = Add64(q4, 0, c)

	t1, t0 = Mul64(z.r[1], x.r[1]); q2, c = Add64(q2, t0, 0); q3, c = Add64(q3, t1, c)
	q5, t0 = Mul64(z.r[1], x.r[3]); q4, c = Add64(q4, t0, c); q5, _ = Add64(q5,  0, c)

	t1, t0 = Mul64(z.r[1], x.r[0]); q1, c = Add64(q1, t0, 0); q2, c = Add64(q2, t1, c)
	t1, t0 = Mul64(z.r[1], x.r[2]); q3, c = Add64(q3, t0, c); q4, c = Add64(q4, t1, c); q5, _ = Add64(q5, 0, c)

	t1, t0 = Mul64(z.r[2], x.r[1]); q3, c = Add64(q3, t0, 0); q4, c = Add64(q4, t1, c)
	q6, t0 = Mul64(z.r[2], x.r[3]); q5, c = Add64(q5, t0, c); q6, _ = Add64(q6,  0, c)

	t1, t0 = Mul64(z.r[2], x.r[0]); q2, c = Add64(q2, t0, 0); q3, c = Add64(q3, t1, c)
	t1, t0 = Mul64(z.r[2], x.r[2]); q4, c = Add64(q4, t0, c); q5, c = Add64(q5, t1, c); q6, _ = Add64(q6, 0, c)

	t1, t0 = Mul64(z.r[3], x.r[1]); q4, c = Add64(q4, t0, 0); q5, c = Add64(q5, t1, c)
	q7, t0 = Mul64(z.r[3], x.r[3]); q6, c = Add64(q6, t0, c); q7, _ = Add64(q7,  0, c)

	t1, t0 = Mul64(z.r[3], x.r[0]); q3, c = Add64(q3, t0, 0); q4, c = Add64(q4, t1, c)
	t1, t0 = Mul64(z.r[3], x.r[2]); q5, c = Add64(q5, t0, c); q6, c = Add64(q6, t1, c); q7, _ = Add64(q7, 0, c)

	// Reduction

	z.reduce8([8]uint64{ q0, q1, q2, q3, q4, q5, q6, q7 })

	return z
}
