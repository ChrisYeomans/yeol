
; puts string length into rax
strlen:
  push   rdi
  push   rcx
  sub    rcx, rcx
  mov    rcx, -1
  sub    al, al
  cld    
  repne  scasb
  neg    rcx
  sub    rcx, 1
  mov    rax, rcx
  pop    rcx
  pop    rdi
  ret

; prints from rdi 
print:
  call   strlen
  mov    rdx, rax
  mov    rsi, rdi
  mov    rdi, 1
  mov    rax, 1
  syscall
  ret

; reads into rsi
read:
  mov    rdi, 64
  mov    rdx, rdi
  mov    rdi, 0
  mov    rax, 0
  syscall
  ret
