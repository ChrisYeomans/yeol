LINE_MAX equ 1024
%include "string.inc"
%include "util.inc"
SECTION .text
global _start
_start:
mov rbp, rsp
    sub rsp, 8
    read 0, line, LINE_MAX
    mov rdi, line
    call strlen
    mov rdi, line
    mov rsi, rax
    call parse_uint
    mov qword [rbp - 8], rax
    mov rax, qword [rbp - 8]
    mov rdi, 1
    mov rsi, rax
    call write_uint
    add rsp, 8
    exit_program 0
SECTION .bss
    line: resb LINE_MAX
