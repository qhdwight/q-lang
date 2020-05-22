.data
.text
.intel_syntax noprefix
.globl main
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
    sub rsp, 64
    
    mov dword ptr [rbp - 8], 2 # Integer literal
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
    # Copy {typeName:i32 varName: stackPos:12} to {typeName:i32 varName:c stackPos:4}
    mov dl, byte ptr [rbp - 12]
    mov byte ptr [rbp - 4], dl
    mov dl, byte ptr [rbp - 11]
    mov byte ptr [rbp - 3], dl
    mov dl, byte ptr [rbp - 10]
    mov byte ptr [rbp - 2], dl
    mov dl, byte ptr [rbp - 9]
    mov byte ptr [rbp - 1], dl
    mov dword ptr [rbp - 24], 3 # Integer literal
    # Copy {typeName:i32 varName: stackPos:24} to {typeName:i32 varName: stackPos:28}
    mov dl, byte ptr [rbp - 24]
    mov byte ptr [rbp - 28], dl
    mov dl, byte ptr [rbp - 23]
    mov byte ptr [rbp - 27], dl
    mov dl, byte ptr [rbp - 22]
    mov byte ptr [rbp - 26], dl
    mov dl, byte ptr [rbp - 21]
    mov byte ptr [rbp - 25], dl
    # Expression base {typeName:i32 varName: stackPos:28}
    # Copy {typeName:i32 varName: stackPos:28} to {typeName:i32 varName: stackPos:20}
    mov dl, byte ptr [rbp - 28]
    mov byte ptr [rbp - 20], dl
    mov dl, byte ptr [rbp - 27]
    mov byte ptr [rbp - 19], dl
    mov dl, byte ptr [rbp - 26]
    mov byte ptr [rbp - 18], dl
    mov dl, byte ptr [rbp - 25]
    mov byte ptr [rbp - 17], dl
    mov dword ptr [rbp - 32], 2 # Integer literal
    # Copy {typeName:i32 varName: stackPos:32} to {typeName:i32 varName: stackPos:36}
    mov dl, byte ptr [rbp - 32]
    mov byte ptr [rbp - 36], dl
    mov dl, byte ptr [rbp - 31]
    mov byte ptr [rbp - 35], dl
    mov dl, byte ptr [rbp - 30]
    mov byte ptr [rbp - 34], dl
    mov dl, byte ptr [rbp - 29]
    mov byte ptr [rbp - 33], dl
    # Expression base {typeName:i32 varName: stackPos:36}
    # Copy {typeName:i32 varName: stackPos:36} to {typeName:i32 varName: stackPos:16}
    mov dl, byte ptr [rbp - 36]
    mov byte ptr [rbp - 16], dl
    mov dl, byte ptr [rbp - 35]
    mov byte ptr [rbp - 15], dl
    mov dl, byte ptr [rbp - 34]
    mov byte ptr [rbp - 14], dl
    mov dl, byte ptr [rbp - 33]
    mov byte ptr [rbp - 13], dl
    # Copy {typeName:test varName: stackPos:20} to {typeName:test varName: stackPos:44}
    mov dl, byte ptr [rbp - 20]
    mov byte ptr [rbp - 44], dl
    mov dl, byte ptr [rbp - 19]
    mov byte ptr [rbp - 43], dl
    mov dl, byte ptr [rbp - 18]
    mov byte ptr [rbp - 42], dl
    mov dl, byte ptr [rbp - 17]
    mov byte ptr [rbp - 41], dl
    mov dl, byte ptr [rbp - 16]
    mov byte ptr [rbp - 40], dl
    mov dl, byte ptr [rbp - 15]
    mov byte ptr [rbp - 39], dl
    mov dl, byte ptr [rbp - 14]
    mov byte ptr [rbp - 38], dl
    mov dl, byte ptr [rbp - 13]
    mov byte ptr [rbp - 37], dl
    # Expression base {typeName:test varName: stackPos:44}
    # Copy {typeName:test varName: stackPos:44} to {typeName:test varName:d stackPos:12}
    mov dl, byte ptr [rbp - 44]
    mov byte ptr [rbp - 12], dl
    mov dl, byte ptr [rbp - 43]
    mov byte ptr [rbp - 11], dl
    mov dl, byte ptr [rbp - 42]
    mov byte ptr [rbp - 10], dl
    mov dl, byte ptr [rbp - 41]
    mov byte ptr [rbp - 9], dl
    mov dl, byte ptr [rbp - 40]
    mov byte ptr [rbp - 8], dl
    mov dl, byte ptr [rbp - 39]
    mov byte ptr [rbp - 7], dl
    mov dl, byte ptr [rbp - 38]
    mov byte ptr [rbp - 6], dl
    mov dl, byte ptr [rbp - 37]
    mov byte ptr [rbp - 5], dl
    mov dword ptr [rbp - 16], 3 # Integer literal
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
    mov eax, dword ptr [rbp - 20]
    mov dword ptr [rbp - 24], 2 # Integer literal
    add eax, dword ptr [rbp - 24]
    mov dword ptr [rbp - 28], 3 # Integer literal
    add eax, dword ptr [rbp - 28]
    mov dword ptr [rbp - 20], eax
    # Copy {typeName:i32 varName: stackPos:20} to {typeName:i32 varName:c stackPos:4}
    mov dl, byte ptr [rbp - 20]
    mov byte ptr [rbp - 4], dl
    mov dl, byte ptr [rbp - 19]
    mov byte ptr [rbp - 3], dl
    mov dl, byte ptr [rbp - 18]
    mov byte ptr [rbp - 2], dl
    mov dl, byte ptr [rbp - 17]
    mov byte ptr [rbp - 1], dl
    mov dword ptr [rbp - 16], 5 # Integer literal
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
    # Copy {typeName:i32 varName: stackPos:20} to {typeName:i32 varName: stackPos:12}
    mov dl, byte ptr [rbp - 20]
    mov byte ptr [rbp - 12], dl
    mov dl, byte ptr [rbp - 19]
    mov byte ptr [rbp - 11], dl
    mov dl, byte ptr [rbp - 18]
    mov byte ptr [rbp - 10], dl
    mov dl, byte ptr [rbp - 17]
    mov byte ptr [rbp - 9], dl
    # Copy {typeName:i32 varName:c stackPos:4} to {typeName:i32 varName: stackPos:16}
    mov dl, byte ptr [rbp - 4]
    mov byte ptr [rbp - 16], dl
    mov dl, byte ptr [rbp - 3]
    mov byte ptr [rbp - 15], dl
    mov dl, byte ptr [rbp - 2]
    mov byte ptr [rbp - 14], dl
    mov dl, byte ptr [rbp - 1]
    mov byte ptr [rbp - 13], dl
    # Expression base {typeName:i32 varName: stackPos:16}
    mov eax, dword ptr [rbp - 16]
    mov dword ptr [rbp - 20], 32 # Integer literal
    add eax, dword ptr [rbp - 20]
    mov dword ptr [rbp - 16], eax
    # Copy {typeName:i32 varName: stackPos:16} to {typeName:i32 varName: stackPos:8}
    mov dl, byte ptr [rbp - 16]
    mov byte ptr [rbp - 8], dl
    mov dl, byte ptr [rbp - 15]
    mov byte ptr [rbp - 7], dl
    mov dl, byte ptr [rbp - 14]
    mov byte ptr [rbp - 6], dl
    mov dl, byte ptr [rbp - 13]
    mov byte ptr [rbp - 5], dl
    # Copy {typeName:i32 varName:c stackPos:4} to {typeName:i32 varName: stackPos:16}
    mov dl, byte ptr [rbp - 4]
    mov byte ptr [rbp - 16], dl
    mov dl, byte ptr [rbp - 3]
    mov byte ptr [rbp - 15], dl
    mov dl, byte ptr [rbp - 2]
    mov byte ptr [rbp - 14], dl
    mov dl, byte ptr [rbp - 1]
    mov byte ptr [rbp - 13], dl
    # Expression base {typeName:i32 varName: stackPos:16}
    mov edi, dword ptr [rbp - 16] # Integer argument
    lea rsi, [rbp - 16] # Buffer pointer argument
    call _digitToChar
    mov rdx, r9 # Size
    mov rax, 1 # Write
    mov rdi, 1 # Standard output
    syscall
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
    mov edi, dword ptr [rbp - 16] # Integer argument
    lea rsi, [rbp - 16] # Buffer pointer argument
    call _digitToChar
    mov rdx, r9 # Size
    mov rax, 1 # Write
    mov rdi, 1 # Standard output
    syscall
    # Copy {typeName:i32 varName: stackPos:8} to {typeName:i32 varName: stackPos:36}
    mov dl, byte ptr [rbp - 8]
    mov byte ptr [rbp - 36], dl
    mov dl, byte ptr [rbp - 7]
    mov byte ptr [rbp - 35], dl
    mov dl, byte ptr [rbp - 6]
    mov byte ptr [rbp - 34], dl
    mov dl, byte ptr [rbp - 5]
    mov byte ptr [rbp - 33], dl
    # Expression base {typeName:i32 varName: stackPos:36}
    mov edi, dword ptr [rbp - 36] # Integer argument
    lea rsi, [rbp - 36] # Buffer pointer argument
    call _digitToChar
    mov rdx, r9 # Size
    mov rax, 1 # Write
    mov rdi, 1 # Standard output
    syscall
    mov eax, 0
    
    add rsp, 64
    pop rbp
    ret
