build:
	go build -ldflags="-s -w -H=windowsgui" cmd/npc_windows/npc.go

release:
	xgo --targets=windows/386,windows/amd64 -ldflags="-s -w -H=windowsgui" cmd/npc_windows/npc.go

