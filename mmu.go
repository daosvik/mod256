// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

func mmu0(z *Modulus) {
	var c, t0, t1, q0, q1, q2, q3, q4 uint64

	q2, q1 = Mul64(z.m[1], z.mu[4])
	q4, q3 = Mul64(z.m[3], z.mu[4])

	t1, q0 = Mul64(z.m[0], z.mu[4]); q1, c = Add64(q1, t1, 0)
	t1, t0 = Mul64(z.m[2], z.mu[4]); q2, c = Add64(q2, t0, c); q3, c = Add64(q3, t1, c); q4, _ = Add64(q4, 0, c)

	if q4 != 0 {
		panic("Error preparing mmu0")
	}

	z.mmu0 = [4]uint64{ q0, q1, q2, q3 }

	return
}

func mmu1(z *Modulus) {
	var c uint64

	z.mmu1[0], c = Add64(z.mmu0[0], z.m[0], 0)
	z.mmu1[1], c = Add64(z.mmu0[1], z.m[1], c)
	z.mmu1[2], c = Add64(z.mmu0[2], z.m[2], c)
	z.mmu1[3], c = Add64(z.mmu0[3], z.m[3], c)

	// mmu0 is the largest multiple of m such that mmu0 < 2^256
	// mmu1 is the smallest multiple of m such that mmu1 >= 2^256
	// => there must be a carry out from the addition above

	if c != 1 {
		panic("Error preparing mmu1")
	}

	return
}
