# Q-Lang

Hello! This is the repository for my work in progress language. It is a fun way for me to learn about assembly and language design.

Here is a sneak peak:

```
pkg main {
    def void -> i32 {
        imp main {
            i32 {
              x = 2;
              y = x + 6;
            }
            out x + y;
        }
    }
    def i32, i32 -> i32 {
        imp calculate(a, b) {
            out a + b;
        }
    }
}
```

The goal is to have a concise language. Heavily influenced from Go and maybe a little bit from Python.
