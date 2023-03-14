SHELL=cmd.exe
AUTH=auth
DISCORD=discord

up_build: build down
	@echo Building Docker images (if necessary) and Starting Docker containers...
	docker-compose up --build -d
	@echo Docker containers are now running!


up:
	@echo Starting Docker containers...
	docker-compose up -d
	@echo Docker containers are now running!

down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

build: build_auth build_discord

build_auth:
	@echo Building ${AUTH} linux binary...
	chdir ./cmd/${AUTH} && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ../../build/${AUTH} ./
	@echo Done!

build_discord:
	@echo Building ${DISCORD} linux binary...
	chdir ./cmd/${DISCORD} && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ../../build/${DISCORD} ./
	@echo Done!