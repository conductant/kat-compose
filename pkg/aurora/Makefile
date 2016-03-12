all: test-aurora

SCHEDULER_HOSTPORT?="192.168.33.7:8081"

test-aurora: gen-thrift
	${GODEP} go test ./... -check.vv -v ${TEST_ARGS}

install-thrift:
	brew install bison
	brew install boost
	brew install thrift
	go get git.apache.org/thrift.git/lib/go/thrift

API_DIR:=org/apache/aurora
_gen-thrift: install-thrift
	-rm -rf ../../vendor/$(API_DIR)
	-mkdir -p ../../vendor/$(API_DIR)
	thrift -r --gen go:package=api,package_prefix=org.apache.aurora. --out ../../vendor/$(API_DIR) ./api.thrift

PACKAGE_PREFIX:=github.com/conductant/kat-compose/pkg/aurora/
gen-thrift: install-thrift
	-rm -rf ./api
	-mkdir -p ./api
	thrift -r --gen go:package=api,package_prefix=$(PACKAGE_PREFIX) --out . ./api.thrift

