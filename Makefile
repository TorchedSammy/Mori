PREFIX ?= /usr
DESTDIR ?=
BINDIR ?= $(PREFIX)/bin

build:
	@go build -ldflags "-s -w"

install:
	@install -v -d "$(DESTDIR)$(BINDIR)/" && install -m 0755 -v mori "$(DESTDIR)$(BINDIR)/mori"
	@echo "Mori has been installed"

uninstall:
	@rm -vrf "$(DESTDIR)$(BINDIR)/mori" 
	@echo "Mori has been removed from the system bin directory! :("

clean:
	@go clean

all: build install

.PHONY: install uninstall build clean
