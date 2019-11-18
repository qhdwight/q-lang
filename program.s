	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 10, 15	sdk_version 10, 15
	.intel_syntax noprefix
	.globl	__Z3betv                ## -- Begin function _Z3betv
	.p2align	4, 0x90
__Z3betv:                               ## @_Z3betv
## %bb.0:
	push	rbp
	mov	rbp, rsp
	pop	rbp
	ret
                                        ## -- End function
	.globl	_main                   ## -- Begin function main
	.p2align	4, 0x90
_main:                                  ## @main
## %bb.0:
	push	rbp
	mov	rbp, rsp
	sub	rsp, 16
	call	__Z3betv
	xor	eax, eax
	mov	dword ptr [rbp - 4], 6
	add	rsp, 16
	pop	rbp
	ret
                                        ## -- End function

.subsections_via_symbols
