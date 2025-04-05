
test:
	go run main test.yeol
	clang test.ll -o test	

build:
	go build main -o yeol


clean:
	rm -f test.asm
	rm -f *.o
	rm -f *.out
	rm -f test
	rm -f test.ll

