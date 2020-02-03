.DEFAULT_GOAL := help

build: ## build
	go build -o bin/hoge -tags=jsoniter handler/main.go
run: ## dev run
	go run -tags=jsoniter handler/main.go


help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

