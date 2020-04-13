.text
.intel_syntax noprefix
.globl _main
_print_num:
  mov rax, rdi

  mov rcx, 0xA
  push rcx
  mov rsi, rsp
  sub rsp, 16

  .to_digit:
  xor edx, edx
  div ecx

  add edx, '0'
  dec rsi
  mov [rsi], dl

  test eax, eax
  jnz .to_digit

_main:
  push rbp
  mov rbp, rsp

  lea rdx, [rsp+16 + 1]
  sub rdx, rsi
  mov rax, 0x2000004 # Write
  mov rdi, 1 # Standard output
  syscall

  mov  eax, 0
  pop  rbp
  ret
