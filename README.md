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
Mod256/Neg                 6.35ns ± 0%       5.39ns ± 1%       4.60ns ± 0%       2.99ns ± 0%       2.59ns ± 0%
Mod256/Dbl                 13.6ns ± 0%       5.76ns ± 0%       11.2ns ± 0%       5.87ns ± 0%       5.70ns ± 0%
Mod256/Sub                 15.6ns ± 1%       7.20ns ± 0%       12.5ns ± 0%       7.62ns ± 0%       5.80ns ± 0%
Mod256/Sqr                  119ns ± 1%       30.2ns ± 0%       84.8ns ± 0%       42.8ns ± 0%       41.9ns ± 0%
Mod256/Inv                 10.6µs ± 0%       2.78µs ± 0%       8.40µs ± 0%       5.35µs ± 0%       4.48µs ± 0%
Mod256/Exp                 37.0µs ± 0%       9.91µs ± 0%       26.0µs ± 0%       13.8µs ± 0%       13.1µs ± 0%
Mod256/ExpPrecomp          12.5µs ± 0%       3.37µs ± 0%       8.58µs ± 0%       4.61µs ± 0%       4.33µs ± 0%

Mod256/Add                 16.1ns ± 2%       7.27ns ± 0%       13.4ns ± 0%       8.58ns ± 0%       6.14ns ± 3%
Mod256/Mul                  131ns ± 0%       36.0ns ± 0%       91.8ns ± 0%       47.4ns ± 0%       45.9ns ± 0%

AddMod/mod256/uint256      40.4ns ± 0%       15.6ns ± 2%       32.3ns ± 0%       22.2ns ± 0%       19.6ns ± 0%
MulMod/mod256/uint256       328ns ± 0%       90.2ns ± 0%        237ns ± 0%        129ns ± 0%        114ns ± 0%
MulMod/mod256/uint256r      145ns ± 0%       40.3ns ± 0%       94.3ns ± 0%       51.8ns ± 0%       52.8ns ± 0%

Speedup factors:
Mod256/Add vs. AddMod             2.51              2.15              2.41              2.59              3.19
Mod256/Mul vs. MulMod             2.50              2.51              2.58              2.72              2.48
Mod256/Mul vs. MulMod(r)          1.11              1.12              1.03              1.09              1.15
```

All benchmarks were performed using Go version 1.17.6, and all report 0 B/op and 0 allocs/op, both for mod256 and uint256.

Summary:
- Modular addition is from 2.15 (M1) to 3.19 (Zen 3) times as fast as uint256.
- Modular multiplication is from 2.48 (Zen 3) to 2.72 (Zen 2) times as fast as multiplication *without* reciprocal cache in uint256.
- Modular multiplication is from 2.7% (Skylake) to 15% (Zen 3) faster than multiplication *with* reciprocal cache in uint256.
