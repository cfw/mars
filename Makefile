GO_MOD=go mod

.PHONY: download
download:
	$(GO_MOD) tidy
