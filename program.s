.data
_string1:
    .string "Enter first number!\n"
_string2:
    .string "Enter second number!\n"
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
    mov dl, byte ptr [rbp - 8]
    mov byte ptr [rbp - 12], dl
    mov dl, byte ptr [rbp - 7]
    mov byte ptr [rbp - 11], dl
    mov dl, byte ptr [rbp - 6]
    mov byte ptr [rbp - 10], dl
    mov dl, byte ptr [rbp - 5]
    mov byte ptr [rbp - 9], dl
    # Expression base {typeName:i32 varName: stackPos:12}
    # Copy {typeName:i32 varName: stackPos:12} to {typeName:i32 varName:a stackPos:4}
    mov dl, byte ptr [rbp - 12]
    mov byte ptr [rbp - 4], dl
    mov dl, byte ptr [rbp - 11]
    mov byte ptr [rbp - 3], dl
    mov dl, byte ptr [rbp - 10]
    mov byte ptr [rbp - 2], dl
    mov dl, byte ptr [rbp - 9]
    mov byte ptr [rbp - 1], dl
    mov dword ptr [rbp - 12], 0 # Integer literal
    # Copy {typeName:i32 varName: stackPos:12} to {typeName:i32 varName: stackPos:16}
    mov dl, byte ptr [rbp - 12]
    mov byte ptr [rbp - 16], dl
    mov dl, byte ptr [rbp - 11]
    mov byte ptr [rbp - 15], dl
    mov dl, byte ptr [rbp - 10]
    mov byte ptr [rbp - 14], dl
    mov dl, byte ptr [rbp - 9]
    mov byte ptr [rbp - 13], dl
    # Expression base {typeName:i32 varName: stackPos:16}
    # Copy {typeName:i32 varName: stackPos:16} to {typeName:i32 varName:b stackPos:8}
    mov dl, byte ptr [rbp - 16]
    mov byte ptr [rbp - 8], dl
    mov dl, byte ptr [rbp - 15]
    mov byte ptr [rbp - 7], dl
    mov dl, byte ptr [rbp - 14]
    mov byte ptr [rbp - 6], dl
    mov dl, byte ptr [rbp - 13]
    mov byte ptr [rbp - 5], dl
    mov dword ptr [rbp - 16], 0 # Integer literal
    # Copy {typeName:i32 varName: stackPos:16} to {typeName:i32 varName: stackPos:20}
    mov dl, byte ptr [rbp - 16]
    mov byte ptr [rbp - 20], dl
    mov dl, byte ptr [rbp - 15]
    mov byte ptr [rbp - 19], dl
    mov dl, byte ptr [rbp - 14]
    mov byte ptr [rbp - 18], dl
    mov dl, byte ptr [rbp - 13]
    mov byte ptr [rbp - 17], dl
    # Expression base {typeName:i32 varName: stackPos:20}
    # Copy {typeName:i32 varName: stackPos:20} to {typeName:i32 varName:c stackPos:12}
    mov dl, byte ptr [rbp - 20]
    mov byte ptr [rbp - 12], dl
    mov dl, byte ptr [rbp - 19]
    mov byte ptr [rbp - 11], dl
    mov dl, byte ptr [rbp - 18]
    mov byte ptr [rbp - 10], dl
    mov dl, byte ptr [rbp - 17]
    mov byte ptr [rbp - 9], dl
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
    mov dword ptr [rbp - 4], eax
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
    # Copy {typeName:i32 varName:a stackPos:4} to {typeName:i32 varName: stackPos:48}
    mov dl, byte ptr [rbp - 4]
    mov byte ptr [rbp - 48], dl
    mov dl, byte ptr [rbp - 3]
    mov byte ptr [rbp - 47], dl
    mov dl, byte ptr [rbp - 2]
    mov byte ptr [rbp - 46], dl
    mov dl, byte ptr [rbp - 1]
    mov byte ptr [rbp - 45], dl
    # Expression base {typeName:i32 varName: stackPos:48}
    mov eax, dword ptr [rbp - 48]
    add eax, dword ptr [rbp - 8]
    mov dword ptr [rbp - 48], eax
    mov edi, dword ptr [rbp - 48] # Integer argument
    lea rsi, [rbp - 48] # Buffer pointer argument
    call _digitToChar
    mov rdx, r9 # Size
    mov rax, 1 # Write system call identifier
    mov rdi, 1 # Standard output file descriptor
    syscall
    mov eax, 0
    
    add rsp, 1024
    pop rbp
    ret
