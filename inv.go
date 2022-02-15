// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

// z = 1/x mod m if it exists, otherwise 0
// Returns true if the inverse exists

// Inv computes the (multiplicative) inverse of a residue, if it exists.
func (z *Residue) Inv() bool {

	var (
		b, c, // Borrow & carry
		a4, a3, a2, a1, a0,
		b4, b3, b2, b1, b0,
		c4, c3, c2, c1, c0,
		d4, d3, d2, d1, d0 uint64
	)

	x := z.r
	y := z.m.m

	u3, u2, u1, u0 := x[3], x[2], x[1], x[0]
	v3, v2, v1, v0 := y[3], y[2], y[1], y[0]

	if (u3 | u2 | u1 | u0) == 0 ||	// u == 0
	   (v3 | v2 | v1 | v0) == 0 ||	// v == 0
	   (u0 | v0) & 1 == 0 {		// 2|gcd(u,v)
		// there is no inverse
		z.r[3], z.r[2], z.r[1], z.r[0] = 0, 0, 0, 0
		return false
	}

	a4, a3, a2, a1, a0 = 0, 0, 0, 0, 1
	b4, b3, b2, b1, b0 = 0, 0, 0, 0, 0
	c4, c3, c2, c1, c0 = 0, 0, 0, 0, 0
	d4, d3, d2, d1, d0 = 0, 0, 0, 0, 1

	done := false

	for !done {
		for u0 & 1 == 0 {

			u0 = (u0 >> 1) | (u1 << 63)
			u1 = (u1 >> 1) | (u2 << 63)
			u2 = (u2 >> 1) | (u3 << 63)
			u3 = (u3 >> 1)

			if (a0 | b0) & 1 == 1 {

				a0, c = Add64(a0, y[0], 0)
				a1, c = Add64(a1, y[1], c)
				a2, c = Add64(a2, y[2], c)
				a3, c = Add64(a3, y[3], c)
				a4, _ = Add64(a4,    0, c)

				b0, b = Sub64(b0, x[0], 0)
				b1, b = Sub64(b1, x[1], b)
				b2, b = Sub64(b2, x[2], b)
				b3, b = Sub64(b3, x[3], b)
				b4, _ = Sub64(b4,    0, b)
			}

			a0 = (a0 >> 1) | (a1 << 63)
			a1 = (a1 >> 1) | (a2 << 63)
			a2 = (a2 >> 1) | (a3 << 63)
			a3 = (a3 >> 1) | (a4 << 63)
			a4 = uint64(int64(a4) >> 1)

			b0 = (b0 >> 1) | (b1 << 63)
			b1 = (b1 >> 1) | (b2 << 63)
			b2 = (b2 >> 1) | (b3 << 63)
			b3 = (b3 >> 1) | (b4 << 63)
			b4 = uint64(int64(b4) >> 1)
		}

		for v0 & 1 == 0 {

			v0 = (v0 >> 1) | (v1 << 63)
			v1 = (v1 >> 1) | (v2 << 63)
			v2 = (v2 >> 1) | (v3 << 63)
			v3 = (v3 >> 1)

			if (c0 | d0) & 1 == 1 {

				c0, c = Add64(c0, y[0], 0)
				c1, c = Add64(c1, y[1], c)
				c2, c = Add64(c2, y[2], c)
				c3, c = Add64(c3, y[3], c)
				c4, _ = Add64(c4,    0, c)

				d0, b = Sub64(d0, x[0], 0)
				d1, b = Sub64(d1, x[1], b)
				d2, b = Sub64(d2, x[2], b)
				d3, b = Sub64(d3, x[3], b)
				d4, _ = Sub64(d4,    0, b)
			}

			c0 = (c0 >> 1) | (c1 << 63)
			c1 = (c1 >> 1) | (c2 << 63)
			c2 = (c2 >> 1) | (c3 << 63)
			c3 = (c3 >> 1) | (c4 << 63)
			c4 = uint64(int64(c4) >> 1)

			d0 = (d0 >> 1) | (d1 << 63)
			d1 = (d1 >> 1) | (d2 << 63)
			d2 = (d2 >> 1) | (d3 << 63)
			d3 = (d3 >> 1) | (d4 << 63)
			d4 = uint64(int64(d4) >> 1)
		}

		t0, b := Sub64(u0, v0, 0)
		t1, b := Sub64(u1, v1, b)
		t2, b := Sub64(u2, v2, b)
		t3, b := Sub64(u3, v3, b)

		if b == 0 { // u >= v

			u3, u2, u1, u0 = t3, t2, t1, t0

			a0, b = Sub64(a0, c0, 0)
			a1, b = Sub64(a1, c1, b)
			a2, b = Sub64(a2, c2, b)
			a3, b = Sub64(a3, c3, b)
			a4, _ = Sub64(a4, c4, b)

			b0, b = Sub64(b0, d0, 0)
			b1, b = Sub64(b1, d1, b)
			b2, b = Sub64(b2, d2, b)
			b3, b = Sub64(b3, d3, b)
			b4, _ = Sub64(b4, d4, b)

		} else { // v > u

			v0, b = Sub64(v0, u0, 0)
			v1, b = Sub64(v1, u1, b)
			v2, b = Sub64(v2, u2, b)
			v3, _ = Sub64(v3, u3, b)

			c0, b = Sub64(c0, a0, 0)
			c1, b = Sub64(c1, a1, b)
			c2, b = Sub64(c2, a2, b)
			c3, b = Sub64(c3, a3, b)
			c4, _ = Sub64(c4, a4, b)

			d0, b = Sub64(d0, b0, 0)
			d1, b = Sub64(d1, b1, b)
			d2, b = Sub64(d2, b2, b)
			d3, b = Sub64(d3, b3, b)
			d4, _ = Sub64(d4, b4, b)
		}

		if (u3 | u2 | u1 | u0) == 0 {
			done = true
		}
	}

	if (v3 | v2 | v1 | (v0 - 1)) != 0 { // gcd(z,m) != 1
		z.r[3], z.r[2], z.r[1], z.r[0] = 0, 0, 0, 0
		return false
	}

	// Add or subtract modulus to find 256-bit inverse (<= 2 iterations expected)

	for (c4 >> 63) != 0 {
		c0, c = Add64(c0, y[0], 0)
		c1, c = Add64(c1, y[1], c)
		c2, c = Add64(c2, y[2], c)
		c3, c = Add64(c3, y[3], c)
		c4, _ = Add64(c4,    0, c)
	}

	for c4 != 0 {
		c0, b = Sub64(c0, y[0], 0)
		c1, b = Sub64(c1, y[1], b)
		c2, b = Sub64(c2, y[2], b)
		c3, b = Sub64(c3, y[3], b)
		c4, _ = Sub64(c4,    0, b)
	}

	z.r[3], z.r[2], z.r[1], z.r[0] = c3, c2, c1, c0

	return true
}
