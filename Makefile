SHELL=cmd.exe
AUTH_APP=auth_app
DISCORD_BOT=discord_bot

build: build_auth build_discord

build_auth:
	@echo Building ${AUTH_APP} linux binary...
	chdir ./cmd/auth && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ../../build/${AUTH_APP} ./
	@echo Done!

build_discord:
	@echo Building ${DISCORD_BOT} linux binary...
	chdir ./cmd/discord && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ../../build/${DISCORD_BOT} ./
	@echo Done!