.data
.text
.intel_syntax noprefix
.globl _main
mov eax, 0
_digitToChar:
    movabs r8, -3689348814741910323
_charLoop:
    mov rax, rdi
    mul r8
    shr rdx, 3
    lea eax, [rdx + rdx]
    lea eax, [rax + 4*rax]
    mov ecx, edi
    sub ecx, eax
    or cl, 48
    mov byte ptr [rsi - 1], cl
    dec rsi
    cmp rdi, 9
    mov rdi, rdx
    ja _charLoop
    mov rax, rsi
    ret
_main:
    push rbp
    mov rbp, rsp
    
    mov dword ptr [rbp - 0], 21
    mov dword ptr [rbp - 4], 0
    mov dword ptr [rbp - 8], 31
    mov eax, dword ptr [rbp - 0]
    add eax, dword ptr [rbp - 8]
    mov dword ptr [rbp - 4], eax
    mov edi, dword ptr [rbp - 4]
    lea rsi, dword ptr [rbp - 8]
    call _digitToChar
    mov rdx, 5 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    
    pop rbp
    ret
