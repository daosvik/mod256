// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

// The ExpBase type contains lookup tables allowing fast repeated modular exponentiation with the same base value.
type ExpBase struct {
	h, l [16]Residue
}

// FromResidue initialises ExpBase from a residue.
// It performs 224 squarings and 22 multiplications.
func (z *ExpBase) FromResidue(x *Residue) *ExpBase {
	var r Residue

	r.Cpy(x)

	z.l[0].m = r.m
	z.l[0].r = [4]uint64{1,0,0,0}

	z.l[1].Cpy(&r)

	for i:=0; i<32; i++ {
		r.Sqr()
	}

	z.l[2].Cpy(&r)
	z.l[3].Cpy(&r).Mul(&z.l[1])

	for i:=0; i<32; i++ {
		r.Sqr()
	}

	z.l[4].Cpy(&r)
	z.l[5].Cpy(&r).Mul(&z.l[1])
	z.l[6].Cpy(&r).Mul(&z.l[2])
	z.l[7].Cpy(&r).Mul(&z.l[3])

	for i:=0; i<32; i++ {
		r.Sqr()
	}

	z.l[ 8].Cpy(&r)
	z.l[ 9].Cpy(&r).Mul(&z.l[1])
	z.l[10].Cpy(&r).Mul(&z.l[2])
	z.l[11].Cpy(&r).Mul(&z.l[3])
	z.l[12].Cpy(&r).Mul(&z.l[4])
	z.l[13].Cpy(&r).Mul(&z.l[5])
	z.l[14].Cpy(&r).Mul(&z.l[6])
	z.l[15].Cpy(&r).Mul(&z.l[7])

	for i:=0; i<32; i++ {
		r.Sqr()
	}

	z.h[0].Cpy(&z.l[0])

	z.h[1].Cpy(&r)

	for i:=0; i<32; i++ {
		r.Sqr()
	}

	z.h[2].Cpy(&r)
	z.h[3].Cpy(&r).Mul(&z.h[1])

	for i:=0; i<32; i++ {
		r.Sqr()
	}

	z.h[4].Cpy(&r)
	z.h[5].Cpy(&r).Mul(&z.h[1])
	z.h[6].Cpy(&r).Mul(&z.h[2])
	z.h[7].Cpy(&r).Mul(&z.h[3])

	for i:=0; i<32; i++ {
		r.Sqr()
	}

	z.h[ 8].Cpy(&r)
	z.h[ 9].Cpy(&r).Mul(&z.h[1])
	z.h[10].Cpy(&r).Mul(&z.h[2])
	z.h[11].Cpy(&r).Mul(&z.h[3])
	z.h[12].Cpy(&r).Mul(&z.h[4])
	z.h[13].Cpy(&r).Mul(&z.h[5])
	z.h[14].Cpy(&r).Mul(&z.h[6])
	z.h[15].Cpy(&r).Mul(&z.h[7])

	return z
}

// ExpPrecomp takes an ExpBase computed from the base value, a 256-bit integer as the exponent, and performs modular exponentiation.
// It performs 31 squarings and 63 multiplications.
func (z *Residue) ExpPrecomp(x *ExpBase, y [4]uint64) *Residue {

	h :=	((y[3] >> 60) & 8) |
		((y[3] >> 29) & 4) |
		((y[2] >> 62) & 2) |
		((y[2] >> 31) & 1)

	l :=	((y[1] >> 60) & 8) |
		((y[1] >> 29) & 4) |
		((y[0] >> 62) & 2) |
		((y[0] >> 31) & 1)

	z.Cpy(&x.h[h]).Mul(&x.l[l])

	for i := 1; i<32; i++ {
		y[3] <<= 1
		y[2] <<= 1
		y[1] <<= 1
		y[0] <<= 1

		h =	((y[3] >> 60) & 8) |
			((y[3] >> 29) & 4) |
			((y[2] >> 62) & 2) |
			((y[2] >> 31) & 1)

		l =	((y[1] >> 60) & 8) |
			((y[1] >> 29) & 4) |
			((y[0] >> 62) & 2) |
			((y[0] >> 31) & 1)

		z.Sqr().Mul(&x.h[h]).Mul(&x.l[l])
	}

	return z
}

// Exp performs modular exponentiation without storing precomputed values for later use.
// It performs 255 squarings and 74 multiplications.
func (z *Residue) Exp(x [4]uint64) *Residue {
	var (
		r Residue
		t [16]Residue
	)

	r.Cpy(z)

	t[0].m = r.m
	t[0].r = [4]uint64{1,0,0,0}

	t[1].Cpy(z)

	for i:=0; i<64; i++ {
		r.Sqr()
	}

	t[2].Cpy(&r)
	t[3].Cpy(&r).Mul(&t[1])

	for i:=0; i<64; i++ {
		r.Sqr()
	}

	t[4].Cpy(&r)
	t[5].Cpy(&r).Mul(&t[1])
	t[6].Cpy(&r).Mul(&t[2])
	t[7].Cpy(&r).Mul(&t[3])

	for i:=0; i<64; i++ {
		r.Sqr()
	}

	t[ 8].Cpy(&r)
	t[ 9].Cpy(&r).Mul(&t[1])
	t[10].Cpy(&r).Mul(&t[2])
	t[11].Cpy(&r).Mul(&t[3])
	t[12].Cpy(&r).Mul(&t[4])
	t[13].Cpy(&r).Mul(&t[5])
	t[14].Cpy(&r).Mul(&t[6])
	t[15].Cpy(&r).Mul(&t[7])

	y := x

	j :=	((y[3] >> 60) & 8) |
		((y[2] >> 61) & 4) |
		((y[1] >> 62) & 2) |
		((y[0] >> 63) & 1)

	z.Cpy(&t[j])

	for i := 1; i<64; i++ {
		y[3] <<= 1
		y[2] <<= 1
		y[1] <<= 1
		y[0] <<= 1

		j =	((y[3] >> 60) & 8) |
			((y[2] >> 61) & 4) |
			((y[1] >> 62) & 2) |
			((y[0] >> 63) & 1)

		z.Sqr().Mul(&t[j])
	}

	return z
}
