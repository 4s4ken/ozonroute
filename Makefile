ifneq ("$(wildcard tutor_makefile.mk)","")
include tutor_makefile.mk
endif
# пожалуйста, не удаляйте и не перемещайте этот импорт, он помогает вашему верному тьютору быстрее смотреть ваше дз
# вы можете описать ваш собственный makefile ниже

BINARY_NAME=homework

.PHONY: update linter build start run

update:
	go mod tidy

linter:
	golangci-lint run

build:
	go build -o $(BINARY_NAME) ./cmd/cli

start:
	./$(BINARY_NAME)

run: update linter build start
