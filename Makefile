TEST_CMD=$(shell type gotestsum 2>/dev/null >/dev/null && echo "gotestsum --" || echo "go test")

mocks:
	mockgen -source=app/service/service.go -destination=app/service/mocks/mocks.go -package=mocks Repository Engine Config

test:
	$(TEST_CMD) -tags=testing ./...