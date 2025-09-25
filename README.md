# llc — a tiny learning language


llc (Learning Language Compiler) is a personal, long‑running project where I build a small programming language end‑to‑end. It’s a place to learn and experiment with lexing, parsing (Pratt parser), ASTs, an interpreter, a bytecode compiler and VM, a REPL, and language features like functions, closures, arrays, hashes, and even macros.

Work in progress: I add features from time to time. It’s intentionally small, readable, and test‑driven.


## Highlights
- End‑to‑end language pipeline: lexer → Pratt parser → AST → interpreter → (experimental) bytecode compiler + VM
- Ergonomic, expression‑oriented syntax with first‑class functions, closures, arrays, hashes, and conditionals
- Small standard library patterns (map/reduce implemented in the language)
- Built‑in functions: len, first, last, rest, push, print
- Macro system with quote/unquote for AST‑level metaprogramming
- Clean CLI and an interactive REPL


## Quick tour of the syntax
```
// Variables and basic types
let name = "Monkey";
let age = 1;
let flags = [true, false, true];
let kv = {"lang": "llc", "year": 2025};

// Functions and recursion
let fib = fn(x) {
  if (x == 0) { 0 }
  else {
    if (x == 1) { return 1; }
    else { fib(x - 1) + fib(x - 2); }
  }
};

// Higher‑order functions (map written in llc)
let map = fn(arr, f) {
  let iter = fn(arr, acc) {
    if (len(arr) == 0) { acc }
    else { iter(rest(arr), push(acc, f(first(arr)))); }
  };
  iter(arr, []);
};

let numbers = [1, 2, 3, 4, 5, 6];
map(numbers, fib);
// => [1, 1, 2, 3, 5, 8]

// Strings and concatenation
"hello, " + name

// Arrays, hashes, and indexing
[1, 2, 3][0];
{"a": 1, "b": 2}["b"]; // 2
```


## Macro system (quote / unquote)
Macros work at the AST level. You can build new syntax using quote/unquote.
```
let unless = macro(condition, consequence, alternative) {
  quote(
    if (!(unquote(condition))) {
      unquote(consequence);
    } else {
      unquote(alternative);
    }
  );
};

unless(10 > 5, print("not greater"), print("greater"));
// expands to: if (!(10 > 5)) { print("not greater") } else { print("greater") }
```


## What works today
Interpreter (tree‑walking)
- integers, booleans, strings
- arrays and hashes + indexing
- prefix and infix operators: !, -, +, -, *, /, <, >, ==, !=, string +
- conditionals (if/else)
- let bindings (global/local)
- first‑class functions, return, closures, higher‑order functions
- built‑ins: len, first, last, rest, push, print
- macros with quote/unquote

Bytecode compiler + VM (used by the REPL)
- initial support: integers and booleans, prefix/infix ops, expression evaluation and pop
- more features are being ported from the interpreter to the VM incrementally


## Try it
Prereqs: a recent Go toolchain (1.21+ should be fine).

Build the CLI
- `go build -o llc .`

Run the REPL (uses the VM)
- `./llc run`
- Type expressions like `1 + 2 * 3`, `true == !false`, etc.

Run a file (uses the interpreter)
- `./llc run examples/hello-world.llc`

Or without building
- `go run . run examples/hello-world.llc`


## Examples
- `examples/hello-world.llc`
- `std/array.llc` shows map and reduce implemented in llc


## Project layout
- `lang/lexer`, `lang/token` — lexical analysis
- `lang/ast` — AST nodes and tree utilities
- `lang/parser` — Pratt parser (operator precedence, calls, indexing)
- `lang/object` — runtime object system (ints, bools, strings, arrays, hashes, functions, macros, etc.)
- `lang/evaluator` — interpreter (tree‑walking) with built‑ins and macros
- `lang/compiler` — bytecode compiler (in progress)
- `lang/code` — instruction encoding/decoding helpers
- `lang/vm` — stack‑based VM (in progress, used by REPL)
- `lang/repl` — interactive shell
- `lang/cli` — cobra‑based CLI (llc run [file])
- `std/` — language‑level utilities (e.g., array.llc with map/reduce)
- `examples/` — small runnable snippets


## Design notes
- Expression‑oriented: blocks and conditionals return values; return short‑circuits from functions.
- Pratt parser keeps the grammar compact but expressive; tests exercise precedence and grouping thoroughly.
- Macros operate on ASTs, not strings; quote/unquote let you build and splice syntax safely.
- Two execution paths on purpose: the interpreter is feature‑complete first; the VM is catching up.


## Roadmap (ongoing)
- Port strings, arrays, conditionals, and function calls to the VM
- Module system/imports and a more complete std library
- Better errors and diagnostics
- Bytecode optimizations and simple compiler passes
- More examples and docs


## Why this project?
- A practical, readable codebase to demonstrate systems skills: parsing, runtimes, VMs, and language design.
- A safe place to try ideas and iterate over time.

If you’re curious or hiring, I’m happy to walk through the code and design trade‑offs.
