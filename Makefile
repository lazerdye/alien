
PACKAGE=github.com/lazerdye/alien

gofmt:
	go fmt $(PACKAGE)/...

test:
	go test -v $(PACKAGE)/...
