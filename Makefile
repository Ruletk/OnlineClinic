.PHONY: tidy proto-auth clean-proto-auth


tidy:
ifeq ($(OS),Windows_NT)
	@echo Running 'go mod tidy' for all services in ./apps
	@for /d %%d in (apps\*) do @if exist %%d\go.mod ( \
		echo Tidying apps/%%~nxd & \
		pushd %%d & go mod tidy & popd \
	)
	@for /d %%d in (pkg\*) do @if exist %%d\go.mod ( \
		echo Tidying pkg/%%~nxd & \
		pushd %%d & go mod tidy & popd \
	)
else
	@echo "Running 'go mod tidy' for all services in ./apps and ./pkg"
	@find ./apps -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "Tidying $$dir"; \
		(cd "$$dir" && go mod tidy); \
	done
	@find ./pkg -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "Tidying $$dir"; \
		(cd "$$dir" && go mod tidy); \
	done
endif

proto-auth:
ifeq ($(OS),Windows_NT)
	@echo Generating Go gRPC files for auth...
	@if not exist pkg\proto\gen\auth mkdir pkg\proto\gen\auth
	protoc --go_out=pkg/proto/gen/auth --go-grpc_out=pkg/proto/gen/auth --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative -I pkg/proto pkg/proto/auth/*.proto
	@echo Auth proto generation complete.
else
	@echo "Generating Go gRPC files for auth..."
	@mkdir -p pkg/proto/gen/auth
	protoc --go_out=pkg/proto/gen/auth --go-grpc_out=pkg/proto/gen/auth \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		-I pkg/proto pkg/proto/auth/*.proto
	@echo "Auth proto generation complete."
endif

clean-proto-auth:
	@echo "Cleaning generated files for auth..."
	@rm -rf pkg/proto/gen/auth



