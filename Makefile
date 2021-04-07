targets=$(shell ls cmds | sed 's:^:./bin/:')
outputs=$(shell ls cmds | sed 's:^:./cmds/:')
gosrc=$(shell find ./ -iname '*.go')

$(targets): $(gosrc)
	mkdir -p ./bin/
	go build -o ./bin/ $(outputs)

.PHONY: test
test:
	go clean -testcache
	go test -v ./...

.PHONY: clean
clean:
	rm -r ./bin/
