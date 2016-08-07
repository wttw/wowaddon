VERSION = `./wowaddon releasetag`
DASHVERSION = `./wowaddon releasetag | tr . -`

.PHONY: clean
clean:
	rm -f *~

.PHONY: distclean
distclean: clean
	rm -f wowaddon wowaddon.exe

.PHONY: install
install:
	go install

wowaddon: *.go
	go build

wowaddon.exe: *.go
	GOOS=windows go build

.PHONY: package
package: wowaddon wowaddon.exe
	rm -f wowaddon-osx-$(DASHVERSION).zip wowaddon-windows-$(DASHVERSION).zip
	zip -9 wowaddon-osx-$(DASHVERSION).zip wowaddon
	zip -9 wowaddon-windows-$(DASHVERSION).zip wowaddon.exe

.PHONY: catalog
catalog:
	cd ../wowslurp && go build && ./wowslurp
	cp ../wowslurp/addoncatalog.json.zip .

.PHONY: ship
ship: package catalog
	./shipit


