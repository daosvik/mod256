// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

func (z *Residue) Sub(x *Residue) *Residue {
	if z.m != x.m {
		if z.m.m != x.m.m {
			panic("Incompatible moduli")
		}
	}

	t0, b := Sub64(z.r[0], x.r[0], 0)
	t1, b := Sub64(z.r[1], x.r[1], b)
	t2, b := Sub64(z.r[2], x.r[2], b)
	t3, b := Sub64(z.r[3], x.r[3], b)

	u0, c := Add64(t0, z.m.mmu1[0], 0)
	u1, c := Add64(t1, z.m.mmu1[1], c)
	u2, c := Add64(t2, z.m.mmu1[2], c)
	u3, _ := Add64(t3, z.m.mmu1[3], c)

	v0, c := Add64(t0, z.m.mmu0[0], 0)
	v1, c := Add64(t1, z.m.mmu0[1], c)
	v2, c := Add64(t2, z.m.mmu0[2], c)
	v3, c := Add64(t3, z.m.mmu0[3], c)

	// Add the larger multiple of m if necessary

	if c == 0 {
		v3, v2, v1, v0 = u3, u2, u1, u0
	}

	// Add if underflow

	if b != 0 {
		t3, t2, t1, t0 = v3, v2, v1, v0
	}

	z.r[3], z.r[2], z.r[1], z.r[0] = t3, t2, t1, t0

	return z
}
