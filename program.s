.data
_string1:
    .string "Hello World!\n"

_string2:
    .string "Test\n"

_string3:
    .string "Ok\n"

.text
.intel_syntax noprefix
.globl _main
_main:
    push rbp
    mov rbp, rsp
    lea rax, [rip + _string1]
    mov rsi, rax # Pointer to string
    mov rdx, 13 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    lea rax, [rip + _string2]
    mov rsi, rax # Pointer to string
    mov rdx, 5 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    lea rax, [rip + _string3]
    mov rsi, rax # Pointer to string
    mov rdx, 3 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    pop rbp
    ret

