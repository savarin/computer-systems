section .text
global pangram
pangram:
		mov		rax, 0				; set up rax as running sum
		jmp .loop

.loop	movzx	rcx, byte [rdi]		; read byte from rdi to rcx
		cmp		rcx, 0				; end if null
		je .final
		add		rdi, 1				; increase index to set up next byte selection
		or		rcx, 0x20			; converts uppercase to lowercase and retains lowercase
		cmp		rcx, 0x61			; compare against 'A'
		jl .loop					; if less, next character
		cmp		rcx, 0x7A			; compare against 'Z'
		jg .loop					; if greater, next character
		sub		rcx, 0x61			; subtract 'A' to obtain exponent n
		mov		rbx, 1				; initialize number for 2^n
		jmp .shift

.shift	cmp		rcx, 0				; loop to shift left to obtain 2^n
		je .or
		shl		rbx, 1
		dec		rcx
		jmp .shift

.or		or		rax, rbx			; binary or to running sum
		jmp .loop

.final	cmp		rax, 67108863		; pangram returns 2^0 + 2^1 + ... + 2^25
		je .true
		mov		rax, 0
		ret

.true	mov		rax, 1
		ret
