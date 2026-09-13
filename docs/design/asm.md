# The low level: C, Radic, and assembly

This document looks at what C and Radic actually mean at the lowest level:
x86-64 assembly in Intel syntax, targeting the System V ABI (the convention
used on Linux). Each construct is shown three ways — C source, Radic source,
and the assembly both compile to — so the mapping from both languages to the
machine is explicit.

The C column is factual: it was produced by `gcc -O0 -masm=intel`. Radic has
no code generator yet, so its column is the **planned** codegen, derived from
the rules in `docs/design/number.md` and the language's 1:1 promise.
Instructions and offsets are representative, not a promise about exact frame
layout.

## The core difference: who decides the machine code

|                       | C                                            | Radic                              |
| --------------------- | -------------------------------------------- | ---------------------------------- |
| contract              | "as if" semantics; compiler is free to pick  | 1:1 to assembly, no implicit optimization |
| at `-O0`              | close to 1:1, honest floor                   | always, by definition               |
| at `-O2`              | inconsistent, everything may change          | nothing can change                  |
| integer width         | fixed per type (`int` = 4 bytes)             | `number` with inferred mold         |
| register operations   | whatever the compiler chose                  | always the native word size (64-bit) |
| memory allocation     | compiler-chosen slots                        | smallest width that fits the value  |

GCC at `-O0` looks a lot like what Radic promises, which is why this document
uses it as the C reference. The two diverge in one place that matters: `int`
is always 4 bytes, and GCC may still reorder or fold even at `-O0`. Radic
fixes both: the mold decides memory width per value, and operations always run
at 64 bits.

## Literals, variables, and the mold

```c
int x = 5;
unsigned short y = 300;
return x + y;
```

```radic
number x = 5;       // mold: i8
number y = 300;     // mold: i16 (doesn't fit u8 either)
return x + y;
```

