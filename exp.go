// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

// The ExpBase type contains lookup tables allowing fast repeated modular exponentiation with the same base value.
type ExpBase struct {
	t [16]Residue
}

// FromResidue initialises ExpBase from a residue.
// It performs 192 squarings and 11 multiplications.
func (z *ExpBase) FromResidue(x *Residue) *ExpBase {
	var t Residue

	t.Cpy(x)

	z.t[0].m = t.m
	z.t[0].r = [4]uint64{1,0,0,0}

	z.t[1].Cpy(&t)

	for i:=0; i<64; i++ {
		t.Sqr()
	}

	z.t[2].Cpy(&t)
	z.t[3].Cpy(&t).Mul(&z.t[1])

	for i:=0; i<64; i++ {
		t.Sqr()
	}

	z.t[4].Cpy(&t)
	z.t[5].Cpy(&t).Mul(&z.t[1])
	z.t[6].Cpy(&t).Mul(&z.t[2])
	z.t[7].Cpy(&t).Mul(&z.t[3])

	for i:=0; i<64; i++ {
		t.Sqr()
	}

	z.t[ 8].Cpy(&t)
	z.t[ 9].Cpy(&t).Mul(&z.t[1])
	z.t[10].Cpy(&t).Mul(&z.t[2])
	z.t[11].Cpy(&t).Mul(&z.t[3])
	z.t[12].Cpy(&t).Mul(&z.t[4])
	z.t[13].Cpy(&t).Mul(&z.t[5])
	z.t[14].Cpy(&t).Mul(&z.t[6])
	z.t[15].Cpy(&t).Mul(&z.t[7])

	return z
}

// ExpPrecomp takes an ExpBase computed from the base value, a 256-bit integer as the exponent, and performs modular exponentiation.
// It performs 63 squarings and 63 multiplications.
func (z *Residue) ExpPrecomp(x *ExpBase, y [4]uint64) *Residue {

	j :=	((y[3] >> 60) & 8) |
		((y[2] >> 61) & 4) |
		((y[1] >> 62) & 2) |
		((y[0] >> 63) & 1)

	z.Cpy(&x.t[j])

	for i := 1; i<64; i++ {
		y[3] <<= 1
		y[2] <<= 1
		y[1] <<= 1
		y[0] <<= 1

		j =	((y[3] >> 60) & 8) |
			((y[2] >> 61) & 4) |
			((y[1] >> 62) & 2) |
			((y[0] >> 63) & 1)

		z.Sqr().Mul(&x.t[j])
	}

	return z
}

// Exp performs modular exponentiation without storing precomputed values for later use.
// It performs 255 squarings and 74 multiplications.
func (z *Residue) Exp(x [4]uint64) *Residue {
	var eb ExpBase

	eb.FromResidue(z)
	z.ExpPrecomp(&eb, x)

	return z
}
