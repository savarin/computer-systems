section .text
global binary_convert
binary_convert:
		mov		rax, 0				; set up rax as running sum
		mov		rbx, 512			; initialize rbx as large exponent of 2
		jmp .loop

.loop	shr		rbx, 1				; divide rbx by 2
		movzx	rcx, byte [rdi]		; read byte from rdi to rcx
		add		rdi, 1				; increase index to set up next byte selection
		cmp		rcx, 49				; compare to '1'
		je .add
		cmp		rcx, 48				; compare to '0'
		je .loop
		jmp .final

.add	add		rax, rbx
		jmp	.loop

.final 	cmp		rbx, 0				; shift right rbx to account for initialization
		je .end
		shr		rbx, 1
		shr		rax, 1				; shift right rax in tandem with rbx
		jmp .final

.end	ret