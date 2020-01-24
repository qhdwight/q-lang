.data
_string1:
    .string "Hello world!\n"

_string2:
    .string "Yeah\n"

.text
.intel_syntax noprefix
.globl _main
_main:
    push rbp
    mov rbp, rsp
    
    mov dword ptr [rbp - 4], 1
    _loopCheck1:
    cmp dword ptr [rbp - 4], 3
    jge _loopContinue1
    jmp _loopBody1
    _loopBody1:
    lea rax, [rip + _string1]
    mov rsi, rax # Pointer to string
    mov rdx, 13 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    mov eax, dword ptr [rbp - 4]
    add eax, 1
    mov dword ptr [rbp - 4], eax
    jmp _loopCheck1
    _loopContinue1:
    mov dword ptr [rbp - 8], 0
    _loopCheck2:
    cmp dword ptr [rbp - 8], 4
    jge _loopContinue2
    jmp _loopBody2
    _loopBody2:
    lea rax, [rip + _string2]
    mov rsi, rax # Pointer to string
    mov rdx, 5 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    mov eax, dword ptr [rbp - 8]
    add eax, 1
    mov dword ptr [rbp - 8], eax
    jmp _loopCheck2
    _loopContinue2:
    mov eax, 0
    
    pop rbp
    ret

