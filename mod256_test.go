// mod256: Arithmetic modulo 193-256 bit moduli
// Copyright 2021-2022 Dag Arne Osvik
// SPDX-License-Identifier: BSD-3-Clause

package mod256

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"math/bits"
	"testing"
)

// Test values are grouped as follows
// 1 bit set	[256]	(64 m values)
// 255 bits set	[256]	(256 m values)
// 0 bits set	[  1]	(0 m values)
// 256 bits set	[  1]	(1 m value)
// Random values[126]	(<= 126 m values)
//
// Only values with m[3] > 0 are used as moduli.

var (
	testval                           [640][4]uint64
	test_all, test_fixed, test_random [][4]uint64
)

func init() {

	nistp256[3] = 0xffffffff00000001
	nistp256[2] = 0x0000000000000000
	nistp256[1] = 0x00000000ffffffff
	nistp256[0] = 0xffffffffffffffff

	nistp224[3] = 0x00000000ffffffff
	nistp224[2] = 0xffffffffffffffff
	nistp224[1] = 0xffffffff00000000
	nistp224[0] = 0x0000000000000001

	test_all	= testval[:]
	test_fixed	= testval[0:514]
	test_random	= testval[514:]

	i := 0

	testval[i] = [4]uint64{0, 0, 0, 0}
	i++

	testval[i] = [4]uint64{0, 0, 0, 0}
	{
		tmp := testval[i-1]
		tmp[0] = ^tmp[0]
		tmp[1] = ^tmp[1]
		tmp[2] = ^tmp[2]
		tmp[3] = ^tmp[3]
		testval[i] = tmp
	}
	i++

	single := i

	for j:=0; j<64; j++ {
		testval[i] = [4]uint64{1 << j, 0, 0, 0}
		i++
	}

	for j:=0; j<64; j++ {
		testval[i] = [4]uint64{0, 1 << j, 0, 0}
		i++
	}

	for j:=0; j<64; j++ {
		testval[i] = [4]uint64{0, 0, 1 << j, 0}
		i++
	}

	for j:=0; j<64; j++ {
		testval[i] = [4]uint64{0, 0, 0, 1 << j}
		i++
	}

	for j:=0; j<256; j++ {
		tmp := testval[single+j]
		tmp[0] = ^tmp[0]
		tmp[1] = ^tmp[1]
		tmp[2] = ^tmp[2]
		tmp[3] = ^tmp[3]
		testval[i] = tmp
		i++
	}

	for i < cap(testval) {
		var (
			b [32]byte
			w [4]uint64
		)

		n, _ := rand.Read(b[:])

		if n < 32 {
			fmt.Print("Could not read 32-byte random value\n")
		}

		// Convert [32]byte to [4]uint64

		for j:=0; j<4; j++ {
			for k:=0; k<8; k++ {
				w[j] <<= 8
				w[j] |= uint64(b[4*j+k])
			}
		}

		testval[i] = w
		i++
	}
}

func TestModulusFromToUint64_OK(t *testing.T) {
	test_mod := test_all

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		x := NewModulusFromUint64(m).ToUint64()

		if x != m {
			t.Fatalf("%v != %v", x, y)
		}
	}
}

func testMmu0Fail(t *testing.T, x [4]uint64) {
	var c uint64

	m := NewModulusFromUint64(x)

	defer func() {
		err := recover()
		if err != nil {
			return
		} else {
			t.Fatalf("Did not fail: %v\n%v", x, m.mu)
		}
	}()

	// increment mu[4] by 1
	m.mu[4], c = bits.Add64(m.mu[4], 1, 0)

	if c != 0 {
		// overflow, so this test does not apply
		// panic to satisfy failure requirement
		panic("")
	}

	mmu0(m)	// Should fail
}

func TestMmu0Fail(t *testing.T) {
	test_mod := test_all
	for _, m := range test_mod {
		if m[3] == 0 {
			continue
		}
		testMmu0Fail(t, m)
	}
}

func testMmu1Fail(t *testing.T, x [4]uint64) {
	var b uint64

	m := NewModulusFromUint64(x)

	defer func() {
		err := recover()
		if err != nil {
			return
		} else {
			t.Fatalf("Did not fail: %v\n%v", x, m.mu)
		}
	}()

	// decrement mu[4] by 1
	m.mu[4], b = bits.Sub64(m.mu[4], 1, 0)

	if b != 0 {
		// overflow, so this test does not apply
		// panic to satisfy failure requirement
		panic("")
	}

	mmu0(m)	// Should not fail
	mmu1(m)	// Should fail
}

func TestMmu1Fail(t *testing.T) {
	test_mod := test_all
	for _, m := range test_mod {
		if m[3] == 0 {
			continue
		}
		testMmu1Fail(t, m)
	}
}

func testResidueFromUint64_OK() {
	var r Residue

	m := NewModulusFromUint64([4]uint64 { 0, 1, 2, 3 })
	r.FromUint64(m, [4]uint64 { 4, 5, 6, 7 })
}

