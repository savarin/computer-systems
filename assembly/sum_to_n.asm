section .text
global sum_to_n
sum_to_n:	
		mov		rax, 0		; set rax as running sum
		jmp .loop

.loop	test	rdi, rdi	; line up rdi for jump if zero
		jz .end
		add 	rax, rdi	; add rdi to rax
		dec		rdi			; decrement rdi
		jmp .loop

.end	ret
