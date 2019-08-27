	.file	"testing.cpp"
	.intel_syntax noprefix
	.text
	.def	__main;	.scl	2;	.type	32;	.endef
	.section .rdata,"dr"
.LC0:
	.ascii "Hello World!\0"
	.text
	.globl	main
	.def	main;	.scl	2;	.type	32;	.endef
main:
	push	rbp
	mov	rbp, rsp
	sub	rsp, 32
	call	__main
	lea	rcx, .LC0[rip]
	call	puts
	mov	eax, 7
	leave
	ret
	.ident	"GCC: (Rev2, Built by MSYS2 project) 8.3.0"
	.def	puts;	.scl	2;	.type	32;	.endef