func TestResidueFromUint64_OK(t *testing.T) {
	defer func() {
		x := recover()
		if x == nil {
			return
		} else {
			t.Fatalf("Did not succeed")
		}
	}()

	testResidueFromUint64_OK()
}

func testResidueFromUint64_NoInit() {
	var (
		m *Modulus
		r Residue
	)

	// m is uninitialised
	r.FromUint64(m, [4]uint64 { 0, 1, 2, 3 })
}

func TestResidueFromUint64_NoInit(t *testing.T) {
	defer func() {
		x := recover()
		if x != nil {
			return
		} else {
			t.Fatalf("Did not fail")
		}
	}()

	testResidueFromUint64_NoInit()
}

func testResidueFromUint64_TooSmallModulus(m *Modulus) {

	m = NewModulusFromUint64([4]uint64 { 3, 2, 1, 0 })
}

func TestResidueFromUint64_TooSmallModulus(t *testing.T) {
	var m *Modulus

	defer func() {
		x := recover()
		if x != nil {
			return
		} else {
			t.Fatalf("Did not fail")
		}
	}()

	testResidueFromUint64_TooSmallModulus(m)
}

func requireSuccess(t *testing.T, f func(a, b *Residue), a, b *Residue) {
//	if err t.Fatalf
//	else return
	defer func() {
		x := recover()
		if x == nil {
			return
		} else {
			t.Fatalf("No success (m1 = %016x%016x%016x%016x, m2 = %016x%016x%016x%016x)\n",
			a.m.m[3], a.m.m[2], a.m.m[1], a.m.m[0], b.m.m[3], b.m.m[2], b.m.m[1], b.m.m[0])
		}
	}()
	f(a,b)
}

func requireFailure(t *testing.T, f func(a, b *Residue), a, b *Residue) {
//	if err recover and return
//	else t.Fatalf
	defer func() {
		x := recover()
		if x != nil {
			return
		} else {
			t.Fatalf("No failure (m1 = %016x%016x%016x%016x, m2 = %016x%016x%016x%016x)\n",
			a.m.m[3], a.m.m[2], a.m.m[1], a.m.m[0], b.m.m[3], b.m.m[2], b.m.m[1], b.m.m[0])
		}
	}()
	f(a,b)
}

func testAdd(a, b *Residue) {
	b.Add(a)
}

func testSub(a, b *Residue) {
	b.Sub(a)
}

func testMul(a, b *Residue) {
	b.Mul(a)
}

func testExp(a, b *Residue) {
	b.Exp(a.r)
}