C (`gcc -O0`, trimmed — `[rbp-4]` is just GCC's chosen stack slot):

```asm
mov dword [rbp-4], 5        ; int is always 4 bytes
mov word  [rbp-6], 300
movzx edx, word  [rbp-6]    ; load y
mov   eax, dword [rbp-4]    ; load x
add   eax, edx              ; 32-bit operation
```

Radic (planned — see the load/store table in `docs/design/number.md`):

```asm
mov byte [rbp-2], 5         ; mold i8  -> 1 byte
mov word [rbp-4], 300       ; mold i16 -> 2 bytes
movsx rax, word [rbp-4]     ; load y, sign-extend to 64 bits
movsx rcx, byte [rbp-2]     ; load x, sign-extend to 64 bits
add   rax, rcx              ; 64-bit operation, always
```

The mold only changes the memory access. Once a value is in a register it is
always the native 64-bit word, which is what keeps the codegen boring and
predictable.

## Arithmetic and bitwise

Both languages use the same mnemonics; only the register width differs (C's
`int`: 32-bit; Radic: always 64-bit).

| operation | C (`int`)       | Radic (`number`) |
| --------- | --------------- | ---------------- |
| `a + b`   | `add eax, edx`  | `add rax, rcx`   |
| `a - b`   | `sub eax, edx`  | `sub rax, rcx`   |
| `a * b`   | `imul eax, edx` | `imul rax, rcx`  |
| `a / b`   | `cdq; idiv esi` | `cqo; idiv rcx`  |
| `a % b`   | `cdq; idiv esi` (rem in `edx`) | `cqo; idiv rcx` (rem in `rdx`) |
| unary `-` | `neg eax`       | `neg rax`        |
| `& \| ^`  | `and/or/xor`    | `and/or/xor`     |
| `~`       | `not`           | `not`            |
| `<<`      | `shl`           | `shl`            |
| `>>`      | `sar` (signed)  | `sar`            |
| `x++ x--` | `inc/dec`       | `inc/dec`        |

Division is the interesting one: `idiv` always divides the double-width
`rdx:rax` by the operand. Signed division therefore needs the sign-extension
first — `cdq` (`eax` → `edx:eax`) in C, `cqo` (`rax` → `rdx:rax`) in Radic.
The quotient lands in `eax/rax`, the remainder in `edx/rdx`.

Real GCC output for `int q = a / b; int r = a % b;` (`-O0`):

```asm
mov eax, DWORD PTR [rbp-36]  ; a
cdq
idiv DWORD PTR [rbp-40]      ; b
mov DWORD PTR [rbp-8], eax   ; q
mov eax, DWORD PTR [rbp-36]
cdq
idiv DWORD PTR [rbp-40]
mov DWORD PTR [rbp-4], edx   ; r
```

## Comparisons and branches

Conditionals and loops lower to the same shape in both languages:

- `cmp a, b` sets the flags; a `jcc` branch reads them
  (`je`/`jne`, `jl`/`jle`, `jg`/`jge`; unsigned variants `jb`/`ja`).
- `if`/`else` becomes a `jcc` to the else branch plus a `jmp` over it.
- An inverted test like `if (!x)` becomes `test x, x` + `je`.
- `while`/`for` become a condition test at the top, a body label, and a
  `jmp` back to the test.

Real GCC output for `if (a > b) { a--; } else { a++; }`:

```asm
mov eax, DWORD PTR [rbp-4]   ; a
cmp eax, DWORD PTR [rbp-8]   ; a > b ?
jle .L2
sub DWORD PTR [rbp-4], 1     ; then: a--
jmp .L3
.L2:
add DWORD PTR [rbp-4], 1     ; else: a++
.L3:
```

Radic emits the identical shape; the loads simply use `movsx` and the tests
run on 64-bit registers. A `for (i = 0; i < 10; i++)` loop compiles as:

```asm
mov  byte [rbp-8], 0         ; i = 0   (mold i8)
jmp  .Ltop
.Lbody:
movsx rax, byte [rbp-8]      ; load i
inc   rax                    ; i++
mov   byte [rbp-8], al       ; store i
.Ltop:
movsx rax, byte [rbp-8]      ; load i
cmp   rax, 10                ; i < 10 ?
jl    .Lbody
```

## Logical operators

C's `&&` and `||` short-circuit: the right side only runs if the left side
decides it matters. The codegen is a chain of tests with jumps to a shared
outcome label. Real GCC output for `if (a > 0 && b > 0)`:

```asm
cmp DWORD PTR [rbp-20], 0    ; a > 0 ?
jle .L2                      ; no -> whole thing is false
cmp DWORD PTR [rbp-24], 0    ; b > 0 ?
jle .L2                      ; no -> false
mov DWORD PTR [rbp-4], 1     ; yes: both true
.L2:
```

Radic's `and`/`or` words map to these logical operators (the parser keeps
them distinct from the bitwise `&`/`|`), but **whether they short-circuit is
still an open question**. C-like short-circuiting costs branches and code
duplication, which is the opposite of "boring"; eager evaluation is simpler
and 1:1, but changes observable behavior (a call on the right side would run
even when the left is false). Until the semantic pass exists, the docs leave
this undecided; the assembly above assumes C-like behavior. The unary `not`
compiles to `test rax, rax` + `sete` (or a `cmp` + `setcc`).

## switch

C compiles `switch` honestly at `-O0` — a straight chain of `cmp`/`je`:

```asm
cmp DWORD PTR [rbp-4], 1
je  .L2
cmp DWORD PTR [rbp-4], 2
je  .L3
jmp .L6                      ; default
.L2:
mov eax, 10
jmp .L5
.L3:
mov eax, 20
jmp .L5
.L6:
mov eax, 0
.L5:
```

At `-O2` GCC replaces this with a **jump table** — an indirect `jmp` through
an array of addresses indexed by the value. That is exactly the kind of
implicit optimization Radic refuses: it would change code size, memory
layout, and branch behavior behind the programmer's back. Radic's
`case 1 { ... }` blocks (implicit break, no fallthrough) are a natural fit
for the plain `cmp`/`je` chain, and that chain is what Radic will always
emit. No jump tables, ever.

