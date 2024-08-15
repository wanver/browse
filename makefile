#!make
ifeq ($(case),)
test:
	@echo "Skipping test execution as case is empty."
else
test:
	go test ./... -v -run "$(case)" -timeout 0
endif

