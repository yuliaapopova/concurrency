TEST_CMD=$(shell type gotestsum 2>/dev/null >/dev/null && echo "gotestsum --" || echo "go test")

mocks:
	mockgen -source=app/compute/parser.go -destination=app/compute/mocks/mocks.go -package=mocks Repository

test:
	$(TEST_CMD) -tags=testing ./...