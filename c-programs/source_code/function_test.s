	.file	"function_test.c"
	.text
	.section	.rodata
.LC0:
	.string	"The sqrt of %d is %f\n"
	.text
	.globl	main
	.type	main, @function
main:
.LFB6:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	subq	$16, %rsp
	movl	$823809, -4(%rbp)
	movl	-4(%rbp), %eax
	movl	%eax, %edi
	call	calcSqrt
	movd	%xmm0, %eax
	movl	%eax, -8(%rbp)
	cvtss2sd	-8(%rbp), %xmm0
	movl	-4(%rbp), %eax
	movl	%eax, %esi
	movl	$.LC0, %edi
	movl	$1, %eax
	call	printf
	movl	$0, %eax
	leave
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE6:
	.size	main, .-main
	.globl	findGreater
	.type	findGreater, @function
findGreater:
.LFB7:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	movl	%edi, -4(%rbp)
	movl	%esi, -8(%rbp)
	movl	-4(%rbp), %eax
	cmpl	-8(%rbp), %eax
	jle	.L4
	movl	-8(%rbp), %eax
	jmp	.L5
.L4:
	movl	-4(%rbp), %eax
.L5:
	popq	%rbp
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE7:
	.size	findGreater, .-findGreater
	.section	.rodata
	.align 8
.LC1:
	.string	"both integers are equal, their GCD will always be themselves"
	.text
	.globl	GCD
	.type	GCD, @function
GCD:
.LFB8:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	subq	$32, %rsp
	movl	%edi, -20(%rbp)
	movl	%esi, -24(%rbp)
	movl	-20(%rbp), %eax
	cmpl	-24(%rbp), %eax
	jle	.L7
	movl	-20(%rbp), %eax
	movl	%eax, -4(%rbp)
	movl	-24(%rbp), %eax
	movl	%eax, -8(%rbp)
	jmp	.L12
.L7:
	movl	-24(%rbp), %eax
	cmpl	-20(%rbp), %eax
	jle	.L9
	movl	-24(%rbp), %eax
	movl	%eax, -4(%rbp)
	movl	-20(%rbp), %eax
	movl	%eax, -8(%rbp)
	jmp	.L12
.L9:
	movl	$.LC1, %edi
	call	puts
	movl	$0, %eax
	jmp	.L10
.L12:
	movl	-4(%rbp), %eax
	movl	%eax, -12(%rbp)
	movl	-4(%rbp), %eax
	cltd
	idivl	-8(%rbp)
	movl	%edx, %eax
	testl	%eax, %eax
	jne	.L11
	movl	-8(%rbp), %eax
	jmp	.L10
.L11:
	movl	-8(%rbp), %eax
	movl	%eax, -4(%rbp)
	movl	-12(%rbp), %eax
	cltd
	idivl	-8(%rbp)
	movl	%edx, -8(%rbp)
	jmp	.L12
.L10:
	leave
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE8:
	.size	GCD, .-GCD
	.globl	calcAbs
	.type	calcAbs, @function
calcAbs:
.LFB9:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	movss	%xmm0, -4(%rbp)
	pxor	%xmm0, %xmm0
	comiss	-4(%rbp), %xmm0
	jbe	.L18
	cvtss2sd	-4(%rbp), %xmm0
	cvtsd2ss	%xmm0, %xmm0
	movss	.LC3(%rip), %xmm1
	xorps	%xmm1, %xmm0
	jmp	.L16
.L18:
	movss	-4(%rbp), %xmm0
.L16:
	popq	%rbp
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE9:
	.size	calcAbs, .-calcAbs
	.globl	calcSqrt
	.type	calcSqrt, @function
