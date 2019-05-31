build:
	go build cmd/npc_windows/npc.go

release:
	xgo --targets=windows/386,windows/amd64 -ldflags="-s -w -H=windowsgui" ./cmd/npc_windows
	mkdir -p out
	mv -f npc_windows-windows-4.0-386.exe npc.exe
	tar -zcf out/win_32_npc.tar.gz npc.exe
	mv -f npc_windows-windows-4.0-amd64.exe npc.exe
	tar -zcf out/win_64_npc.tar.gz npc.exe

