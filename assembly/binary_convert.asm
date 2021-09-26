section .text
global binary_convert
binary_convert:
		xor		rax, rax			; set up rax as running sum
		mov		rbx, 512			; initialize rbx as large exponent of 2
		jmp .loop

.loop	shr		rbx, 1				; divide rbx by 2
		movzx	rcx, byte [rdi]		; read byte from rdi to rcx
		cmp		rcx, 0				; end if null
		je		.final
		add		rdi, 1				; increase index to set up next byte selection
		cmp		rcx, 49				; compare to '1'
		je .add
		cmp		rcx, 48				; compare to '0'
		je .loop
		xor		rax, rax			; if neither '0' or '1', reset and end
		jmp .end

.add	add		rax, rbx
		jmp	.loop

.final 	cmp		rbx, 0				; check if shift right has reduced rbx to 0
		je .end
		shr		rbx, 1				; shift right rbx to account for initialization
		shr		rax, 1				; shift right rax in tandem with rbx
		jmp .final

.end	ret
