gcc -S -fverbose-asm -O0 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -masm=intel program.cpp
gcc program.s -o program
