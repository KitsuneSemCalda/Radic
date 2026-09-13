# The `number` type and its mold

Radic aims to be stupidly simple and compile 1:1 to assembly, without implicit
optimizations. This document defines how the `number` type works: it is a single
type whose bit width (the *mold*) is inferred by the compiler, with an escape
hatch to pin an exact width when it matters.

## Philosophy

Writing down an explicit width for every integer (`i8`, `u32`, ...) is noise for
a language whose whole point is simplicity. Because codegen is 1:1 to assembly,
the compiler always has full knowledge of every value at compile time, so it can
decide the mold itself and stamp it onto the value's characteristics.

This is deliberately "babysitting": the compiler handles widths so the
programmer doesn't have to. To keep that honest, the language offers no implicit
optimization at runtime — the mold is decided entirely at compile time and never
changes during execution. When a programmer actually needs exact control (ABI,
interop, struct layout), a pinned width is available.

## The default: `number`

`number` is a single type, **signed by default**. Unsigned only emerges from
context. The compiler infers a mold (width + signedness) for each value.

- **Registers and operations** always work on the target's native word size
  (64-bit on x86-64). No operand-size prefixes, no partial-register tricks.
- **Memory and variables** are allocated at the smallest width that accommodates
  the value (1, 2, 4, or 8 bytes).
- The mold never truncates and never errors at runtime; it only lives in the
  compiler's analysis.

### Literal inference

The mold of a numeric literal is chosen by its value:

| value                    | mold                | example    |
| ------------------------ | ------------------- | ---------- |
| negative                 | smallest signed     | `-5` → i8  |
| fits signed at width     | signed at width     | `5` → i8   |
| doesn't fit signed, fits unsigned at width | unsigned at width | `200` → u8  |
| otherwise                | smallest signed next width | `300` → i16 |

Value ladder: i8/u8 → i16/u16 → i32/u32 → i64/u64. Values beyond u64 are a
compile-time error.

### Growth by fixed point

The mold of a variable is the **fixed point of all requirements across every
code path**. Each assignment and operation records the minimum width it needs;
the final mold is the largest requirement seen anywhere.

Examples:

```radic
x = number 5         // mold 8-bit
x = x + 300          // grows to 16-bit (i16)
```

Never a runtime error, never silent truncation — growth happens before codegen.

### Struct fields

A struct field has a fixed layout, so the mold of a `number` field is the
**maximum across all uses of the field in the whole program**. A field used
somewhere with a large literal becomes that wide. When that implicit layout is
unacceptable, pin the field (below).

## The escape hatch: pinned widths

When the exact width matters (ABI, interop, memory-mapped I/O, packed
protocols), the following are recognized as builtin type names just like
`number`:

```
i8  i16  i32  i64
u8  u16  u32  u64
```

They can be used at any type site: struct/union fields, parameters, returns,
variable declarations.

```radic
field: u32          // exactly 4 bytes, no inference, no growth
func parse() -> u16 // fixed return width
```

A pinned width is **fixed**: it never participates in inference and never grows.
Assigning a value that does not fit the pinned width is a **compile-time error**
— never truncation, never growth. This is the guarantee that gives real
low-level control.

The unsigned flip described above only applies to inferred (`number`) sites;
pinned widths respect exactly what was written.

## Codegen mapping

| width | load (signed)        | load (unsigned)     | store            |
| ----- | -------------------- | ------------------- | ---------------- |
| 1     | `movsx reg, byte [m]` | `movzx reg, byte [m]` | `mov byte [m], al` |
| 2     | `movsx reg, word [m]` | `movzx reg, word [m]` | `mov word [m], ax`  |
| 4     | `movsxd reg, dword [m]` | `mov reg, dword [m]`  | `mov dword [m], eax` |
| 8     | `mov reg, qword [m]`  | `mov reg, qword [m]`  | `mov qword [m], rax` |

Arithmetic and comparisons operate on registers at native word size, so they emit
plain 64-bit instructions regardless of the mold. The mold only affects memory
accesses — which is precisely what keeps codegen boring and 1:1.

## Constraints

- No runtime cost: the mold is decided during analysis, before codegen.
- No implicit truncation or silent size changes after analysis.
- Pinned widths take precedence over inference at their site.
- A `number`-typed value can be widened by fixed-point growth, but never
  narrowed by assignment without an explicit pinned width.

## Roadmap notes

- The AST keeps types lexical as `Name{ Lexeme, Builtin }`; the parser only
  produces `number` or a pinned-width name.
- The semantic analysis pass uses the future `internal/types` package to
  resolve molds: `MoldOfLiteral` for inference and pinned-width parsing for the
  escape hatch.
- This document is the source of truth for that pass; see also
  `internal/token` for the existing type keywords.