func TestResidueCompatibility(t *testing.T) {
	var (
		r1, r2 Residue
		count  int
	)

	test_mod := test_fixed
	test_ops := test_random

	for i, m := range test_mod {
		if m[3] == 0 {
			continue
		}

		m1 := NewModulusFromUint64(m)

		for _, a := range test_ops {
			r1.FromUint64(m1, a)

			for j, n := range test_mod {
				if n[3] == 0 {
					continue
				}

				m2 := NewModulusFromUint64(n)
				r2.FromUint64(m2, a)

				if i==j {
					requireSuccess(t, testAdd, &r1, &r2)
					requireSuccess(t, testSub, &r1, &r2)
					requireSuccess(t, testMul, &r1, &r2)
				} else {
					requireFailure(t, testAdd, &r1, &r2)
					requireFailure(t, testSub, &r1, &r2)
					requireFailure(t, testMul, &r1, &r2)
				}
				count+=4
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestResidueFromToUint64(t *testing.T) {
	var (
		r           Residue
		bm, b, bmod big.Int
		count       int
	)

	test_mod := test_all
	test_ops := test_all

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		bm.SetString(fmt.Sprintf("%016x%016x%016x%016x", m[3], m[2], m[1], m[0]), 16)

		for _, a := range test_ops {
			r.FromUint64(mod, a).ToUint64()

			b.SetString(fmt.Sprintf("%016x%016x%016x%016x", a[3], a[2], a[1], a[0]), 16)
			bmod.Mod(&b, &bm)

			b.SetString(fmt.Sprintf("%016x%016x%016x%016x", r.r[3], r.r[2], r.r[1], r.r[0]), 16)

			if bmod.Cmp(&b) != 0 {
				t.Fatalf("%v != %v", bmod, b)
			}

			count++
		}
	}
	t.Logf("%v tests\n", count)
}

func TestReciprocal(t *testing.T) {
	check := func(m [4]uint64, e, mu [5]uint64) {
		if mu != e {
			t.Errorf("Computed 2**512//0x%016x%016x%016x%016x\n", m[3], m[2], m[1], m[0])
			t.Errorf("expected 0x%016x%016x%016x%016x%016x\n", e[4], e[3], e[2], e[1], e[0])
			t.Fatalf("received 0x%016x%016x%016x%016x%016x\n", mu[4], mu[3], mu[2], mu[1], mu[0])
		}
	}

	var (
		count                              int
		str_x, str_m, str_mu               string
		big_1, big_x, big_m, big_mu, big_e big.Int
		m                                  [4]uint64
		e, mu                              [5]uint64
	)

	m = [4]uint64{0x0000000000000001, 0x0000000000000000, 0x0000000000000000, 0x0000000000000001}
	e = [5]uint64{0, 0, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}
	mu = reciprocal(m)

	check(m, e, mu)

	m = [4]uint64{0x0000000000000000, 0x0000000000000000, 0xffffffffffffffff, 0xffffffffffffffff}
	e = [5]uint64{1, 0, 1, 0, 1}
	mu = reciprocal(m)

	check(m, e, mu)

	m = [4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffffe, 0xffffffffffffffff}

	for i:=uint64(2); i<134217728; i++ {
		e = [5]uint64{i, 0, 1, 0, 1}
		mu = reciprocal(m)

		check(m, e, mu)
		m[0]--

		count++
	}

	m = [4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffffd, 0xffffffffffffffff}
	e = [5]uint64{5, 0, 2, 0, 1}
	mu = reciprocal(m)

	check(m, e, mu)
	count++

	m = [4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffffc, 0xffffffffffffffff}
	e = [5]uint64{0xa, 0, 3, 0, 1}
	mu = reciprocal(m)

	check(m, e, mu)
	count++

	m = [4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffffb, 0xffffffffffffffff}
	e = [5]uint64{0x11, 0, 4, 0, 1}
	mu = reciprocal(m)

	check(m, e, mu)
	count++

	m = [4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffff7, 0xffffffffffffffff}
	e = [5]uint64{0x41, 0, 8, 0, 1}
	mu = reciprocal(m)

	check(m, e, mu)
	count++

	m = [4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffef, 0xffffffffffffffff}
	e = [5]uint64{0x101, 0, 0x10, 0, 1}
	mu = reciprocal(m)

	check(m, e, mu)
	count++

	m = [4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffff7f, 0xffffffffffffffff}
	e = [5]uint64{0x4001, 0, 0x80, 0, 1}
	mu = reciprocal(m)

	check(m, e, mu)
	count++

	test_mod := test_all

	big_1.SetUint64(1)

	str_x = fmt.Sprintf("%016x%016x%016x%016x%016x%016x%016x%016x%016x", 1, 0, 0, 0, 0, 0, 0, 0, 0)
	big_x.SetString(str_x,  16)

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mu = reciprocal(m)

		str_m  = fmt.Sprintf("%016x%016x%016x%016x", m[3], m[2], m[1], m[0])
		str_mu = fmt.Sprintf("%016x%016x%016x%016x%016x", mu[4], mu[3], mu[2], mu[1], mu[0])

		big_m.SetString(str_m,  16)
		big_mu.SetString(str_mu, 16)

		big_e.Div(&big_x, &big_m)

		if m[0] | m[1] | m[2] | (m[3] & (m[3]-1)) == 0 {
			// m is a power of 2
			// mu is then one less than 2^512/m
			big_e.Sub(&big_e, &big_1)
		}

		if big_e.Cmp(&big_mu) != 0 {
			t.Errorf("Computed 2**512 // %v\n", str_m)
			t.Errorf("expected %v\n", big_e.Text(16))
			t.Fatalf("received %v\n", big_mu.Text(16))
		}
		count++
	}

	t.Logf("%v tests\n", count)
}

func TestEqNeq(t *testing.T) {
	var (
		r             [4]Residue
		a, b, c, x, y Residue
		u, v          bool
		count         int
	)

	// All unequal

	m1 := NewModulusFromUint64([4]uint64 { 1, 2, 3, 4 })
	m2 := NewModulusFromUint64([4]uint64 { 5, 6, 7, 8 })

	r[0].FromUint64(m1, [4]uint64 { 1, 3, 5, 7 })
	r[1].FromUint64(m1, [4]uint64 { 2, 4, 6, 8 })

	r[2].FromUint64(m2, [4]uint64 { 1, 3, 5, 7 })
	r[3].FromUint64(m2, [4]uint64 { 2, 4, 6, 8 })

	for i:=0; i<4; i++ {
		for j:=0; j<4; j++ {
			x.Cpy(&r[i])
			y.Cpy(&r[j])

			if (i == j) && !x.Eq(&y) {
				t.Fatalf("Eq(%v, %v)", r[i], r[j])
			}

			x.Cpy(&r[i])
			y.Cpy(&r[j])

			if (i != j) && x.Eq(&y) {
				t.Fatalf("Eq(%v, %v)", r[i], r[j])
			}

			x.Cpy(&r[i])
			y.Cpy(&r[j])

			if (i == j) && x.Neq(&y) {
				t.Fatalf("Neq(%v, %v)", r[i], r[j])
			}

			x.Cpy(&r[i])
			y.Cpy(&r[j])

			if (i != j) && !x.Neq(&y) {
				t.Fatalf("Neq(%v, %v)", r[i], r[j])
			}
			count+=2
		}
	}

	// All equal

	r[0].FromUint64(m1, [4]uint64 { 1, 3,  5,  7 })
	r[1].FromUint64(m1, [4]uint64 { 2, 5,  8, 11 })
	r[2].FromUint64(m1, [4]uint64 { 3, 7, 11, 15 })
	r[3].FromUint64(m1, [4]uint64 { 4, 9, 14, 19 })

	for i:=0; i<4; i++ {
		for j:=0; j<4; j++ {
			x.Cpy(&r[i])
			y.Cpy(&r[j])

			if !x.Eq(&y) {
				t.Fatalf("Eq(%v, %v)", r[i], r[j])
			}

			x.Cpy(&r[i])
			y.Cpy(&r[j])

			if x.Neq(&y) {
				t.Fatalf("Neq(%v, %v)", r[i], r[j])
			}
			count+=2
		}
	}

	test_mod := test_all
	test_ops := test_fixed

	// a == a+m (mod m)
	// Eq(a,b) != Neq(a,b) for all a,b

	for i, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		m1 := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(m1, _a)

			for k, _b := range test_ops {
				b.FromUint64(m1, _b)
				c.Cpy(&a)

				u = c.Eq(&b)
				v = c.Neq(&b)

				if u == v {
					t.Errorf("a==b = %v", u)
					t.Errorf("a!=b = %v", v)
					t.Errorf("(%v,%v,%v)\n", i, j, k);
					t.Fatalf("%v\n%v\n%v\n", a, b, m);
				}

				count++
			}
		}
	}

	t.Logf("%v tests\n", count)
}

