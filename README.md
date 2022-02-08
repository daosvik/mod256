# Library for arithmetic modulo any 193- to 256-bit modulus

This library only supports modular arithmetic.

- No boolean operations or other bitwise operations like shifts are supported.
- Speed is gained by postponing full reduction to canonical (least non-negative) residues.
- Internally any 256-bit representative of each residue class is used.
- Reduction to canonical residue is only performed when converting for external use or testing for equality of residues.
- Residues are treated as being different when their moduli are different. E.g. 2 mod 3 is not the same as 2 mod 4.
- Arrays of uint64 are treated as little-endian. Hence the array [4]uint64{ 1, 0, 0, 0 } contains the value 1.

The library is alloc-free, and code coverage is at 99.9%.

## Security

This library is **not** meant to protect sensitive data like cryptographic keys.

Although some operations should be constant-time on most architectures, the library does **not** protect from e.g. timing or cache attacks.

## Types

### Modulus

Contains a modulus `m` as well as derived values that help speed up computations.
The allowed range for `m` is `2^192` to `2^256-1`.

### Residue

Contains a representative of a residue class, and the pointer to its modulus.

### ExpBase

Contains lookup tables allowing fast repeated modular exponentiation with the same base value.

## Operations

- **Modulus.FromUint64()**
- **Modulus.ToUint64()** returns an array with the modulus.

- **Residue.FromUint64()**
- **Residue.ToUint64()** returns an array with the canonical representative of the residue class.

- **Residue.Eq()** compares one residue to another, returns true when equal.
- **Residue.Neq()** compares one residue to another, returns true when different.

- **Residue.Cpy()** copies one residue to another. Both the residue value and the modulus pointer are copied.
- **Residue.Neg()** computes the negation (additive inverse) of a residue.
- **Residue.Add()** computes the sum of two residues.
- **Residue.Sub()** computes the sum of a residue and the negation of a second residue.
- **Residue.Dbl()** computes the double of a residue.
- **Residue.Mul()** computes the product of two residues.
- **Residue.Sqr()** computes the square of a residue.
- **Residue.Inv()** computes the (multiplicative) inverse of a residue, if it exists.
- **Residue.Exp()** performs modular exponentiation without storing precomputed values for later use.
- **Residue.ExpPrecomp()** takes an ExpBase computed from the base value, a 256-bit integer as the exponent, and performs modular exponentiation.

- **ExpBase.FromResidue()** computes and stores lookup tables for fast exponentiation of the given residue.

## Benchmarks

The speeds achieved in mod256 compared to [uint256](https://github.com/holiman/uint256) below illustrate the performance advantage of designing specifically for modular arithmetic.

Note that the aim of uint256 is to provide a replacement for big.Int for 256-bit integers. Therefore there is only partially overlapping functionality.

All benchmarks were performed using Go version 1.17.6. For brevity, only each median of 9 tests is shown here.

### Ice Lake (mobile) @ 1.5 GHz

```
BenchmarkMod256/Neg-8         	187632590	         6.402 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Dbl-8         	87988447	        13.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sub-8         	78145950	        15.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Add-8         	69903696	        16.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sqr-8         	10063964	       119.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Mul-8         	 9191982	       130.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Inv-8         	  113192	     10662 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Exp-8         	   31624	     37431 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/ExpPrecomp-8  	   72931	     16314 ns/op	       0 B/op	       0 allocs/op

BenchmarkAddMod/mod256/uint256-8      	29427345	        40.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256-8      	 3651633	       328.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256r-8     	 8301348	       144.8 ns/op	       0 B/op	       0 allocs/op
```

### M1 @ 3.2 GHz

```
BenchmarkMod256/Neg-8         	225061627	         5.398 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Dbl-8         	208388655	         5.758 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sub-8         	166325428	         7.203 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Add-8         	165019262	         7.273 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sqr-8         	39676798	        30.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Mul-8         	33286484	        36.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Inv-8         	  433461	      2785 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Exp-8         	  129596	      9896 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/ExpPrecomp-8  	  269511	      4452 ns/op	       0 B/op	       0 allocs/op

BenchmarkAddMod/mod256/uint256-8      	76248772	        15.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256-8      	13249834	        90.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256r-8     	29764826	        40.28 ns/op	       0 B/op	       0 allocs/op
```

### Skylake @ 2.2 GHz

```
BenchmarkMod256/Neg-4         	260609025	         4.576 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Dbl-4         	100000000	        11.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sub-4         	95406411	        12.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Add-4         	89778802	        13.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sqr-4         	14057278	        84.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Mul-4         	13072905	        91.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Inv-4         	  143182	      8390 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Exp-4         	   45506	     26153 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/ExpPrecomp-4  	  105787	     11336 ns/op	       0 B/op	       0 allocs/op

BenchmarkAddMod/mod256/uint256-4      	37145136	        32.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256-4      	 5047970	       237.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256r-4     	12708961	        94.22 ns/op	       0 B/op	       0 allocs/op
```

### Zen 2 @ 3.6 GHz

```
BenchmarkMod256/Neg-12         	333289629	         3.844 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Dbl-12         	150186574	         8.503 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sub-12         	124515450	        10.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Add-12         	100000000	        10.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sqr-12         	20554772	        53.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Mul-12         	23398146	        57.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Inv-12         	  152364	      7486 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Exp-12         	   64242	     17318 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/ExpPrecomp-12  	  150073	      7408 ns/op	       0 B/op	       0 allocs/op

BenchmarkAddMod/mod256/uint256-12      	53721606	        22.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256-12      	 9285478	       128.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256r-12     	23170990	        51.72 ns/op	       0 B/op	       0 allocs/op
```

### Zen 3 @ 3.9 GHz

```
BenchmarkMod256/Neg-12         	457737588	         2.598 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Dbl-12         	210516388	         5.698 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sub-12         	211873234	         5.660 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Add-12         	197642300	         6.095 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Sqr-12         	27790299	        42.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Mul-12         	26322823	        45.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Inv-12         	  269335	      4465 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/Exp-12         	   90829	     13183 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod256/ExpPrecomp-12  	  210309	      5669 ns/op	       0 B/op	       0 allocs/op

BenchmarkAddMod/mod256/uint256-12      	60854953	        19.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256-12      	10438305	       114.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulMod/mod256/uint256r-12     	22693557	        52.73 ns/op	       0 B/op	       0 allocs/op
```

## Testing

Benchmarks cover most properties of [commutative rings](https://en.wikipedia.org/wiki/Commutative_ring).
The one exception is transitivity of equality, as no meaningful test of this was found.

Structured as well as random test values are used.
