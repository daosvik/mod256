# Library for arithmetic modulo any 193-bit to 256-bit modulus

This library only supports [modular arithmetic](https://en.wikipedia.org/wiki/Modular_arithmetic).

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

## Testing

Tests cover most properties of [commutative rings](https://en.wikipedia.org/wiki/Commutative_ring).
The one exception is transitivity of equality, as no meaningful test for this was found.

Structured as well as random test values are used.

## Benchmarks

The speeds achieved in mod256 compared to [uint256](https://github.com/holiman/uint256) below illustrate the performance
advantage of designing specifically for modular arithmetic.

Note that the aim of uint256 is to provide a replacement for [big.Int](https://pkg.go.dev/math/big) for 256-bit integers.
Therefore there is only partially overlapping functionality, which is for modular addition and modular multiplication.

```
                       Ice Lake 1.5GHz         M1 3.2GHz    Skylake 2.2GHz      Zen 2 3.6GHz      Zen 3 3.9GHz
Mod256/Neg                 3.94ns ± 0%       3.60ns ± 0%       3.35ns ± 0%       2.16ns ± 0%       1.83ns ± 0%
Mod256/Dbl                 13.6ns ± 0%       5.76ns ± 0%       11.2ns ± 0%       5.86ns ± 0%       5.73ns ± 1%
Mod256/Sub                 14.1ns ± 0%       6.49ns ± 0%       10.5ns ± 0%       6.73ns ± 0%       5.76ns ± 0%
Mod256/Sqr                  104ns ± 0%       24.4ns ± 0%       70.8ns ± 0%       37.3ns ± 0%       36.2ns ± 0%
Mod256/Inv                 10.6µs ± 0%       2.79µs ± 0%       8.38µs ± 0%       5.34µs ± 0%       4.50µs ± 0%
Mod256/Exp                 34.8µs ± 0%       8.14µs ± 0%       24.4µs ± 0%       12.8µs ± 0%       12.3µs ± 0%
Mod256/ExpPrecomp          10.8µs ± 0%       2.57µs ± 0%       7.28µs ± 0%       3.91µs ± 0%       3.76µs ± 0%

Mod256/Add                 13.9ns ± 1%       6.64ns ± 0%       10.7ns ± 0%       6.73ns ± 0%       5.85ns ± 3%
Mod256/Mul                  112ns ± 0%       28.0ns ± 0%       77.5ns ± 0%       39.8ns ± 0%       37.8ns ± 0%

AddMod/mod256/uint256      40.4ns ± 0%       15.6ns ± 2%       32.3ns ± 0%       22.2ns ± 0%       19.6ns ± 0%
MulMod/mod256/uint256       328ns ± 0%       90.2ns ± 0%        237ns ± 0%        129ns ± 0%        114ns ± 0%
MulMod/mod256/uint256r      145ns ± 0%       40.3ns ± 0%       94.3ns ± 0%       51.8ns ± 0%       52.8ns ± 0%

Speedup factors:
Mod256/Add vs. AddMod             2.91              2.35              3.02              3.30              3.35
Mod256/Mul vs. MulMod             2.93              3.22              3.06              3.24              3.02
Mod256/Mul vs. MulMod(r)          1.29              1.44              1.22              1.30              1.40
```

All benchmarks were performed using Go version 1.17.6, and all report 0 B/op and 0 allocs/op, both for mod256 and uint256.

Summary:
- Modular addition is from 2.35 (M1) to 3.35 (Zen 3) times as fast as uint256.
- Modular multiplication is from 2.93 (Ice Lake) to 3.24 (Zen 2) times as fast as multiplication *without* reciprocal cache in uint256.
- Modular multiplication is from 22% (Skylake) to 44% (M1) faster than multiplication *with* reciprocal cache in uint256.