func TestReflexivity(t *testing.T) {
	var (
		a     Residue
		u     bool
		count int
	)

	test_mod := test_all
	test_ops := test_all

	// a == a for all a

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for _, _a := range test_ops {
			a.FromUint64(mod, _a)

			u = a.Eq(&a)

			if u == false {
				t.Errorf("a=%v\n", a);
				t.Errorf("m=%v\n", m);
				t.Fatalf("a!=a (mod %v)\n", m);
			}

			count++
		}
	}

	t.Logf("%v tests\n", count)
}

func TestSymmetry(t *testing.T) {
	var (
		a, b  Residue
		u, v  bool
		count int
	)

	test_mod := test_all
	test_ops := test_random

	for i, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			for k, _b := range test_ops {
				b.FromUint64(mod, _b)

				// (a==b) == (b==a)

				u = a.Eq(&b)
				v = b.Eq(&a)

				if u != v {
					t.Errorf("a==b = %v", u)
					t.Errorf("b==a = %v", v)
					t.Fatalf("(%v,%v,%v) count = %v\n", i, j, k, count);
				}

				count++

				// (a!=b) == (b!=a)

				u = a.Neq(&b)
				v = b.Neq(&a)

				if u != v {
					t.Errorf("a!=b = %v", u)
					t.Errorf("b!=a = %v", v)
					t.Fatalf("(%v,%v,%v) count = %v\n", i, j, k, count);
				}

				count++
			}
		}
	}

	t.Logf("%v tests\n", count)
}

func TestTransitivity(t *testing.T) {
	t.Log("Not implemented")
}

func TestAdditiveIdentity(t *testing.T) {
	var (
		zero, zero_m, zero_p, a, u, v, w, x, y, z Residue
		count                                     int
	)

	test_mod := test_all
	test_ops := test_all

	// a+0 == 0+a == a

	for _, m := range test_mod {
		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)		// modulus
		zero_p.FromUint64(mod, m)	// residue value m == 0 (mod m)
		zero.Cpy(&zero_p).Sub(&zero_p)	// 0
		zero_m.Cpy(&zero).Sub(&zero_p)	// 0-m

		for _, _a := range test_ops {
			a.FromUint64(mod, _a)

			u.Cpy(&a).Add(&zero_m)
			v.Cpy(&a).Add(&zero  )
			w.Cpy(&a).Add(&zero_p)

			x.Cpy(&zero_m).Add(&a)
			y.Cpy(&zero  ).Add(&a)
			z.Cpy(&zero_p).Add(&a)

			if a.Neq(&u) || a.Neq(&v) || a.Neq(&w) || a.Neq(&x) || a.Neq(&y) || a.Neq(&z) {
				t.Errorf("m=%v\n", m);
				t.Errorf("a=%v\n", a);
				t.Errorf("u=%v\n", u);
				t.Errorf("v=%v\n", v);
				t.Errorf("w=%v\n", w);
				t.Errorf("x=%v\n", x);
				t.Errorf("y=%v\n", y);
				t.Fatalf("z=%v\n", z);
			}

			count++
		}
	}

	t.Logf("%v tests\n", count)
}

