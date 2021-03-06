// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

// reciprocal computes a 320-bit value mu representing 2^512/m
// (or equivalently 1/m in a 0.320-bit fixed point representation)
//
// Notes:
// - starts with a 32-bit division, refines with newton-raphson iterations
// - mu * m < 2^512
// - mu * m + m >= 2^512
// - mu * m + m == 2^512 iff m is a power of 2
func reciprocal(m [4]uint64) (mu [5]uint64) {

	// Note: specialized for m[3] != 0

	s := LeadingZeros64(m[3])
	p := 63 - s

	// 0 or a power of 2?

	// Check if at least one bit is set in m[2], m[1] or m[0],
	// or at least two bits in m[3]
	// If not, m is 0 or a power of 2

	if m[0] | m[1] | m[2] | (m[3] & (m[3]-1)) == 0 {
		mu[4] = ^uint64(0) >> uint(p & 63)
		mu[3] = ^uint64(0)
		mu[2] = ^uint64(0)
		mu[1] = ^uint64(0)
		mu[0] = ^uint64(0)

		return mu
	}

	// Maximise division precision by left-aligning divisor

	var (
		y  [4]uint64 // left-aligned copy of m
		r0 uint32    // estimate of 2^31/y
	)

	y = shiftleft256(m, uint(s))	// 1/2 < y < 1

	// Extract most significant 32 bits

	yh := uint32(y[3] >> 32)


	if yh == 0x80000000 { // Avoid overflow in division
		r0 = 0xffffffff
	} else {
		r0, _ = Div32(0x80000000, 0, yh)
	}

	// First iteration: 32 -> 64

	t1 := uint64(r0)		// 2^31/y
	t1 *= t1			// 2^62/y^2
	t1, _ = Mul64(t1, y[3])	// 2^62/y^2 * 2^64/y / 2^64 = 2^62/y

	r1 := uint64(r0) << 32		// 2^63/y
	r1 -= t1			// 2^63/y - 2^62/y = 2^62/y
	r1 *= 2				// 2^63/y

	if (r1 | (y[3]<<1)) == 0 {
		r1 = ^uint64(0)
	}

	// Second iteration: 64 -> 128

	// square: 2^126/y^2
	a2h, a2l := Mul64(r1, r1)

	// multiply by y: e2h:e2l:b2h = 2^126/y^2 * 2^128/y / 2^128 = 2^126/y
	b2h, _   := Mul64(a2l, y[2])
	c2h, c2l := Mul64(a2l, y[3])
	d2h, d2l := Mul64(a2h, y[2])
	e2h, e2l := Mul64(a2h, y[3])

	b2h, c   := Add64(b2h, c2l, 0)
	e2l, c    = Add64(e2l, c2h, c)
	e2h, _    = Add64(e2h,   0, c)

	_,   c    = Add64(b2h, d2l, 0)
	e2l, c    = Add64(e2l, d2h, c)
	e2h, _    = Add64(e2h,   0, c)

	// subtract: t2h:t2l = 2^127/y - 2^126/y = 2^126/y
	t2l, b   := Sub64( 0, e2l, 0)
	t2h, _   := Sub64(r1, e2h, b)

	// double: r2h:r2l = 2^127/y
	r2l, c   := Add64(t2l, t2l, 0)
	r2h, _   := Add64(t2h, t2h, c)

	if (r2h | r2l | (y[3]<<1)) == 0 {
		r2h = ^uint64(0)
		r2l = ^uint64(0)
	}

	// Third iteration: 128 -> 192

	// square r2 (keep 256 bits): 2^190/y^2
	a3h, a3l := Mul64(r2l, r2l)
	b3h, b3l := Mul64(r2l, r2h)
	c3h, c3l := Mul64(r2h, r2h)

	a3h, c    = Add64(a3h, b3l, 0)
	c3l, c    = Add64(c3l, b3h, c)
	c3h, _    = Add64(c3h,   0, c)

	a3h, c    = Add64(a3h, b3l, 0)
	c3l, c    = Add64(c3l, b3h, c)
	c3h, _    = Add64(c3h,   0, c)

	// multiply by y: q = 2^190/y^2 * 2^192/y / 2^192 = 2^190/y

	x0 := a3l
	x1 := a3h
	x2 := c3l
	x3 := c3h

	var q0, q1, q2, q3, q4, t0 uint64

	q0, _  = Mul64(x2, y[0])
	q1, t0 = Mul64(x3, y[0]); q0, c = Add64(q0, t0, 0); q1, _ = Add64(q1,  0, c)


	t1, _  = Mul64(x1, y[1]); q0, c = Add64(q0, t1, 0)
	q2, t0 = Mul64(x3, y[1]); q1, c = Add64(q1, t0, c); q2, _ = Add64(q2,  0, c)

	t1, t0 = Mul64(x2, y[1]); q0, c = Add64(q0, t0, 0); q1, c = Add64(q1, t1, c); q2, _ = Add64(q2, 0, c)


	t1, t0 = Mul64(x1, y[2]); q0, c = Add64(q0, t0, 0); q1, c = Add64(q1, t1, c)
	q3, t0 = Mul64(x3, y[2]); q2, c = Add64(q2, t0, c); q3, _ = Add64(q3,  0, c)

	t1, _  = Mul64(x0, y[2]); q0, c = Add64(q0, t1, 0)
	t1, t0 = Mul64(x2, y[2]); q1, c = Add64(q1, t0, c); q2, c = Add64(q2, t1, c); q3, _ = Add64(q3, 0, c)


	t1, t0 = Mul64(x1, y[3]); q1, c = Add64(q1, t0, 0); q2, c = Add64(q2, t1, c)
	q4, t0 = Mul64(x3, y[3]); q3, c = Add64(q3, t0, c); q4, _ = Add64(q4,  0, c)

	t1, t0 = Mul64(x0, y[3]); q0, c = Add64(q0, t0, 0); q1, c = Add64(q1, t1, c)
	t1, t0 = Mul64(x2, y[3]); q2, c = Add64(q2, t0, c); q3, c = Add64(q3, t1, c); q4, _ = Add64(q4, 0, c)

	// subtract: t3 = 2^191/y - 2^190/y = 2^190/y
	_,   b  = Sub64(  0, q0, 0)
	_,   b  = Sub64(  0, q1, b)
	t3l, b := Sub64(  0, q2, b)
	t3m, b := Sub64(r2l, q3, b)
	t3h, _ := Sub64(r2h, q4, b)

	// double: r3 = 2^191/y
	r3l, c := Add64(t3l, t3l, 0)
	r3m, c := Add64(t3m, t3m, c)
	r3h, _ := Add64(t3h, t3h, c)

	// Fourth iteration: 192 -> 320

	// square r3

	a4h, a4l := Mul64(r3l, r3l)
	b4h, b4l := Mul64(r3l, r3m)
	c4h, c4l := Mul64(r3l, r3h)
	d4h, d4l := Mul64(r3m, r3m)
	e4h, e4l := Mul64(r3m, r3h)
	f4h, f4l := Mul64(r3h, r3h)

	b4h, c = Add64(b4h, c4l, 0)
	e4l, c = Add64(e4l, c4h, c)
	e4h, _ = Add64(e4h,   0, c)

	a4h, c = Add64(a4h, b4l, 0)
	d4l, c = Add64(d4l, b4h, c)
	d4h, c = Add64(d4h, e4l, c)
	f4l, c = Add64(f4l, e4h, c)
	f4h, _ = Add64(f4h,   0, c)

	a4h, c = Add64(a4h, b4l, 0)
	d4l, c = Add64(d4l, b4h, c)
	d4h, c = Add64(d4h, e4l, c)
	f4l, c = Add64(f4l, e4h, c)
	f4h, _ = Add64(f4h,   0, c)

	// multiply by y

	x1, x0  = Mul64(d4h, y[0])
	x3, x2  = Mul64(f4h, y[0])
	t1, t0  = Mul64(f4l, y[0]); x1, c = Add64(x1, t0, 0); x2, c = Add64(x2, t1, c)
							      x3, _ = Add64(x3,  0, c)

	t1, t0  = Mul64(d4h, y[1]); x1, c = Add64(x1, t0, 0); x2, c = Add64(x2, t1, c)
	x4, t0 := Mul64(f4h, y[1]); x3, c = Add64(x3, t0, c); x4, _ = Add64(x4,  0, c)
	t1, t0  = Mul64(d4l, y[1]); x0, c = Add64(x0, t0, 0); x1, c = Add64(x1, t1, c)
	t1, t0  = Mul64(f4l, y[1]); x2, c = Add64(x2, t0, c); x3, c = Add64(x3, t1, c)
							      x4, _ = Add64(x4,  0, c)

	t1, t0  = Mul64(a4h, y[2]); x0, c = Add64(x0, t0, 0); x1, c = Add64(x1, t1, c)
	t1, t0  = Mul64(d4h, y[2]); x2, c = Add64(x2, t0, c); x3, c = Add64(x3, t1, c)
	x5, t0 := Mul64(f4h, y[2]); x4, c = Add64(x4, t0, c); x5, _ = Add64(x5,  0, c)
	t1, t0  = Mul64(d4l, y[2]); x1, c = Add64(x1, t0, 0); x2, c = Add64(x2, t1, c)
	t1, t0  = Mul64(f4l, y[2]); x3, c = Add64(x3, t0, c); x4, c = Add64(x4, t1, c)
							      x5, _ = Add64(x5,  0, c)

	t1, t0  = Mul64(a4h, y[3]); x1, c = Add64(x1, t0, 0); x2, c = Add64(x2, t1, c)
	t1, t0  = Mul64(d4h, y[3]); x3, c = Add64(x3, t0, c); x4, c = Add64(x4, t1, c)
	x6, t0 := Mul64(f4h, y[3]); x5, c = Add64(x5, t0, c); x6, _ = Add64(x6,  0, c)
	t1, t0  = Mul64(a4l, y[3]); x0, c = Add64(x0, t0, 0); x1, c = Add64(x1, t1, c)
	t1, t0  = Mul64(d4l, y[3]); x2, c = Add64(x2, t0, c); x3, c = Add64(x3, t1, c)
	t1, t0  = Mul64(f4l, y[3]); x4, c = Add64(x4, t0, c); x5, c = Add64(x5, t1, c)
							      x6, _ = Add64(x6,  0, c)

	// subtract
	_,   b	 = Sub64(  0, x0, 0)
	_,   b	 = Sub64(  0, x1, b)
	r4l, b	:= Sub64(  0, x2, b)
	r4k, b	:= Sub64(  0, x3, b)
	r4j, b	:= Sub64(r3l, x4, b)
	r4i, b	:= Sub64(r3m, x5, b)
	r4h, _	:= Sub64(r3h, x6, b)

	// Multiply candidate for 1/4y by y, with full precision

	x0 = r4l
	x1 = r4k
	x2 = r4j
	x3 = r4i
	x4 = r4h

	q1, q0	 = Mul64(x0, y[0])
	q3, q2	 = Mul64(x2, y[0])
	q5, q4	:= Mul64(x4, y[0])

	t1, t0	 = Mul64(x1, y[0]); q1, c = Add64(q1, t0, 0); q2, c = Add64(q2, t1, c)
	t1, t0	 = Mul64(x3, y[0]); q3, c = Add64(q3, t0, c); q4, c = Add64(q4, t1, c); q5, _ = Add64(q5, 0, c)

	t1, t0	 = Mul64(x0, y[1]); q1, c = Add64(q1, t0, 0); q2, c = Add64(q2, t1, c)
	t1, t0	 = Mul64(x2, y[1]); q3, c = Add64(q3, t0, c); q4, c = Add64(q4, t1, c)
	q6, t0	:= Mul64(x4, y[1]); q5, c = Add64(q5, t0, c); q6, _ = Add64(q6,  0, c)

	t1, t0	 = Mul64(x1, y[1]); q2, c = Add64(q2, t0, 0); q3, c = Add64(q3, t1, c)
	t1, t0	 = Mul64(x3, y[1]); q4, c = Add64(q4, t0, c); q5, c = Add64(q5, t1, c); q6, _ = Add64(q6, 0, c)

	t1, t0	 = Mul64(x0, y[2]); q2, c = Add64(q2, t0, 0); q3, c = Add64(q3, t1, c)
	t1, t0	 = Mul64(x2, y[2]); q4, c = Add64(q4, t0, c); q5, c = Add64(q5, t1, c)
	q7, t0	:= Mul64(x4, y[2]); q6, c = Add64(q6, t0, c); q7, _ = Add64(q7,  0, c)

	t1, t0	 = Mul64(x1, y[2]); q3, c = Add64(q3, t0, 0); q4, c = Add64(q4, t1, c)
	t1, t0	 = Mul64(x3, y[2]); q5, c = Add64(q5, t0, c); q6, c = Add64(q6, t1, c); q7, _ = Add64(q7, 0, c)

	t1, t0	 = Mul64(x0, y[3]); q3, c = Add64(q3, t0, 0); q4, c = Add64(q4, t1, c)
	t1, t0	 = Mul64(x2, y[3]); q5, c = Add64(q5, t0, c); q6, c = Add64(q6, t1, c)
	q8, t0	:= Mul64(x4, y[3]); q7, c = Add64(q7, t0, c); q8, _ = Add64(q8,  0, c)

	t1, t0	 = Mul64(x1, y[3]); q4, c = Add64(q4, t0, 0); q5, c = Add64(q5, t1, c)
	t1, t0	 = Mul64(x3, y[3]); q6, c = Add64(q6, t0, c); q7, c = Add64(q7, t1, c); q8, _ = Add64(q8, 0, c)

	// Final adjustments: increment/decrement the result to get the correct reciprocal

	// subtract q from 1/4
	q0, b = Sub64(0, q0, 0)
	q1, b = Sub64(0, q1, b)
	q2, b = Sub64(0, q2, b)
	q3, b = Sub64(0, q3, b)
	q4, b = Sub64(0, q4, b)
	q5, b = Sub64(0, q5, b)
	q6, b = Sub64(0, q6, b)
	q7, b = Sub64(0, q7, b)
	q8, b = Sub64(uint64(1) << 62, q8, b)

	// decrement the result
	x0, t := Sub64(r4l, 1, 0)
	x1, t  = Sub64(r4k, 0, t)
	x2, t  = Sub64(r4j, 0, t)
	x3, t  = Sub64(r4i, 0, t)
	x4, _  = Sub64(r4h, 0, t)

	// commit the decrement if the subtraction underflowed (reciprocal was too large)
	if b != 0 {
		r4h, r4i, r4j, r4k, r4l = x4, x3, x2, x1, x0
	}

	// subtract y from q
	q0, b = Sub64(q0, y[0], 0)
	q1, b = Sub64(q1, y[1], b)
	q2, b = Sub64(q2, y[2], b)
	q3, b = Sub64(q3, y[3], b)
	q4, b = Sub64(q4,    0, b)
	q5, b = Sub64(q5,    0, b)
	q6, b = Sub64(q6,    0, b)
	q7, b = Sub64(q7,    0, b)
	q8, b = Sub64(q8,    0, b)

	// increment the result
	x0, t = Add64(r4l, 1, 0)
	x1, t = Add64(r4k, 0, t)
	x2, t = Add64(r4j, 0, t)
	x3, t = Add64(r4i, 0, t)
	x4, _ = Add64(r4h, 0, t)

	// commit the increment if the subtraction did not underflow (reciprocal was too small)
	if b == 0 {
		r4h, r4i, r4j, r4k, r4l = x4, x3, x2, x1, x0
	}

	// Shift to correct bit alignment, truncating excess bits

	p = p - 1

	x0, c = Add64(r4l, r4l, 0)
	x1, c = Add64(r4k, r4k, c)
	x2, c = Add64(r4j, r4j, c)
	x3, c = Add64(r4i, r4i, c)
	x4, _ = Add64(r4h, r4h, c)

	if p < 0 {
		r4h, r4i, r4j, r4k, r4l = x4, x3, x2, x1, x0
		p = 0	// avoid negative shift below
	}

	// Shift right 0-62 bits
	{
		r := uint(p)		// right shift
		l := uint(64 - r)	// left shift

		x0 = (r4l >> r) | (r4k << l)
		x1 = (r4k >> r) | (r4j << l)
		x2 = (r4j >> r) | (r4i << l)
		x3 = (r4i >> r) | (r4h << l)
		x4 = (r4h >> r)
	}

	if p > 0 {
		r4h, r4i, r4j, r4k, r4l = x4, x3, x2, x1, x0
	}

	mu[0] = r4l
	mu[1] = r4k
	mu[2] = r4j
	mu[3] = r4i
	mu[4] = r4h

	return mu
}
