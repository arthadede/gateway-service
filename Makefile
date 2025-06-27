project_name = histweety-gateway-service
image_name = histweety/gateway-service:latest

run-local: ## Run the app locally
	go run main.go

requirements: ## Generate go.mod & go.sum files
	go mod tidy

clean-packages: ## Clean packages
	go clean -modcache

build: ## Generate docker image
	docker build -t $(image_name) .

push: ## Push image
	docker push $(image_name)