func TestMultiplicativeIdentity(t *testing.T) {
	var (
		zero, one, one_m, one_p, a, u, v, w, x, y, z Residue
		count                                        int
	)

	test_mod := test_all
	test_ops := test_all

	// a*1 == 1*a == a

	for _, m := range test_mod {
		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)		// modulus
		zero.FromUint64(mod, m)	// residue value m == 0 (mod m)
		one.FromUint64(mod, [4]uint64{ 1, 0, 0, 0 })

		one_m.Cpy(&one).Sub(&zero)	// 1-m
		one_p.Cpy(&one).Add(&zero)	// 1+m

		for _, _a := range test_ops {
			a.FromUint64(mod, _a)

			u.Cpy(&a).Mul(&one_m)
			v.Cpy(&a).Mul(&one  )
			w.Cpy(&a).Mul(&one_p)

			x.Cpy(&one_m).Mul(&a)
			y.Cpy(&one  ).Mul(&a)
			z.Cpy(&one_p).Mul(&a)

			if a.Neq(&u) || a.Neq(&v) || a.Neq(&w) || a.Neq(&x) || a.Neq(&y) || a.Neq(&z) {
				t.Errorf("m=%v\n", m);
				t.Errorf("a=%v\n", a);
				t.Errorf("u=%v\n", u);
				t.Errorf("v=%v\n", v);
				t.Errorf("w=%v\n", w);
				t.Errorf("x=%v\n", x);
				t.Errorf("y=%v\n", y);
				t.Fatalf("z=%v\n", z);
			}

			count++
		}
	}

	t.Logf("%v tests\n", count)
}

func TestAdditiveInverse(t *testing.T) {
	var (
		a, b, u, v, w Residue
		count         int
	)

	test_mod := test_all
	test_ops := test_all

	// a+(-a) == (-a)+a == 0

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for _, _a := range test_ops {
			a.FromUint64(mod, _a)

			v.Cpy(&a).Neg()

			// a+(-a)
			u.Cpy(&a).Add(&v)

			// (-a)+a
			v.Add(&a)

			if u.Neq(&v) {
				t.Errorf("a=%v\n", a);
				t.Errorf("m=%v\n", m);
				t.Errorf("a+(-a)=%v\n", u);
				t.Fatalf("(-a)+a=%v\n", v);
			}

			count++
		}
	}

	// a == -(-a))

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for _, _a := range test_ops {
			a.FromUint64(mod, _a)

			u.Cpy(&a).Neg().Neg()

			if u.Neq(&a) {
				t.Errorf("a=%v\n", a);
				t.Errorf("m=%v\n", m);
				t.Fatalf("-(-a)=%v\n", u);
			}

			count++
		}
	}

	test_mod = test_all
	test_ops = test_random

	// a-b == (-b)+a == -(b-a)

	for i, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			for k, _b := range test_ops {
				b.FromUint64(mod, _b)

				// a-b
				u.Cpy(&a).Sub(&b)

				// (-b)+a
				v.Cpy(&b).Neg().Add(&a)

				// -(b-a)
				w.Cpy(&b).Sub(&a).Neg()

				if u.Neq(&v) || u.Neq(&w) || v.Neq(&w) {
					t.Errorf("a-b = %v", u)
					t.Errorf("-b+a = %v", v)
					t.Errorf("-(b-a) = %v", w)
					t.Fatalf("(%v,%v,%v) count = %v\n", i, j, k, count);
				}

				count++
			}
		}
	}

	t.Logf("%v tests\n", count)
}

func TestMultiplicativeInverse(t *testing.T) {

	var (
		a, u, v, one              Residue
		invertible, noninvertible int
	)

	test_mod := test_fixed
	test_ops := test_fixed

	// if 1/a exists, then
	//	a*(1/a) == (1/a)*a == 1

	for i, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)
		one.FromUint64(mod, [4]uint64{ 1, 0, 0, 0 })

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			// 1/a
			u.Cpy(&a)
			if !u.Inv() {
				noninvertible++
				continue
			}

			// a * 1/a
			v.Cpy(&a)
			v.Mul(&u)

			if v.Neq(&one) {
				t.Errorf("%v,%v\n", i, j);
				t.Fatalf("%v^2%%%v\n%v\n%v\n", a, m, u, v);
			}

			invertible++
		}
	}
	t.Logf("%v invertible, %v noninvertible\n", invertible, noninvertible)
}

