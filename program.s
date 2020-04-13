.data
.text
.intel_syntax noprefix
.globl _main
_main:
    push rbp
    mov rbp, rsp
    
    mov dword ptr [rbp - 4], 2
    mov dword ptr [rbp - 8], 0
    mov dword ptr [rbp - 12], 0
    mov dword ptr [rbp - 16], 3
    mov eax, dword ptr [rbp - 12]
    add eax, dword ptr [rbp - 16]
    mov dword ptr [rbp - 8], eax
    mov eax, 4
    
    pop rbp
    ret

