GO_MOD=go mod

.PHONY: tidy
tidy:
	$(GO_MOD) tidy
.PHONY: download
download:
	$(GO_MOD) download
