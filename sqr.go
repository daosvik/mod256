// mod256: Arithmetic modulo 193-256 bit moduli 
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

func (z *Residue) Sqr() (*Residue) {
	// Squaring

	var c, t0, t1, q0, q1, q2, q3, q4, q5, q6, q7 uint64

	q4, q3 = Mul64(z.r[0], z.r[3])

	t1, q2 = Mul64(z.r[0], z.r[2]); q3, c = Add64(q3, t1, 0)
	q5, t0 = Mul64(z.r[1], z.r[3]); q4, c = Add64(q4, t0, c); q5, c = Add64(q5, 0, c)

	t1, q1 = Mul64(z.r[0], z.r[1]); q2, c = Add64(q2, t1, 0)
	t1, t0 = Mul64(z.r[1], z.r[2]); q3, c = Add64(q3, t0, c); q4, c = Add64(q4, t1, c)
	q6, t0 = Mul64(z.r[2], z.r[3]); q5, c = Add64(q5, t0, c); q6, c = Add64(q6, 0, c)

	q1, c = Add64(q1, q1, 0)
	q2, c = Add64(q2, q2, c)
	q3, c = Add64(q3, q3, c)
	q4, c = Add64(q4, q4, c)
	q5, c = Add64(q5, q5, c)
	q6, c = Add64(q6, q6, c)
	q7, _ = Add64( 0,  0, c)

	t1, q0 = Mul64(z.r[0], z.r[0]); q1, c = Add64(q1, t1, 0)
	t1, t0 = Mul64(z.r[1], z.r[1]); q2, c = Add64(q2, t0, c); q3, c = Add64(q3, t1, c)
	t1, t0 = Mul64(z.r[2], z.r[2]); q4, c = Add64(q4, t0, c); q5, c = Add64(q5, t1, c)
	t1, t0 = Mul64(z.r[3], z.r[3]); q6, c = Add64(q6, t0, c); q7, _ = Add64(q7, t1, c)

	// Reduction

	z.reduce8([8]uint64{ q0, q1, q2, q3, q4, q5, q6, q7 })

	return z
}
