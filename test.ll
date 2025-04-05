@tmp = global [4 x i8] c"%d\0A\00"

declare i32 @printf(i8* %format, ...)

define i32 @main() {
0:
	%1 = alloca i32
	store i32 12, i32* %1
	%2 = getelementptr [3 x i8], [4 x i8]* @tmp, i32 0, i32 0
	%3 = call i32 (i8*, ...) @printf(i8* %2, i32 12, i32 0)
	ret i32 0
}
