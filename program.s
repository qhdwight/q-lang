.data
_message:
	.string "Hello World!\n"

.text
.intel_syntax noprefix
.globl	_main
_main:
	push rbp
	mov	rbp, rsp


