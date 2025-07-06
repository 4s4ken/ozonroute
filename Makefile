ifneq ("$(wildcard tutor_makefile.mk)","")
include tutor_makefile.mk
endif
# пожалуйста, не удаляйте и не перемещайте этот импорт, он помогает вашему верному тьютору быстрее смотреть ваше дз
# вы можете описать ваш собственный makefile ниже

include .env


.PHONY: bin-deps .vendor-proto/validate .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options generate run \
        update linter start up goose-up\
        coverage

BINARY_NAME=./bin/homework
LOCAL_BIN := $(CURDIR)/bin
OUT_PATH := ./pkg

update:
	go mod tidy

linter:
	golangci-lint run

bin-deps: export GOBIN := $(LOCAL_BIN)
bin-deps: export PROTOC_VERSION := protoc-31.1-win64
bin-deps:
	curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v31.1/$(PROTOC_VERSION).zip
	unzip -o $(PROTOC_VERSION).zip -d $(LOCAL_BIN)
	rm $(PROTOC_VERSION).zip

	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest


.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
	https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
	cd vendor.protogen/tmp && \
	git sparse-checkout set --no-cone validate && \
	git checkout
	mkdir -p vendor.protogen/validate
	mv vendor.protogen/tmp/validate vendor.protogen/
	rm -rf vendor.protogen/tmp

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

generate: export GOBIN := $(LOCAL_BIN)
generate:
	mkdir -p $(OUT_PATH)
	$(LOCAL_BIN)/bin/protoc --proto_path=api --proto_path=vendor.protogen \
		--go_out=$(OUT_PATH) --go_opt=paths=source_relative --plugin protoc-gen-go="${GOBIN}/protoc-gen-go.exe" \
        --go-grpc_out=$(OUT_PATH) --go-grpc_opt=paths=source_relative --plugin protoc-gen-go-grpc="${GOBIN}/protoc-gen-go-grpc.exe" \
        --grpc-gateway_out=$(OUT_PATH) --grpc-gateway_opt=paths=source_relative --plugin protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway.exe \
        --openapiv2_out=$(OUT_PATH) --plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2.exe \
		api/pvz/pvz.proto
	go mod tidy

run: bin-deps .vendor-proto/validate .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options generate

DATABASE_DSN=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

goose-add:
	goose -dir ./migrations postgres "$(DATABASE_DSN)" create $(NAME) sql

goose-up:
	goose -dir ./migrations postgres "$(DATABASE_DSN)" up

goose-down:
	goose -dir ./migrations postgres "$(DATABASE_DSN)" down

goose-status:
	goose -dir ./migrations postgres "$(DATABASE_DSN)" status

test-unit:
	go test -v -count=1 ./internal/app/... ./internal/domain/... -tags=unit

test-integration:
	go test -v -count=1 ./internal/storage/dbsuite -tags=integration

test-e2e:
	go build -o internal/e2etests/server.exe ./cmd/grpc/server
	go test -v -timeout=5m ./internal/e2etests

test: test-unit test-integration test-e2e

coverage:
	go test -coverpkg=gitlab.ozon.dev/gleb_evseev/homework/internal/app,gitlab.ozon.dev/gleb_evseev/homework/internal/storage/db \
		-coverprofile=coverage.txt \
		./internal/app \
		./internal/storage/dbsuite
	go tool cover -func=coverage.txt
	go tool cover -html=coverage.txt -o coverage.html

.PHONY: run-server run-cli

run-cli:
	go run ./cmd/cli/main.go

run-server:
	go run ./cmd/grpc/server/main.go

run-solo-client:
	go run ./cmd/grpc/client/solo_req/main.go

run-par-client:
	go run ./cmd/grpc/client/paralel_req/main.go



