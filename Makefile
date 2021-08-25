setup:
	@echo "--- Setup and generate config yaml files ---"
	@cp config.sample.yaml config.yaml

httpd:
	@echo "--- Run httpd server ---"
	@go run cmd/httpd/main.go

help:
	@echo "make setup: setup and generate config yaml files"
	@echo "make httpd: run httpd server"
