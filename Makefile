
VERSION=$(shell git describe --always --tags --long --dirty | sed -e 's/-/./' -e 's/-g/-/')

all: test
	@echo "Building: $(VERSION)"
	@go build -o ./stitch -ldflags \
		"-X 'github.com/nirosys/stitch/cmd/stitch/subcmd/repl.ReplVersion=$(VERSION)'" \
		./cmd/stitch

test:
	@echo "Testing.."
	@go test ./...