## Struct and union

A struct field is fixed memory at a known offset from the struct's base.
C writes and reads it directly from the slot it chose:

```asm
mov DWORD PTR [rbp-8], 3     ; p.x = 3
mov DWORD PTR [rbp-4], 4     ; p.y = 4
```

In Radic each `number` field gets the mold = maximum across all uses in the
program (see `docs/design/number.md`), so the same writes at i8 are a byte
apart:

```asm
mov byte [rbp-2], 3          ; p.x = 3
mov byte [rbp-1], 4          ; p.y = 4
```

Two things are still to be specified for Radic: aggregate layout at the mold
level (padding, if any, between fields) and exactly how a `union`'s shared
storage picks its mold. The parser accepts both `struct` and `union`
declarations today; the semantic pass will decide.

## Pointers

The three pointer operators all appear in both languages:

| operation | semantics | assembly |
| --------- | --------- | -------- |
| `&v`      | take address | `lea rax, [rbp-offset]` |
| `*p`      | dereference  | `mov rax, [rax]` / `mov [rax], rcx` |
| `p->mt`   | add offset + dereference | `mov eax, [rax + 4]` |

Real GCC output for `int *p = &v; *p = 6;` (with the stack-canary noise
removed):

```asm
mov  DWORD PTR [rbp-36], 5   ; v = 5
lea  rax, [rbp-36]           ; &v
mov  QWORD PTR [rbp-32], rax ; p = &v
mov  rax, QWORD PTR [rbp-32] ; load p
mov  DWORD PTR [rax], 6      ; *p = 6
```

Pointers in Radic are only half-born: the lexer and parser already handle the
unary `*`, `&` and the `->` member operator (see `docs/design/c_radic.md`),
but there is no pointer type syntax yet, so you cannot declare a pointer
variable. That part is pending. When it lands, the lowering is exactly the
standard `lea` / memory-operand pattern above — pointer arithmetic, arrays,
and subscripting are not planned as implicit conveniences.

## Functions and the ABI

Both languages target the System V ABI on x86-64:

- integer/pointer arguments in `rdi, rsi, rdx, rcx, r8, r9`, then the stack;
- the return value in `rax` (or its 32-bit half `eax`);
- the callee owns its own stack frame: `push rbp; mov rbp, rsp` prologue,
  optionally `sub rsp, N`, and `leave; ret` at the end.

Real GCC output for `int total = sum(3, 4);` in `main`:

```asm
mov esi, 4                   ; arg 2
mov edi, 3                   ; arg 1
call sum
mov DWORD PTR [rbp-4], eax   ; total = eax
```

```asm
sum:
push rbp
mov  rbp, rsp
mov  DWORD PTR [rbp-4], edi
mov  DWORD PTR [rbp-8], esi
mov  edx, DWORD PTR [rbp-4]
mov  eax, DWORD PTR [rbp-8]
add  eax, edx
pop  rbp
ret
```

Radic's planned codegen follows the same ABI. Arguments are widened to the
full 32-bit register as the ABI demands, and a `-> number` return leaves its
value in `rax`. Radic adds one quirk of its own: parameter and return molds
are inferred from how the values are actually used, but the *convention*
never changes — registers, `call`, `ret`, `leave`. Nothing about a function
call can surprise you.

## enum

Enums are just constants. Both languages lower a variant to an immediate
store followed by whatever the expression needs:

```asm
mov DWORD PTR [rbp-4], 2     ; enum Color c = Green;  (Green == Red+1)
```

Radic's enum variants can carry explicit values too; a constant like that
becomes a `mov imm` at the mold width chosen for it.

## One program, three views

The same program from `docs/design/c_radic.md`, now with the machine code
both compilers would produce.

