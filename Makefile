TEST_CMD=$(shell type gotestsum 2>/dev/null >/dev/null && echo "gotestsum --" || echo "go test")

mocks:
	mockgen -source=app/service/service.go -destination=app/service/mocks/mocks.go -package=mocks Repository Engine Config WAL Storage Segment
	mockgen -source=app/storage/wal/log_manager.go -destination=app/storage/wal/log_manager_mocks.go -package=wal Segment
	mockgen -source=app/storage/storage.go -destination=app/storage/storage_mocks.go -package=storage Engine WAL
test:
	$(TEST_CMD) -tags=testing ./...