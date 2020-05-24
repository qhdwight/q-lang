.data
_string1:
    .string "Hello world!\n"
.text
.intel_syntax noprefix
.globl _main
_main:
    push rbp
    mov rbp, rsp
    sub rsp, 1024
    
    mov dword ptr [rbp - 4], 0 # Integer literal
    mov dword ptr [rbp - 8], 3 # Integer literal
    mov eax, dword ptr [rbp - 4]
    mov dword ptr [rbp - 12], eax # Counter
    _loopCheck1:
    mov eax, dword ptr [rbp - 12]
    cmp eax, dword ptr [rbp - 8]
    jge _loopContinue1
    jmp _loopBody1
    _loopBody1:
    lea rax, [rip + _string1]
    mov rsi, rax # Pointer to string
    mov rdx, 13 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    mov eax, dword ptr [rbp - 12]
    inc eax
    mov dword ptr [rbp - 12], eax
    jmp _loopCheck1
    _loopContinue1:
    
    add rsp, 1024
    pop rbp
    ret
