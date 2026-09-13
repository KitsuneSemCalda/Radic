# Radic vs C: syntax comparison

This document compares Radic's syntax with C's, side by side, based on what
the current parser and lexer actually implement. Things not yet supported
are marked as such, so they don't become a false promise.

## Functions

```c
int add(int a, int b) {
    return a + b;
}
```

```radic
func add(number a, number b) -> number {
    return a + b;
}
```

- C writes the return type before the name; Radic uses `->` after the
  parameter list.
- Parameters are C-style in both: type first, then the name.
- `return;` with no value is accepted by Radic's parser, even in a function
  that declares a return type (`internal/parser/parser.go`).

## Variables and types

```c
int x = 5;
unsigned short y = 300;
```

```radic
number x = 5;       // inferred mold: 8-bit (i8)
number y = 300;     // grows by fixed point: 16-bit (i16)
```

- In C, width/signedness are written by hand; in Radic the mold (width +
  signedness) is inferred during analysis and never truncates (see
  `docs/design/number.md`).
- Still no pointer type in type position: `int* p` does not exist — see
  Pointers.
- `string` and `bool` are Radic builtins; C uses `char*` and
  `_Bool`/`<stdbool.h>`.
- Literals: Radic has `true`, `false` and `nil`; C uses `true/false` via
  `<stdbool.h>` and `NULL`.

## if / while / for

The three constructs follow C's form, with mandatory parentheses:

```c
if (x > 0) {
    x--;
} else {
    x++;
}
for (int i = 0; i < 10; i++) {
    ...
}
while (x > 0) {
    ...
}
```

```radic
if (x > 0) {
    x--;
} else {
    x++;
}
for (i = 0; i < 10; i++) {
    ...
}
while (x > 0) {
    ...
}
```

Radic's `for` loop accepts optional init/cond/post, same as C's
(`internal/parser/parser.go`, `parseFor`).

## switch

```c
switch (x) {
    case 1: foo(); break;
    case 2: bar();
    default: baz();
}
```

```radic
switch (x) {
    case 1 { foo(); }
    case 2 { bar(); }
    default { baz(); }
}
```

- In Radic each `case` is a `{ }` block with an *implicit break*: there is no
  fallthrough, no `break` needed.
- The discriminator uses parentheses, as in C (`switch (x) {`).

## Logical operators

| C        | Radic   |
| -------- | ------- |
| `&&`     | `and`   |
| `\|\|`   | `or`    |
| `!`      | `not`   |

Bitwise and compound assignment are identical: `&`, `|`, `^`, `~`, `<<`, `>>`,
`+=`, `-=`, `*=`, `/=`, `%=`, `++`, `--` (prefix and postfix).

## struct / enum / union

```c
typedef struct { int x; int y; } Point;
Point p;
p.x = 3;
```

```radic
struct Point {
    number x;
    number y;
}
p = Point...      // (see note below)
p.x = 3;
```

- Radic has no `typedef`: `struct Name { ... }` is already direct. A field's
  name comes from `number` (or a pinned width, e.g. `field: u32`).
- The parser does not yet implement aggregate initialization (`Point{3, 4}`),
  so the example is limited to field-by-field assignment.
- `enum Name { A = 1, B, C }` has the same form; the last variant doesn't
  require a trailing comma in Radic.
- Member access uses `.` and `->`, as in C.

## Pointers

Radic has the unary operators `*` (deref), `&` (address-of) and `->`
(dereferenced member), but still does **not** have pointer type syntax
(`int* p`) nor array syntax (`int a[10]`). Pointers today only appear via
operators, with no explicit pointer type declaration — pending functionality.

## include

```c
#include <stdio.h>
```

```radic
include <stdio.h>
```

- C uses the preprocessor; Radic has the `include` keyword, still **not
  implemented** — the parser responds with an error
  (`internal/parser/parser.go`).
- Radic has no macros or preprocessor.

## Full program side by side

```c
struct Point { int x; int y; };

int sum(int a, int b) {
    return a + b;
}

void main() {
    struct Point p;
    p.x = 3;
    p.y = 4;

    if ((p.x > 0) && (p.y > 0)) {
        int total = sum(p.x, p.y);
    }
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
    }
}
```

Note: the Radic example above expresses intent, not the parser's current
state — struct construction, variable declaration for structs, and pointers
in type position are still to be implemented.

## Summary

| feature                | C                                      | Radic                                |
| ----------------------- | -------------------------------------- | ------------------------------------ |
| function return         | `int f(...)`                           | `func f(...) -> number`              |
| integer width           | explicit (`int`, `unsigned short`)     | `number` with inferred mold          |
| logical operators       | `&& \|\| !`                            | `and or not`                         |
| switch                  | `case 1: ... break;` (fallthrough)     | `case 1 { ... }` (implicit break)    |
| typedef                 | yes                                    | no (direct names)                    |
| preprocessor / macros   | yes                                    | no (`include` pending)               |
| pointer/array types     | yes                                    | pending                              |
| `bool`/`string`/`error` | via libs                               | builtins                             |
