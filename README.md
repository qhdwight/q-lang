# Q-Lang

Hello! This is the repository for my work in progress language. It is a fun way for me to learn about assembly and language design.

Here is a sneak-peak of what will be possible:

```
pkg main {
    def void -> u32 {
        imp main {
            u32 {
              x <- 2;
              y <- x + 6;
            }
            out x + y;
        }
    }
    def u32, u32 -> u32 {
        imp calculate(a, b) {
            out a + b;
        }
    }
}
```

## Working examples

[Calculator example](calculator.qq)

The goal is to have a concise language. Heavily influenced from Go and maybe a little bit from Python.
