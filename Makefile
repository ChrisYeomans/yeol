
test:
	go run main test.yeol
	nasm -felf64 -g test.asm
	ld -o test test.o

build:
	go build main -o yeol

yeol:
	./yeol test.yeol
	nasm -felf64 -g test.asm
	ld -o test test.o

clean:
	rm -f test.asm
	rm -f *.o
	rm -f *.out
	rm -f test

