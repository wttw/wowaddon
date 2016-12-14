VERSION = `./wowaddon releasetag`
DASHVERSION = `./wowaddon releasetag | tr . -`

.PHONY: clean
clean:
	rm -f *~

.PHONY: distclean
distclean: clean
	rm -f wowaddon wowaddon.exe linux/wowaddon
	-rmdir linux

.PHONY: install
install:
	go install

wowaddon: *.go
	go build

wowaddon.exe: *.go
	GOOS=windows go build

linux/wowaddon: *.go
	mkdir -p linux
	GOOS=linux go build -o linux/wowaddon

README.md: wowaddon mkreadme
	perl mkreadme

.PHONY: package
package: wowaddon wowaddon.exe linux/wowaddon
	rm -f wowaddon-osx-$(DASHVERSION).zip wowaddon-windows-$(DASHVERSION).zip
	zip -9 wowaddon-osx-$(DASHVERSION).zip wowaddon
	zip -9 wowaddon-windows-$(DASHVERSION).zip wowaddon.exe
	cd linux && zip -9 ../wowaddon-linux-`../wowaddon releasetag | tr . -`.zip wowaddon

.PHONY: catalog
catalog: wowaddon
	rm -f addoncatalog.json.zip
	cd ../wowslurp && go build && APP_VERSION=`../wowaddon/wowaddon releasetag` ./wowslurp
	cd ../wowslurp && zip -9 ../wowaddon/addoncatalog.json.zip addoncatalog.json

.PHONY: ship
ship: package catalog
	./shipit


