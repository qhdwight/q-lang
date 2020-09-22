# [Try it Online!](http://www.quintin.cloud/qlang)
##### May take a bit to load the first time (sets up Docker image)

# Q-Lang

Hello! This is the repository for my work in progress language. It is a fun way for me to learn about assembly and language design.

It currently compiles Q-Lang into x86 assembly and then assembles it with GCC (so you must have that in your path). Mac and Linux only for now.

## Examples

### [Calculator](examples/calculator.qq)

```
pkg main {
    dat inputs_d {
        u32 { a; b; }
    }
    def u32, u32 -> u32 {
        imp calculate a,b {
            out a + b;
        }
    }
    def nil -> u32 {
        imp main {
            inputs_d { i <- ?; }
            wln 'Enter first number!';
            rln $i.a;

            wln 'Enter second number!';
            rln $i.b;

            u32 { c <- calculate i.a, i.b; }
            wln c;

            iff c = 25 {
                itr i.a...i.b {
                    wln 'Secret!';
                }
            } els {
                wln 'No secret!';
            }
            out 0;
        }
    }
}
```

## Philosophy

The goal is to have a concise language. Heavily influenced from Go and maybe a little bit from Python.

## Overview

### 1. Hello world!

```
pkg main {
    def nil -> u32 {
        imp main {
            wln 'Hello world!';
            out 0;
        }
    }
}
```

For now, only one package is supported, and it has to be `pkg main`.

`def nil -> u32` defines that the following block must contain implementations that take no parameters and return a 32-bit unsigned integer.

`imp main` defines a function named main.

`wln 'Hello world!';` Prints to standard output. Remember semicolons end statements and newlines do not implcitly do so!

`out 0` give back ("return") a zeroed unsigned 32-bit integer.

### 1.5. Compiling q-lang files

`./q --input <qq file>`

It will generate an executable named `program` in the current working directory, which can be via `./program`

I provided a script `run.zsh` which can be used to automatically run the compiled program `./run.zsh <qq file>`

### 2. Defining data ("struct" in other languages)

Place inside of main package for now.

```
dat test_d {
    u32 { a; b; }
}
```

`dat test_d` defines a data type named `test_d`.

`u32 { a; b; }` defines two fields named a and b that are unsigned 32-bit integers.

### 3. Defining variables inside of functions

```
u32 { x <- 1; y <- 2; }
test_d { t <- test_d{a <- 3; b <- 2+2; }; }
# Prints 1
wln x;
# Prints 2
wln y;
# Prints 3
wln t.a;
# Prints 4
wln t.b;
```

The `<-` language construct is an assignment operator. It moves what is on the right hand side into a named variable.

`u32 { x <- 3; y <- 4; }` allows us to define mutliple variables of the same type easily.

`test_d { t <- test_d{a <- 1; b <- 2; }; }` defines a variable named `t` with data type `test_d`, which is defined in #2. We can also set the fields inside of the data, which are unsigned 32-bit integers with names `a` and `b`. Remember semicolon placements!

Uninitialized data can be achieved by using the `?` language construct. `u32 { x <- ?; }` for example does not garuntee that `x` has any reasonable value.

### 4. More complicated functions
```
pkg main {
    def u32, u32 -> u32 {
        imp calculate a,b {
            out a + b;
        }
    }
}
```

Multiple inputs types can be achieved by comma separation. Notice that in the definition, only type names are present, and in actual implementations we provide names that are bound to the scope of the funciton.

### 5. Iteration

```
itr 0...3 {
    wln 'Hello world!';
}
```

The `...` language construct can be used to iterate between two unsigned 32-bit integers. Inclusive to exclusive, so the above prints three times.

### 6. Control flow

```
iff 25 = 25 {
    wln 'Secret!';
} els {
    wln 'No secret!';
}
```

### 7. Math

The `+`, `-`, `*`, and `/` operators are currently supported between 32-bit unsigned integers and 32-bit floating points.
