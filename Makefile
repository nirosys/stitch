
VERSION=$(shell git describe --always --tags --long --dirty | sed -e 's/-/./' -e 's/-g/-/')

all: test
	@if test -n "${GITHUB_ACTIONS}"; then \
		echo "##[set-output name=version;]${VERSION}"; \
	fi
	@echo "Building: $(VERSION)"
	@CGO_ENABLED=0 go build -o ./stitch -ldflags \
		"-X 'github.com/nirosys/stitch/cmd/stitch/subcmd/repl.ReplVersion=$(VERSION)'" \
		./cmd/stitch

test:
	@echo "Testing.."
	@CGO_ENABLED=0 go test ./...
