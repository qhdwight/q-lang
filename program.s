	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 10, 15	sdk_version 10, 15
	.intel_syntax noprefix
	.globl	_main                   ## -- Begin function main
	.p2align	4, 0x90
_main:                                  ## @main
## %bb.0:
	push	rbp
	mov	rbp, rsp
	mov	dword ptr [rbp - 4], 0
	mov	dword ptr [rbp - 8], 0
LBB0_1:                                 ## =>This Inner Loop Header: Depth=1
	cmp	dword ptr [rbp - 8], 5
	jge	LBB0_4
## %bb.2:                               ##   in Loop: Header=BB0_1 Depth=1
	jmp	LBB0_3
LBB0_3:                                 ##   in Loop: Header=BB0_1 Depth=1
	mov	eax, dword ptr [rbp - 8]
	add	eax, 1
	mov	dword ptr [rbp - 8], eax
	jmp	LBB0_1
LBB0_4:
	mov	eax, dword ptr [rbp - 4]
	pop	rbp
	ret
                                        ## -- End function

.subsections_via_symbols
