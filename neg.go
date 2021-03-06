// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	. "math/bits"
)

// Neg computes the negation (additive inverse) of a residue.
func (z *Residue) Neg() *Residue {
	t0, b := Sub64(z.m.mmu0[0], z.r[0], 0)
	t1, b := Sub64(z.m.mmu0[1], z.r[1], b)
	t2, b := Sub64(z.m.mmu0[2], z.r[2], b)
	t3, b := Sub64(z.m.mmu0[3], z.r[3], b)

	if b == 0 {
		z.r[3], z.r[2], z.r[1], z.r[0] = t3, t2, t1, t0
		return z
	}

	t0, b = Sub64(z.m.mmu1[0], z.r[0], 0)
	t1, b = Sub64(z.m.mmu1[1], z.r[1], b)
	t2, b = Sub64(z.m.mmu1[2], z.r[2], b)
	t3, _ = Sub64(z.m.mmu1[3], z.r[3], b)

	z.r[3], z.r[2], z.r[1], z.r[0] = t3, t2, t1, t0

	return z
}
