main:
    mov rax, 10

for:
    add rbx, 3
    dec rax
    cmp rax, 0
    jz fim
    jmp for

fim:
