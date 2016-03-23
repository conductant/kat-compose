all: test-aurora

clean:
	-rm .thrift-*

SCHEDULER_HOSTPORT?="192.168.33.7:8081"

test-aurora:
	SCHEDULER_HOSTPORT=$(SCHEDULER_HOSTPORT) \
	${GODEP} go test . -logtostderr -check.vv -v ${TEST_ARGS}


install-thrift:
ifeq ($(wildcard ./has-thrift),)
	brew install bison
	brew install boost
	brew install thrift
	go get git.apache.org/thrift.git/lib/go/thrift
endif

API_DIR:=github.com/conductant/kat-compose/pkg/aurora
TARGET:=$(GOPATH)/src/$(API_DIR)
PACKAGE_PREFIX:=$(API_DIR)/
gen-thrift: install-thrift
ifeq ($(wildcard ./gen-thrift),)
	-rm -rf $(TARGET)/api
	-mkdir -p $(TARGET)/api
	thrift -r --gen go:package=api,package_prefix=$(PACKAGE_PREFIX) --out $(TARGET) ./api.thrift
endif