func TestCommutativeAdd(t *testing.T) {
	var (
		a, b, u, v Residue
		count      int
	)

	test_mod := test_random
	test_ops := test_all

	// a+b == b+a

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			for _, _b := range test_ops[:j] {
				b.FromUint64(mod, _b)

				u.Cpy(&a)
				u.Add(&b)

				v.Cpy(&b)
				v.Add(&a)

				if u.Neq(&v) {
					t.Fatalf("%v\n+%v\n%%%v\n=%v\n=%v\n", a, b, m, u, v);
				}

				count++
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestCommutativeMul(t *testing.T) {
	var (
		a, b, u, v Residue
		count      int
	)

	test_mod := test_random
	test_ops := test_all

	// a*b == b*a

	for i, _m := range test_mod {

		if _m[3] == 0 {
			continue
		}

		m := NewModulusFromUint64(_m)

		for j, _a := range test_ops {
			a.FromUint64(m, _a)

			for k, _b := range test_ops[:j] {
				b.FromUint64(m, _b)

				u.Cpy(&a)
				u.Mul(&b)

				v.Cpy(&b)
				v.Mul(&a)

				if u.Neq(&v) {
					u.Cpy(&a)
					u.Mul(&b)

					v.Cpy(&b)
					v.Mul(&a)

					t.Errorf("ERROR: %v/%v/%v\n", i, j, k)
					t.Errorf("a   = 0x%016x%016x%016x%016x", a.r[3], a.r[2], a.r[1], a.r[0])
					t.Errorf("b   = 0x%016x%016x%016x%016x", b.r[3], b.r[2], b.r[1], b.r[0])
					t.Errorf("m   = 0x%016x%016x%016x%016x", m.m[3], m.m[2], m.m[1], m.m[0])
					t.Errorf("a*b = 0x%016x%016x%016x%016x", u.r[3], u.r[2], u.r[1], u.r[0])
					t.Fatalf("b*a = 0x%016x%016x%016x%016x", v.r[3], v.r[2], v.r[1], v.r[0])
				}

				count++
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestAssociativeAdd(t *testing.T) {
	var (
		a, b, c, u, v Residue
		count         int
	)

	test_mod := test_random
	test_ops := test_random

	// (a+b)+c == a+(b+c)

	for _, _m := range test_mod {

		if _m[3] == 0 {
			continue
		}

		m := NewModulusFromUint64(_m)

		for j, _a := range test_ops {
			a.FromUint64(m, _a)

			for k, _b := range test_ops[:j] {
				b.FromUint64(m, _b)

				for _, _c := range test_ops[:k] {
					c.FromUint64(m, _c)

					u.Cpy(&a)
					u.Add(&b)
					u.Add(&c)

					v.Cpy(&c)
					v.Add(&b)
					v.Add(&a)

					if u.Neq(&v) {
						t.Fatalf("%v+%v+%v%%%v\n%v\n%v\n", a, b, c, m, u, v);
					}

					count++
				}
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestAssociativeMul(t *testing.T) {
	var (
		a, b, c, u, v Residue
		count         int
	)

	test_mod := test_random
	test_ops := test_random

	// (ab)c == a(bc)

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			for k, _b := range test_ops[:j] {
				b.FromUint64(mod, _b)

				for _, _c := range test_ops[:k] {
					c.FromUint64(mod, _c)

					u.Cpy(&a)
					u.Mul(&b)
					u.Mul(&c)

					v.Cpy(&c)
					v.Mul(&b)
					v.Mul(&a)

					if u.Neq(&v) {
						t.Fatalf("%v*%v*%v%%%v\n%v\n%v\n", a, b, c, m, u, v);
					}

					count++
				}
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestDistributiveLeft(t *testing.T) {
	var (
		a, b, c, u, v, w Residue
		count            int
	)

	test_mod := test_random
	test_ops := test_random

	// a(b+c) == ab+ac

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			for k, _b := range test_ops[:j] {
				b.FromUint64(mod, _b)

				for _, _c := range test_ops[:k] {
					c.FromUint64(mod, _c)

					// ab+ac
					u.Cpy(&a)
					u.Mul(&b)

					v.Cpy(&a)
					v.Mul(&c)

					u.Add(&v)

					// a(b+c)
					v.Cpy(&a)
					w.Cpy(&b)
					w.Add(&c)
					v.Mul(&w)

					if u.Neq(&v) {
						t.Fatalf("%v*%v*%v%%%v\n%v\n%v\n", a, b, c, m, u, v);
					}

					count++
				}
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestDistributiveRight(t *testing.T) {
	var (
		a, b, c, u, v Residue
		count         int
	)

	test_mod := test_random
	test_ops := test_random

	// (a+b)c == ac+bc

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			for k, _b := range test_ops[:j] {
				b.FromUint64(mod, _b)

				for _, _c := range test_ops[:k] {
					c.FromUint64(mod, _c)

					// ac+bc
					u.Cpy(&a)
					u.Mul(&c)

					v.Cpy(&b)
					v.Mul(&c)

					u.Add(&v)

					// (a+b)c
					v.Cpy(&a)
					v.Add(&b)
					v.Mul(&c)

					if u.Neq(&v) {
						t.Fatalf("%v*%v*%v%%%v\n%v\n%v\n", a, b, c, m, u, v);
					}

					count++
				}
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestDouble(t *testing.T) {
	var (
		a, b, u, v Residue
		count      int
	)

	test_mod := test_all
	test_ops := test_random

	// 2a == a+a

	for _, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for _, _a := range test_ops {
			a.FromUint64(mod, _a)

			// 2a
			u.Cpy(&a).Dbl()

			// a+a
			v.Cpy(&a)
			v.Add(&a)

			if u.Neq(&v) {
				t.Fatalf("2*%v%%%v\n%v\n%v\n", a, m, u, v);
			}

			count++
		}
	}

	// 2(a+b) == 2a + 2b

	for i, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			for k, _b := range test_ops {
				b.FromUint64(mod, _b)

				// 2(a+b)
				u.Cpy(&a).Add(&b).Dbl()

				// 2a + 2b
				v.Cpy(&a).Dbl()
				b.Dbl()
				v.Add(&b)

				if u.Neq(&v) {
					b.FromUint64(mod, _b)
					t.Errorf("a: %v\nb: %v\nm: %v\nu: %v\nv: %v\n", a, b, m, u, v);

					// 2(a+b)
					u.Cpy(&a)
					t.Errorf("a: %v\n", u);
					u.Add(&b)
					t.Errorf("a+b: %v\n", u);
					u.Dbl()
					t.Errorf("2(a+b): %v\n", u);

					// 2a + 2b
					v.Cpy(&a)
					t.Errorf("a: %v\n", v);
					v.Dbl()
					t.Errorf("2a: %v\n", v);
					b.Dbl()
					t.Errorf("2b: %v\n", b);
					v.Add(&b)
					t.Errorf("2a+2b: %v\n", v);

					t.Fatalf("(%v,%v,%v)", i, j, k);
				}

				count++
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestSquare(t *testing.T) {
	var (
		a, b, u, v, w Residue
		count         int
	)

	test_mod := test_all
	test_ops := test_random

	// a^2 == a*a

	for i, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			// a^2
			u.Cpy(&a)
			u.Sqr()

			// a*a
			v.Cpy(&a)
			v.Mul(&a)

			// a*a
			w.Cpy(&a)
			w.Mul(&w)

			if u.Neq(&v) || u.Neq(&w) {
				t.Errorf("%v,%v\n", i, j);
				t.Errorf("0x%016x%016x%016x%016x^2 %% 0x%016x%016x%016x%016x\n", a.r[3], a.r[2], a.r[1], a.r[0], m[3], m[2], m[1], m[0]);
				t.Errorf("0x%016x%016x%016x%016x\n", u.r[3], u.r[2], u.r[1], u.r[0]);
				t.Errorf("0x%016x%016x%016x%016x\n", v.r[3], v.r[2], v.r[1], v.r[0]);
				t.Fatalf("0x%016x%016x%016x%016x\n", w.r[3], w.r[2], w.r[1], w.r[0]);
			}

			count++
		}
	}

	// (ab)^2 == a^2 * b^2

	for i, m := range test_mod {

		if m[3] == 0 {
			continue
		}

		mod := NewModulusFromUint64(m)

		for j, _a := range test_ops {
			a.FromUint64(mod, _a)

			for k, _b := range test_ops {
				b.FromUint64(mod, _b)

				// (a*b)^2
				u.Cpy(&a).Mul(&b).Sqr()

				// a^2 * b^2
				v.Cpy(&a).Sqr()
				b.Sqr()
				v.Mul(&b)

				if u.Neq(&v) {
					b.FromUint64(mod, _b)

					// (a*b)^2
					u.Cpy(&a)
					t.Errorf("a = %v\n", u);
					u.Mul(&b)
					t.Errorf("ab = %v\n", u);
					u.Sqr()
					t.Errorf("(ab)^2 = %v\n", u);

					t.Errorf("%%%v = %v\n", m, u);

					// a^2 * b^2
					v.Cpy(&a)
					t.Errorf("a = %v\n", v);
					v.Sqr()
					t.Errorf("a^2 = %v\n", v);
					t.Errorf("b = %v\n", b);
					b.Sqr()
					t.Errorf("b^2 = %v\n", b);
					v.Mul(&b)
					t.Errorf("a^2*b^2 = %v\n", v);

					t.Errorf("%%%v = %v\n", m, v);

					t.Fatalf("(%v,%v,%v) count = %v\n", i, j, k, count);
				}

				count++
			}
		}
	}
	t.Logf("%v tests\n", count)
}

func TestExponentiation(t *testing.T) {
	var (
		a, b, c Residue
		eb      ExpBase
		count   int
	)

	test_ops := test_fixed

	// a^m == a

	mod := NewModulusFromUint64(nistp256)

	for i, _a := range test_ops {
		a.FromUint64(mod, _a)

		invertible := b.Cpy(&a).Inv()

		if !invertible {
			// nistp256 is prime, so only 0 lacks an inverse,
			// and 0 is the only fixed value = 0 (mod nistp256)
			if _a[3] | _a[2] | _a[1] | _a[0] != 0 {
				t.Fatalf("0X%016X%016X%016X%016X\n",
					a.r[3], a.r[2], a.r[1], a.r[0])
			}
			continue
		}

		// b = a^m
		eb.FromResidue(&a)
		b.ExpPrecomp(&eb, nistp256)

		// c = a^m
		c.Cpy(&a).Exp(nistp256)

		if a.Neq(&b) || a.Neq(&c) {
			t.Errorf("%v", i)
			t.Errorf("0X%016X%016X%016X%016X\n",
				a.r[3], a.r[2], a.r[1], a.r[0])
			t.Errorf("0X%016X%016X%016X%016X\n",
				nistp256[3], nistp256[2], nistp256[1], nistp256[0])
			t.Errorf("0X%016X%016X%016X%016X\n",
				b.r[3], b.r[2], b.r[1], b.r[0])
			t.Fatalf("0X%016X%016X%016X%016X\n",
				c.r[3], c.r[2], c.r[1], c.r[0])
		}

		count++
	}
	t.Logf("%v tests\n", count)
}

var (
	nistp256 [4]uint64
	nistp224 [4]uint64
	x, y     Residue
)

func BenchmarkMod256(b *testing.B) {
	b.Run("Neg", benchmarkNeg)
	b.Run("Dbl", benchmarkDbl)
	b.Run("Sub", benchmarkSub)
	b.Run("Add", benchmarkAdd)

	b.Run("Sqr", benchmarkSqr)
	b.Run("Mul", benchmarkMul)
	b.Run("Inv", benchmarkInv)
	b.Run("Exp", benchmarkExp)
	b.Run("ExpPrecomp", benchmarkExpPrecomp)
}

func benchmarkNeg(b *testing.B) {
	m := NewModulusFromUint64(nistp256)

	x.FromUint64(m, [4]uint64{257, 479, 487, 491})
	y.FromUint64(m, [4]uint64{997, 499, 503, 509})

	//b.ResetTimer()

	for i := 0; i < b.N; i+=2 {
		x.Neg()
		y.Neg()
	}
}

func benchmarkDbl(b *testing.B) {
	m := NewModulusFromUint64(nistp256)

	x.FromUint64(m, [4]uint64{257, 479, 487, 491})
	y.FromUint64(m, [4]uint64{997, 499, 503, 509})

	//b.ResetTimer()

	for i := 0; i < b.N; i+=2 {
		x.Dbl()
		y.Dbl()
	}
}

func benchmarkSub(b *testing.B) {
	m := NewModulusFromUint64(nistp256)

	x.FromUint64(m, [4]uint64{257, 479, 487, 491})
	y.FromUint64(m, [4]uint64{997, 499, 503, 509})

	//b.ResetTimer()

	for i := 0; i < b.N; i+=2 {
		x.Sub(&y)
		y.Sub(&x)
	}
}

func benchmarkAdd(b *testing.B) {
	m := NewModulusFromUint64(nistp256)

	x.FromUint64(m, [4]uint64{257, 479, 487, 491})
	y.FromUint64(m, [4]uint64{997, 499, 503, 509})

	//b.ResetTimer()

	for i := 0; i < b.N; i+=2 {
		x.Add(&y)
		y.Add(&x)
	}
}

func benchmarkSqr(b *testing.B) {
	m := NewModulusFromUint64(nistp256)

	x.FromUint64(m, [4]uint64{257, 479, 487, 491})
	y.FromUint64(m, [4]uint64{997, 499, 503, 509})

	//b.ResetTimer()

	for i := 0; i < b.N; i+=2 {
		x.Sqr()
		y.Sqr()
	}
}

func benchmarkMul(b *testing.B) {
	m := NewModulusFromUint64(nistp256)

	x.FromUint64(m, [4]uint64{257, 479, 487, 491})
	y.FromUint64(m, [4]uint64{997, 499, 503, 509})

	//b.ResetTimer()

	for i := 0; i < b.N; i+=2 {
		x.Mul(&y)
		y.Mul(&x)
	}
}

func benchmarkInv(b *testing.B) {
	var (
		a     Residue
		count int
	)

	test_ops := test_all

	mod := NewModulusFromUint64(nistp256)

OuterLoop:
	for {
		for _, _a := range test_ops {
			a.FromUint64(mod, _a)

			a.Inv()
			a.Inv()

			count += 2

			if count >= b.N {
				break OuterLoop
			}
		}
	}
}

func benchmarkExp(b *testing.B) {
	var (
		a     Residue
		count int
	)

	test_mod := test_all
	test_ops := test_all

	// a ^ a % m

OuterLoop:
	for {
		for _, m := range test_mod {

			if m[3] == 0 {
				continue
			}

			mod := NewModulusFromUint64(m)

			for _, _a := range test_ops {
				a.FromUint64(mod, _a)

				// a = a^a
				a.Exp(_a)
				a.Exp(_a)

				count += 2

				if count >= b.N {
					break OuterLoop
				}
			}
		}
	}
}

func benchmarkExpPrecomp(b *testing.B) {
	var (
		a, u  Residue
		eb    ExpBase
		count int
	)

	test_mod := test_all
	test_ops := test_all

	// a ^ e % m

OuterLoop:
	for {
		for _, m := range test_mod {

			if m[3] == 0 {
				continue
			}

			mod := NewModulusFromUint64(m)

			for _, _a := range test_ops {
				a.FromUint64(mod, _a)

				eb.FromResidue(&a)
				for _, e := range test_ops {

					u.ExpPrecomp(&eb, e)
					u.ExpPrecomp(&eb, e)

					count += 2

					if count >= b.N {
						break OuterLoop
					}
				}
			}
		}
	}
}
