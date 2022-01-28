// mod256: Arithmetic modulo 193-256 bit moduli 
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

func (z *Residue) Add(x *Residue) (*Residue) {
	if z.m != x.m {
		if z.m.m != x.m.m {
			panic("Incompatible moduli")
		}
	}

	t0, c := Add64(z.r[0], x.r[0], 0)
	t1, c := Add64(z.r[1], x.r[1], c)
	t2, c := Add64(z.r[2], x.r[2], c)
	t3, c := Add64(z.r[3], x.r[3], c)

	u0, b := Sub64(t0, z.m.mmu1[0], 0)
	u1, b := Sub64(t1, z.m.mmu1[1], b)
	u2, b := Sub64(t2, z.m.mmu1[2], b)
	u3, b := Sub64(t3, z.m.mmu1[3], b)

	v0, b := Sub64(t0, z.m.mmu0[0], 0)
	v1, b := Sub64(t1, z.m.mmu0[1], b)
	v2, b := Sub64(t2, z.m.mmu0[2], b)
	v3, b := Sub64(t3, z.m.mmu0[3], b)

	// Subtract the larger multiple of m if necessary

	if b == 0 {
		v3, v2, v1, v0 = u3, u2, u1, u0
	}

	// Subtract if overflow

	if c != 0 {
		t3, t2, t1, t0 = v3, v2, v1, v0
	}

	z.r[3], z.r[2], z.r[1], z.r[0] = t3, t2, t1, t0

	return z
}
