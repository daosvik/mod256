// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

// reduce8 computes a 256-bit residue of x modulo z.m and stores it in z
func (z *Residue) reduce8(x [8]uint64) *Residue {

	// NB: Most variable names in the comments match the pseudocode for
	// 	Barrett reduction in the Handbook of Applied Cryptography.

	mu := z.m.mu
	m := z.m.m

	// q1 = x/2^192

	x0 := x[3]
	x1 := x[4]
	x2 := x[5]
	x3 := x[6]
	x4 := x[7]

	// q2 = q1 * mu; q3 = q2 / 2^320

	var q0, q1, q2, q3, q4, q5, t0, t1, c uint64

	q0, _  = Mul64(x3, mu[0])
	q1, t0 = Mul64(x4, mu[0]); q0, c = Add64(q0, t0, 0); q1, _ = Add64(q1,  0, c)


	t1, _  = Mul64(x2, mu[1]); q0, c = Add64(q0, t1, 0)
	q2, t0 = Mul64(x4, mu[1]); q1, c = Add64(q1, t0, c); q2, _ = Add64(q2,  0, c)

	t1, t0 = Mul64(x3, mu[1]); q0, c = Add64(q0, t0, 0); q1, c = Add64(q1, t1, c); q2, _ = Add64(q2, 0, c)


	t1, t0 = Mul64(x2, mu[2]); q0, c = Add64(q0, t0, 0); q1, c = Add64(q1, t1, c)
	q3, t0 = Mul64(x4, mu[2]); q2, c = Add64(q2, t0, c); q3, _ = Add64(q3,  0, c)

	t1, _  = Mul64(x1, mu[2]); q0, c = Add64(q0, t1, 0)
	t1, t0 = Mul64(x3, mu[2]); q1, c = Add64(q1, t0, c); q2, c = Add64(q2, t1, c); q3, _ = Add64(q3, 0, c)


	t1, _  = Mul64(x0, mu[3]); q0, c = Add64(q0, t1, 0)
	t1, t0 = Mul64(x2, mu[3]); q1, c = Add64(q1, t0, c); q2, c = Add64(q2, t1, c)
	q4, t0 = Mul64(x4, mu[3]); q3, c = Add64(q3, t0, c); q4, _ = Add64(q4,  0, c)

	t1, t0 = Mul64(x1, mu[3]); q0, c = Add64(q0, t0, 0); q1, c = Add64(q1, t1, c)
	t1, t0 = Mul64(x3, mu[3]); q2, c = Add64(q2, t0, c); q3, c = Add64(q3, t1, c); q4, _ = Add64(q4, 0, c)


	t1, t0 = Mul64(x0, mu[4]); _,  c = Add64(q0, t0, 0); q1, c = Add64(q1, t1, c)
	t1, t0 = Mul64(x2, mu[4]); q2, c = Add64(q2, t0, c); q3, c = Add64(q3, t1, c)
	q5, t0 = Mul64(x4, mu[4]); q4, c = Add64(q4, t0, c); q5, _ = Add64(q5,  0, c)

	t1, t0 = Mul64(x1, mu[4]); q1, c = Add64(q1, t0, 0); q2, c = Add64(q2, t1, c)
	t1, t0 = Mul64(x3, mu[4]); q3, c = Add64(q3, t0, c); q4, c = Add64(q4, t1, c); q5, _ = Add64(q5, 0, c)

	// Drop the fractional part of q3

	q0 = q1
	q1 = q2
	q2 = q3
	q3 = q4
	q4 = q5

	// r1 = x mod 2^320

	x0 = x[0]
	x1 = x[1]
	x2 = x[2]
	x3 = x[3]
	x4 = x[4]

	// r2 = q3 * m mod 2^320

	var r0, r1, r2, r3, r4 uint64

	r4, r3 = Mul64(q0, m[3])
	_,  t0 = Mul64(q1, m[3]); r4, _ = Add64(r4, t0, 0)


	t1, r2 = Mul64(q0, m[2]); r3, c = Add64(r3, t1, 0)
	_,  t0 = Mul64(q2, m[2]); r4, _ = Add64(r4, t0, c)

	t1, t0 = Mul64(q1, m[2]); r3, c = Add64(r3, t0, 0); r4, _ = Add64(r4, t1, c)


	t1, r1 = Mul64(q0, m[1]); r2, c = Add64(r2, t1, 0)
	t1, t0 = Mul64(q2, m[1]); r3, c = Add64(r3, t0, c); r4, _ = Add64(r4, t1, c)

	t1, t0 = Mul64(q1, m[1]); r2, c = Add64(r2, t0, 0); r3, c = Add64(r3, t1, c)
	_,  t0 = Mul64(q3, m[1]); r4, _ = Add64(r4, t0, c)


	t1, r0 = Mul64(q0, m[0]); r1, c = Add64(r1, t1, 0)
	t1, t0 = Mul64(q2, m[0]); r2, c = Add64(r2, t0, c); r3, c = Add64(r3, t1, c)
	_,  t0 = Mul64(q4, m[0]); r4, _ = Add64(r4, t0, c)

	t1, t0 = Mul64(q1, m[0]); r1, c = Add64(r1, t0, 0); r2, c = Add64(r2, t1, c)
	t1, t0 = Mul64(q3, m[0]); r3, c = Add64(r3, t0, c); r4, _ = Add64(r4, t1, c)


	// r = r1 - r2

	var b uint64

	r0, b = Sub64(x0, r0, 0)
	r1, b = Sub64(x1, r1, b)
	r2, b = Sub64(x2, r2, b)
	r3, b = Sub64(x3, r3, b)
	r4, _ = Sub64(x4, r4, b)

	if r4 == 0 {
		z.r[3], z.r[2], z.r[1], z.r[0] = r3, r2, r1, r0
		return z
	}

	// q = r - m
	r0, b = Sub64(r0, m[0], 0)
	r1, b = Sub64(r1, m[1], b)
	r2, b = Sub64(r2, m[2], b)
	r3, b = Sub64(r3, m[3], b)
	r4, b = Sub64(r4,    0, b)

	// q = r - m
	x0, b = Sub64(r0, m[0], 0)
	x1, b = Sub64(r1, m[1], b)
	x2, b = Sub64(r2, m[2], b)
	x3, b = Sub64(r3, m[3], b)
	x4, b = Sub64(r4,    0, b)

	// commit if no borrow
	if b == 0 {
		r4, r3, r2, r1, r0 = x4, x3, x2, x1, x0
	}

	z.r[3], z.r[2], z.r[1], z.r[0] = r3, r2, r1, r0

	return z
}
