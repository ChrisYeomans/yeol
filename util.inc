macro read fd, buff, count {
    mov rax, 0
    mov rdi, fd
    mov rsi, buff
    mov rdx, count
    syscall
}

macro exit_program {
  xor    rdi, rdi ; exit code 0
  mov    rax, 60 ; sys_exit
  syscall
}

macro write fd, buf, count {
   mov rax, 1
   mov rdi, fd
   mov rsi, buf
   mov rdx, count
   syscall 
}

;; Parse unsigned integer from a sized string
;;   rdi - void *buf
;;   rsi - size_t n
parse_uint:
    xor rax, rax
    xor rbx, rbx
    mov rcx, 10
.next_digit:
    cmp rsi, 0
    jle .done

    mov bl, byte [rdi]
    cmp rbx, '0'
    jl .done
    cmp rbx, '9'
    jg .done
    sub rbx, '0'

    mul rcx
    add rax, rbx

    inc rdi
    dec rsi
    jmp .next_digit
.done:
    ret


;; Write an integer to a file
;;   rdi - int fd
;;   rsi - uint64_t int x
write_uint:
    test rsi, rsi
    jz .base_zero

    mov rcx, 10     ;; 10 literal for division
    mov rax, rsi    ;; keeping track of rsi in rax cause it's easier to div it like that
    mov r10, 0      ;; counter of how many digits we already converted
.next_digit:
    test rax, rax
    jz .done
    mov rdx, 0
    div rcx
    add rdx, '0'
    dec rsp
    mov byte [rsp], dl
    inc r10
    jmp .next_digit
.done:
    write rdi, rsp, r10
    add rsp, r10
    ret
.base_zero:
    dec rsp
    mov byte [rsp], '0'
    write rdi, rsp, 1
    inc rsp
    ret
