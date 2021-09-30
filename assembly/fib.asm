section .text
global fib
fib:	mov		rax, rdi	; base case
		cmp		rdi, 1		; return 1 if n <= 1
		jle .final

		push	rbx			; push rbx to stack
		push	rcx			; push rcx to stack

		mov		rbx, rdi	; store n in rbx
		sub		rdi, 1
		call fib			; call fib(n-1)
		mov		rcx, rax	; store result in rcx

		mov		rdi, rbx	; restore n
		sub		rdi, 2
		call fib			; call fib(n-2)
		add		rax, rcx	; fib(n-1) + fib(n-2)

		pop		rcx			; pop rcx off stack
		pop		rbx			; pop rcb off stack

.final: ret
