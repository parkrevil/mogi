# Development commands
.PHONY: dev-suction-server dev-suction-client

# Go services
dev-suction-server:
	go run apps/suction-server/main.go

dev-suction-client:
	go run apps/suction-client/main.go

# Utility commands
.PHONY: clean install

clean:
	@echo "Cleaning builds..."
	rm -rf apps/*/bin/

install:
	@echo "Installing dependencies..."
	go mod tidy
