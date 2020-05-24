.data
_string1:
    .string "Enter first number!\n"
_string2:
    .string "Enter second number!\n"
_string3:
    .string "Secret!\n"
_string4:
    .string "No secret!\n"
.text
.intel_syntax noprefix
.globl _main
_asciiToInt32:
    xor edx, edx
    xor eax, eax
    dec esi
    _whileBody:
    cmp esi, edx
    jle _whileEnd
    imul eax, eax, 10
    movsx ecx, byte ptr [rdi+rdx]
    inc rdx
    lea eax, [rax-48+rcx]
    jmp _whileBody
    _whileEnd:
    ret
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
_calculate:
    push rbp
    mov rbp, rsp
    sub rsp, 1024
    mov qword ptr [rbp - 8], rdi
    
    # Copy {typeName:u32 varName:a stackPos:12} to {typeName:u32 varName: stackPos:20}
    mov edx, dword ptr [rbp - 12]
    mov dword ptr [rbp - 20], edx
    # Expression base {typeName:u32 varName: stackPos:20}
    mov eax, dword ptr [rbp - 20]
    add eax, dword ptr [rbp - 16]
    mov dword ptr [rbp - 20], eax
    mov rdi, qword ptr [rbp - 8] # Caller return variable stack position
    mov eax, dword ptr [rbp - 20]
    mov dword ptr [rdi + 0], eax
    
    add rsp, 1024
    pop rbp
    ret
_main:
    push rbp
    mov rbp, rsp
    sub rsp, 1024
    
    lea rax, [rip + _string1]
    mov rsi, rax # Pointer to string
    mov rdx, 20 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    lea rsi, [rbp - 24] # Pointer to ASCII buffer
    mov rdx, 16 # Size
    mov rax, 0x2000003 # Read system call identifier
    mov rdi, 0 # Standard input file descriptor
    syscall
    mov esi, eax # Length of characters
    lea rdi, [rbp - 24] # Pointer to ASCII buffer
    call _asciiToInt32
    mov dword ptr [rbp - 8], eax
    lea rax, [rip + _string2]
    mov rsi, rax # Pointer to string
    mov rdx, 21 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    lea rsi, [rbp - 40] # Pointer to ASCII buffer
    mov rdx, 16 # Size
    mov rax, 0x2000003 # Read system call identifier
    mov rdi, 0 # Standard input file descriptor
    syscall
    mov esi, eax # Length of characters
    lea rdi, [rbp - 40] # Pointer to ASCII buffer
    call _asciiToInt32
    mov dword ptr [rbp - 4], eax
    lea rdi, dword ptr [rbp - 48] # Return stack pointer
    mov eax, dword ptr [rbp - 8]
    mov dword ptr [rsp - 28], eax # Copy parameter a to called function stack frame
    mov eax, dword ptr [rbp - 4]
    mov dword ptr [rsp - 32], eax # Copy parameter b to called function stack frame
    call _calculate
    # Copy {typeName:u32 varName: stackPos:48} to {typeName:u32 varName: stackPos:52}
    mov edx, dword ptr [rbp - 48]
    mov dword ptr [rbp - 52], edx
    # Expression base {typeName:u32 varName: stackPos:52}
    # Copy {typeName:u32 varName: stackPos:52} to {typeName:u32 varName:c stackPos:44}
    mov edx, dword ptr [rbp - 52]
    mov dword ptr [rbp - 44], edx
    # Copy {typeName:u32 varName:c stackPos:44} to {typeName:u32 varName: stackPos:48}
    mov edx, dword ptr [rbp - 44]
    mov dword ptr [rbp - 48], edx
    # Expression base {typeName:u32 varName: stackPos:48}
    mov edi, dword ptr [rbp - 48] # Integer argument
    lea rsi, [rbp - 48] # Buffer pointer argument
    call _digitToChar
    mov rdx, r9 # Size
    mov rax, 0x2000004 # Write system call identifier
    mov rdi, 1 # Standard output file descriptor
    syscall
    mov dword ptr [rbp - 48], 25 # Integer literal
    mov eax, dword ptr [rbp - 44]
    cmp eax, dword ptr [rbp - 48]
    jne _iff0
    mov eax, dword ptr [rbp - 8]
    mov dword ptr [rbp - 52], eax # Counter
    _loopCheck1:
    mov eax, dword ptr [rbp - 52]
    cmp eax, dword ptr [rbp - 4]
    jge _loopContinue1
    jmp _loopBody1
    _loopBody1:
    lea rax, [rip + _string3]
    mov rsi, rax # Pointer to string
    mov rdx, 8 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    mov eax, dword ptr [rbp - 52]
    inc eax
    mov dword ptr [rbp - 52], eax
    jmp _loopCheck1
    _loopContinue1:
    jmp _iff1
    _iff0:
    lea rax, [rip + _string4]
    mov rsi, rax # Pointer to string
    mov rdx, 11 # Size
    mov rax, 0x2000004 # Write
    mov rdi, 1 # Standard output
    syscall
    _iff1:
    mov dword ptr [rbp - 52], 0 # Integer literal
    # Copy {typeName:u32 varName: stackPos:52} to {typeName:u32 varName: stackPos:56}
    mov edx, dword ptr [rbp - 52]
    mov dword ptr [rbp - 56], edx
    # Expression base {typeName:u32 varName: stackPos:56}
    mov eax, dword ptr [rbp - 56] # Program exit code
    
    add rsp, 1024
    pop rbp
    ret
