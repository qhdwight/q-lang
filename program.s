.data
_string1:
    .string "Hello world!\n"

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
    mov eax, 1
    
    pop rbp
    ret

