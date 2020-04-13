.data
_string1:
    .string "Hello World\n"

.text
.intel_syntax noprefix
.globl _main
_main:
    push rbp
    mov rbp, rsp
    
    mov dword ptr [rbp - 0], 0
    _loopCheck1:
    cmp dword ptr [rbp - 0], 3
    jge _loopContinue1
    jmp _loopBody1
    _loopBody1:
    lea rax, [rip + _string1]
    mov rsi, rax # Pointer to string
    mov rdx, 12 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    mov eax, dword ptr [rbp - 0]
    add eax, 1
    mov dword ptr [rbp - 0], eax
    jmp _loopCheck1
    _loopContinue1:
    
    pop rbp
    ret