```c
struct Point { int x; int y; };

int sum(int a, int b) {
    return a + b;
}

int main(void) {
    struct Point p;
    p.x = 3;
    p.y = 4;

    if ((p.x > 0) && (p.y > 0)) {
        int total = sum(p.x, p.y);
        return total;
    }
    return 0;
}
```

```radic
struct Point {
    number x;
    number y;
}

func sum(number a, number b) -> number {
    return a + b;
}

func main() {
    p = Point;            // local allocation; fields assigned below
    p.x = 3;
    p.y = 4;

    if (p.x > 0 and p.y > 0) {
        number total = sum(p.x, p.y);
        return total;
    }
    return 0;
}
```

C as compiled by GCC at `-O0` (`[rbp-X]` slots are GCC's choice):

```asm
main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16
    mov  dword [rbp-8], 3       ; p.x = 3
    mov  dword [rbp-4], 4       ; p.y = 4
    mov  eax, dword [rbp-8]     ; p.x > 0?
    test eax, eax
    jle  .L4
    mov  eax, dword [rbp-4]     ; p.y > 0?
    test eax, eax
    jle  .L4
    mov  edx, dword [rbp-4]     ; arg 2 = p.y
    mov  eax, dword [rbp-8]     ; arg 1 = p.x
    mov  esi, edx
    mov  edi, eax
    call sum
    mov  dword [rbp-12], eax    ; total = sum(...)
    mov  eax, dword [rbp-12]
    jmp  .L6
.L4:
    mov  eax, 0                 ; return 0
.L6:
    leave
    ret
```

The same program as Radic would emit it (planned; assumes C-like
short-circuit for `and`):

```asm
main:
    push rbp
    mov  rbp, rsp
    sub  rsp, 16
    mov  byte [rbp-2], 3        ; p.x = 3   (mold i8)
    mov  byte [rbp-1], 4        ; p.y = 4   (mold i8)
    movsx rax, byte [rbp-2]     ; p.x > 0?  (64-bit test)
    cmp  rax, 0
    jle  .Lelse
    movsx rax, byte [rbp-1]     ; p.y > 0?
    cmp  rax, 0
    jle  .Lelse
    movsx edi, byte [rbp-2]     ; arg 1 = p.x (ABI)
    movsx esi, byte [rbp-1]     ; arg 2 = p.y
    call sum
    mov  qword [rbp-8], rax     ; total = sum(...)  (native word)
    jmp  .Ldone
.Lelse:
    mov  rax, 0                 ; return 0
.Ldone:
    leave
    ret

sum:
    push rbp
    mov  rbp, rsp
    movsx rax, edi           ; a  (already in a register via the ABI)
    movsx rcx, esi           ; b
    add   rax, rcx           ; a + b
    pop  rbp
    ret
```

The two listings tell the whole story of this document: same shape, same
mnemonics, same ABI — but C sits under a compiler that rewrites it on a
whim, and Radic promises that what you wrote is what the machine runs, down
to the width of a single byte in memory.

## Summary

| feature          | C at `-O0`                            | Radic (planned)                       |
| ---------------- | ------------------------------------- | ------------------------------------- |
| `int`/`number`   | 4-byte slots, 32-bit ops              | mold-width memory, 64-bit ops         |
| small values     | still 4 bytes                         | smallest mold that fits (1/2/4/8)     |
| load/store width | fixed by type                         | `movsx`/`movzx` load, `mov` store     |
| division         | `cdq; idiv` (32-bit)                  | `cqo; idiv` (64-bit)                  |
| `switch`         | chain at `-O0`, jump table at `-O2`   | always a `cmp`/`je` chain             |
| short-circuit    | yes (`&&`, `||`)                      | undecided (`and`/`or`)                |
| `& * ->`         | standard memory operands              | standard, once pointer types land     |
| function calls   | System V ABI                          | System V ABI                          |
| implicit opt     | free (folding, jump tables, ...)      | none, ever                            |

Related: `docs/design/c_radic.md` (C vs Radic syntax) and
`docs/design/number.md` (the mold and its load/store table).