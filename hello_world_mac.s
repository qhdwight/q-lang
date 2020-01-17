.data
_message:
	.string "Hello World!\n"

.text
.intel_syntax noprefix
.globl	_main
_main:
	push rbp
	mov	rbp, rsp

	lea	rax, [rip + _message]
	mov rsi, rax # Pointer to string
	mov rdx, 13 # Size
	mov rax, 0x2000004 # Write
	mov rdi, 1 # Standard output
	syscall

	mov	eax, 2
	pop	rbp
	ret
