
clean:
	@echo clean
	rm -rf build/*

build: clean
	@echo build
	cd tools/transformers && \
	go build -o ../../build/ . ;

