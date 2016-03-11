all: clean kat-compose

clean:
	-rm kat-compose

# Common makefile that extracts git version info and generates the LDFLAGS variable.
include ../hack/make/version.mk

kat-compose:
	${GODEP} go build -ldflags "$(LDFLAGS)" kat-compose.go

