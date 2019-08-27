	.intel_syntax noprefix
	.text
	.globl	main
	.def	main
main:
	push	rbp
	mov	rbp, rsp
	sub	rsp, 32
	mov	eax, 7
	leave
	ret