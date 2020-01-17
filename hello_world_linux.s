.data
_message:
	.string "Hello World!\n"

.text
.intel_syntax noprefix
.globl main	
main:
	push rbp
	mov	rbp, rsp

	lea	rax, [rip + _message]
	mov rsi, rax # Pointer to string
	mov rdx, 13 # Size
	mov rax, 1 # Write
	mov rdi, 1 # Standard output
	syscall

	mov	eax, 0
	pop	rbp
	ret

