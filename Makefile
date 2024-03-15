default: help

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@go build -o build/lynkly cmd/lynkly/main.go

run: build ## Run the application
	@./build/lynkly &

stop:
	@killall lynkly || true

restart: stop run ## Restart the application

fswatch-check:
	@command -v fswatch >/dev/null 2>&1 || { echo >&2 "You need fswatch. Run 'brew install fswatch'"; exit 1; }

watch: fswatch-check run ## Watch for changes and restart the application
	@echo "Watching *.go files for changes..."
	@fswatch -l 2 -e ".*" -i "\.go$$" -o . | xargs -n1 -I {} sh -c 'echo "Changes detected, restarting..."; make restart'


.PHONY: help build run stop restart watch
