.PHONY: tidy

# tidy_all: запустить go mod tidy во всех сервисах в ./apps
tidy:
ifeq ($(OS),Windows_NT)
	@echo Running 'go mod tidy' for all services in ./apps
	@for /d %%d in (apps\*) do @if exist %%d\go.mod ( \
		echo "Tidying apps/%%~nxd" & \
		pushd %%d & go mod tidy & popd \
	)
	@for /d %%d in (pkg\*) do @if exist %%d\go.mod ( \
		echo Tidying apps/%%~nxd & \
		pushd %%d & go mod tidy & popd \
	)
else
	@echo "Running 'go mod tidy' for all services in ./apps"
	@find ./apps -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "Tidying apps/$$dir"; \
		(cd "$$dir" && go mod tidy); \
	done
	@find ./pkg -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "Tidying apps/$$dir"; \
		(cd "$$dir" && go mod tidy); \
	done
endif

