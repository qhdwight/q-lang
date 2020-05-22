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
.globl main
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
main:
    push rbp
    mov rbp, rsp
    sub rsp, 1024
    
    mov dword ptr [rbp - 8], 0 # Integer literal
    # Copy {typeName:i32 varName: stackPos:8} to {typeName:i32 varName: stackPos:12}
    mov edx, dword ptr [rbp - 8]
    mov dword ptr [rbp - 12], edx
    # Expression base {typeName:i32 varName: stackPos:12}
    # Copy {typeName:i32 varName: stackPos:12} to {typeName:i32 varName:c stackPos:4}
    mov edx, dword ptr [rbp - 12]
    mov dword ptr [rbp - 4], edx
    # Copy {typeName:test varName: stackPos:20} to {typeName:test varName: stackPos:28}
    mov edx, dword ptr [rbp - 20]
    mov dword ptr [rbp - 28], edx
    mov edx, dword ptr [rbp - 16]
    mov dword ptr [rbp - 24], edx
    # Expression base {typeName:test varName: stackPos:28}
    # Copy {typeName:test varName: stackPos:28} to {typeName:test varName:z stackPos:12}
    mov edx, dword ptr [rbp - 28]
    mov dword ptr [rbp - 12], edx
    mov edx, dword ptr [rbp - 24]
    mov dword ptr [rbp - 8], edx
    lea rax, [rip + _string1]
    mov rsi, rax # Pointer to string
    mov rdx, 20 # Size
    mov rax, 1 # Write
    mov rdi, 1 # Standard output
    syscall
    lea rsi, [rbp - 28] # Pointer to ASCII buffer
    mov rdx, 16 # Size
    mov rax, 0 # Read system call identifier
    mov rdi, 0 # Standard input file descriptor
    syscall
    mov esi, eax # Length of characters
    lea rdi, [rbp - 28] # Pointer to ASCII buffer
    call _asciiToInt32
    mov dword ptr [rbp - 12], eax
    lea rax, [rip + _string2]
    mov rsi, rax # Pointer to string
    mov rdx, 21 # Size
    mov rax, 1 # Write
    mov rdi, 1 # Standard output
    syscall
    lea rsi, [rbp - 44] # Pointer to ASCII buffer
    mov rdx, 16 # Size
    mov rax, 0 # Read system call identifier
    mov rdi, 0 # Standard input file descriptor
    syscall
    mov esi, eax # Length of characters
    lea rdi, [rbp - 44] # Pointer to ASCII buffer
    call _asciiToInt32
    mov dword ptr [rbp - 8], eax
    # Copy {typeName:i32 varName: stackPos:12} to {typeName:i32 varName: stackPos:48}
    mov edx, dword ptr [rbp - 12]
    mov dword ptr [rbp - 48], edx
    # Expression base {typeName:i32 varName: stackPos:48}
    mov eax, dword ptr [rbp - 48]
    add eax, dword ptr [rbp - 8]
    mov dword ptr [rbp - 48], eax
    # Copy {typeName:i32 varName: stackPos:48} to {typeName:i32 varName:c stackPos:4}
    mov edx, dword ptr [rbp - 48]
    mov dword ptr [rbp - 4], edx
    # Copy {typeName:i32 varName:c stackPos:4} to {typeName:i32 varName: stackPos:48}
    mov edx, dword ptr [rbp - 4]
    mov dword ptr [rbp - 48], edx
    # Expression base {typeName:i32 varName: stackPos:48}
    mov edi, dword ptr [rbp - 48] # Integer argument
    lea rsi, [rbp - 48] # Buffer pointer argument
    call _digitToChar
    mov rdx, r9 # Size
    mov rax, 1 # Write system call identifier
    mov rdi, 1 # Standard output file descriptor
    syscall
    mov dword ptr [rbp - 48], 25 # Integer literal
    mov eax, dword ptr [rbp - 4]
    cmp eax, dword ptr [rbp - 48]
    jne _iff0
    mov eax, dword ptr [rbp - 12]
    mov dword ptr [rbp - 52], eax # Counter
    _loopCheck1:
    mov eax, dword ptr [rbp - 52]
    cmp eax, dword ptr [rbp - 8]
    jge _loopContinue1
    jmp _loopBody1
    _loopBody1:
    lea rax, [rip + _string3]
    mov rsi, rax # Pointer to string
    mov rdx, 8 # Size
    mov rax, 1 # Write
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
    mov rax, 1 # Write
    mov rdi, 1 # Standard output
    syscall
    _iff1:
    mov eax, 0
    
    add rsp, 1024
    pop rbp
    ret
