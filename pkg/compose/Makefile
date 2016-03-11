all: test-compose

test-compose:
	${GODEP} go test ./...  -logtostderr -check.vv -v ${TEST_ARGS}
