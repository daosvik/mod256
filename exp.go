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

	r.Copy(x)

	z.l[0].m = r.m
	z.l[0].r = [4]uint64{1,0,0,0}

	z.l[1].Copy(&r)

	for i:=0; i<32; i++ {
		r.Square()
	}

	z.l[2].Copy(&r)
	z.l[3].Copy(&r).Mul(&z.l[1])

	for i:=0; i<32; i++ {
		r.Square()
	}

	z.l[4].Copy(&r)
	z.l[5].Copy(&r).Mul(&z.l[1])
	z.l[6].Copy(&r).Mul(&z.l[2])
	z.l[7].Copy(&r).Mul(&z.l[3])

	for i:=0; i<32; i++ {
		r.Square()
	}

	z.l[ 8].Copy(&r)
	z.l[ 9].Copy(&r).Mul(&z.l[1])
	z.l[10].Copy(&r).Mul(&z.l[2])
	z.l[11].Copy(&r).Mul(&z.l[3])
	z.l[12].Copy(&r).Mul(&z.l[4])
	z.l[13].Copy(&r).Mul(&z.l[5])
	z.l[14].Copy(&r).Mul(&z.l[6])
	z.l[15].Copy(&r).Mul(&z.l[7])

	for i:=0; i<32; i++ {
		r.Square()
	}

	z.h[0].Copy(&z.l[0])

	z.h[1].Copy(&r)

	for i:=0; i<32; i++ {
		r.Square()
	}

	z.h[2].Copy(&r)
	z.h[3].Copy(&r).Mul(&z.h[1])

	for i:=0; i<32; i++ {
		r.Square()
	}

	z.h[4].Copy(&r)
	z.h[5].Copy(&r).Mul(&z.h[1])
	z.h[6].Copy(&r).Mul(&z.h[2])
	z.h[7].Copy(&r).Mul(&z.h[3])

	for i:=0; i<32; i++ {
		r.Square()
	}

	z.h[ 8].Copy(&r)
	z.h[ 9].Copy(&r).Mul(&z.h[1])
	z.h[10].Copy(&r).Mul(&z.h[2])
	z.h[11].Copy(&r).Mul(&z.h[3])
	z.h[12].Copy(&r).Mul(&z.h[4])
	z.h[13].Copy(&r).Mul(&z.h[5])
	z.h[14].Copy(&r).Mul(&z.h[6])
	z.h[15].Copy(&r).Mul(&z.h[7])

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

	z.Copy(&x.h[h]).Mul(&x.l[l])

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

		z.Square().Mul(&x.h[h]).Mul(&x.l[l])
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

	r.Copy(z)

	t[0].m = r.m
	t[0].r = [4]uint64{1,0,0,0}

	t[1].Copy(z)

	for i:=0; i<64; i++ {
		r.Square()
	}

	t[2].Copy(&r)
	t[3].Copy(&r).Mul(&t[1])

	for i:=0; i<64; i++ {
		r.Square()
	}

	t[4].Copy(&r)
	t[5].Copy(&r).Mul(&t[1])
	t[6].Copy(&r).Mul(&t[2])
	t[7].Copy(&r).Mul(&t[3])

	for i:=0; i<64; i++ {
		r.Square()
	}

	t[ 8].Copy(&r)
	t[ 9].Copy(&r).Mul(&t[1])
	t[10].Copy(&r).Mul(&t[2])
	t[11].Copy(&r).Mul(&t[3])
	t[12].Copy(&r).Mul(&t[4])
	t[13].Copy(&r).Mul(&t[5])
	t[14].Copy(&r).Mul(&t[6])
	t[15].Copy(&r).Mul(&t[7])

	y := x

	j :=	((y[3] >> 60) & 8) |
		((y[2] >> 61) & 4) |
		((y[1] >> 62) & 2) |
		((y[0] >> 63) & 1)

	z.Copy(&t[j])

	for i := 1; i<64; i++ {
		y[3] <<= 1
		y[2] <<= 1
		y[1] <<= 1
		y[0] <<= 1

		j =	((y[3] >> 60) & 8) |
			((y[2] >> 61) & 4) |
			((y[1] >> 62) & 2) |
			((y[0] >> 63) & 1)

		z.Square().Mul(&t[j])
	}

	return z
}
