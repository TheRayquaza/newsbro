.PHONY: help build test run docker-build docker-up docker-down clean fmt lint

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build:
	cargo build --release

test:
	cargo test

run:
	cargo run

clean:
	cargo clean
	docker-compose down -v

fmt:
	cargo fmt

lint:
	cargo clippy -- -D warnings

check: fmt lint test

dev:
	RUST_LOG=debug cargo watch -x run
