# Development commands
.PHONY: dev-server dev-client

# Go services
dev-server:
	cd apps/suction-server && bunx nodemon

dev-client:
	cd apps/suction-client && bunx nodemon

# Utility commands
.PHONY: clean install

clean:
	@echo "Cleaning builds..."
	rm -rf apps/{suction-client,suction-server}/{bin,tmp}

install:
	@echo "Installing dependencies..."
	go work sync
