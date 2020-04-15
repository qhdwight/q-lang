.data
.text
.intel_syntax noprefix
.globl _main
mov eax, 0
_digitToChar:
    movabs r8, -3689348814741910323
    xor r9, r9
    mov byte ptr [rsi - 1], 10
    dec rsi
    inc r9
_charLoop:
    movsxd rax, edi
    mul r8
    shr rdx, 3
    lea eax, [rdx + rdx]
    lea eax, [rax + 4*rax]
    mov ecx, edi
    sub ecx, eax
    or cl, 48
    mov byte ptr [rsi - 1], cl
    dec rsi
    inc r9
    cmp rdi, 9
    mov rdi, rdx
    ja _charLoop
    ret
_main:
    push rbp
    mov rbp, rsp
    sub rsp, 64
    
    mov dword ptr [rbp - 4], 43
    mov dword ptr [rbp - 8], 2
    mov dword ptr [rbp - 20], 0 # Counter
    _loopCheck1:
    cmp dword ptr [rbp - 20], 5
    jge _loopContinue1
    jmp _loopBody1
    _loopBody1:
    mov dword ptr [rbp - 24], 0
    mov eax, dword ptr [rbp - 4]
    add eax, dword ptr [rbp - 8]
    mov dword ptr [rbp - 24], eax
    mov edi, dword ptr [rbp - 24] # Integer argument
    lea rsi, [rbp - 32] # Buffer pointer argument
    call _digitToChar
    mov rdx, r9 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    mov eax, dword ptr [rbp - 20]
    inc eax
    mov dword ptr [rbp - 20], eax
    jmp _loopCheck1
    _loopContinue1:
    
    add rsp, 64
    pop rbp
    ret