calcSqrt:
.LFB10:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	subq	$48, %rsp
	movl	%edi, -36(%rbp)
	movl	$5, -16(%rbp)
	movl	-36(%rbp), %eax
	movl	%eax, -20(%rbp)
	movl	-36(%rbp), %eax
	movl	%eax, %edi
	call	largestWholeSquare
	movl	%eax, -24(%rbp)
	cvtsi2ss	-24(%rbp), %xmm0
	movss	-4(%rbp), %xmm1
	addss	%xmm1, %xmm0
	movss	%xmm0, -4(%rbp)
	movl	$0, -8(%rbp)
	jmp	.L20
.L27:
	movl	-8(%rbp), %eax
	addl	$1, %eax
	cvtsi2sd	%eax, %xmm1
	movsd	.LC4(%rip), %xmm0
	call	pow
	movapd	%xmm0, %xmm1
	movsd	.LC5(%rip), %xmm0
	divsd	%xmm1, %xmm0
	cvtsd2ss	%xmm0, %xmm2
	movss	%xmm2, -28(%rbp)
	movl	$0, -12(%rbp)
	jmp	.L21
.L26:
	movss	-4(%rbp), %xmm0
	addss	-28(%rbp), %xmm0
	cvtss2sd	%xmm0, %xmm0
	movsd	.LC6(%rip), %xmm1
	call	pow
	movapd	%xmm0, %xmm1
	cvtsi2sd	-36(%rbp), %xmm0
	comisd	%xmm1, %xmm0
	jnb	.L29
	jmp	.L25
.L29:
	movss	-4(%rbp), %xmm0
	addss	-28(%rbp), %xmm0
	movss	%xmm0, -4(%rbp)
	addl	$1, -12(%rbp)
.L21:
	cmpl	$9, -12(%rbp)
	jle	.L26
.L25:
	addl	$1, -8(%rbp)
.L20:
	movl	-8(%rbp), %eax
	cmpl	-16(%rbp), %eax
	jle	.L27
	movss	-4(%rbp), %xmm0
	leave
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE10:
	.size	calcSqrt, .-calcSqrt
	.globl	largestWholeSquare
	.type	largestWholeSquare, @function
largestWholeSquare:
.LFB11:
	.cfi_startproc
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register 6
	movl	%edi, -36(%rbp)
	movl	$0, -4(%rbp)
	movl	-36(%rbp), %eax
	movl	%eax, -8(%rbp)
	jmp	.L31
.L35:
	movl	-8(%rbp), %edx
	movl	-4(%rbp), %eax
	addl	%edx, %eax
	movl	%eax, %edx
	shrl	$31, %edx
	addl	%edx, %eax
	sarl	%eax
	movl	%eax, -16(%rbp)
	movl	-16(%rbp), %eax
	imull	-16(%rbp), %eax
	cltq
	movq	%rax, -24(%rbp)
	movl	-36(%rbp), %eax
	cltq
	cmpq	%rax, -24(%rbp)
	jne	.L32
	movq	-24(%rbp), %rax
	jmp	.L33
.L32:
	movl	-36(%rbp), %eax
	cltq
	cmpq	%rax, -24(%rbp)
	jbe	.L34
	movl	-16(%rbp), %eax
	subl	$1, %eax
	movl	%eax, -8(%rbp)
	jmp	.L31
.L34:
	movl	-16(%rbp), %eax
	movl	%eax, -12(%rbp)
	movl	-16(%rbp), %eax
	addl	$1, %eax
	movl	%eax, -4(%rbp)
.L31:
	movl	-4(%rbp), %eax
	cmpl	-8(%rbp), %eax
	jle	.L35
	movl	-12(%rbp), %eax
.L33:
	popq	%rbp
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE11:
	.size	largestWholeSquare, .-largestWholeSquare
	.section	.rodata
	.align 16
.LC3:
	.long	2147483648
	.long	0
	.long	0
	.long	0
	.align 8
.LC4:
	.long	0
	.long	1076101120
	.align 8
.LC5:
	.long	0
	.long	1072693248
	.align 8
.LC6:
	.long	0
	.long	1073741824
	.ident	"GCC: (GNU) 8.5.0 20210514 (Red Hat 8.5.0-18)"
	.section	.note.GNU-stack,"",@progbits
