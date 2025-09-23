# llc


llc (learning language compiler) is my own programming language implementation to learn how programming languages are built starting from lexing and up to compiler optimizations.

It not a compiled language yet it's interpreted language, and it doesn't support a lot of obvious things, but it will at some point.

It's supposed to be life-lasting project. I'll implement new features to it when I'll have inspiration and free time.

Now it supports:
- integers
- booleans
- strings
- arrays
- hashes
- prefix-, infix- and index operators
- conditionals
- global and local bindings
- first-class functions
- return statements
- closures


Quick look at syntax:
```
let name = "Monkey";
let age = 1;
let inspirations = ["Scheme", "Lisp", "JavaScript", "Clojure"];

let fibonacci = fn(x) {
    if (x == 0) {
        0
    } else {
        if (x == 1) {
            return 1;
        } else {
            fibonacci(x - 1) + fibonacci(x - 2);
        }
    }
};


let map = fn(arr, f) {
    let iter = fn(arr, accumulated) {
        if (len(arr) == 0) {
            accumulated
        } else {
            iter(rest(arr), push(accumulated, f(first(arr))));
        }
    };
    iter(arr, []);
};

let numbers = [1, 1 + 1, 4 - 1, 2 * 2, 2 + 3, 12 / 2];
map(numbers, fibonacci);
// => returns: [1, 1, 2, 3, 5, 8]
```
