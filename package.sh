#/bin/sh

mkdir -p out
while read -r os arch; do
    env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -o ./npc cmd/npc/npc.go
    out_name=out/${os}_${arch}_client.tar.gz
    tar -zcf ${out_name} npc -C conf npc.conf
    echo ${out_name} ok

    env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -o ./nps cmd/nps/nps.go
    out_name=out/${os}_${arch}_serverr.tar.gz
    tar -zcf ${out_name} nps conf
    echo ${out_name} ok
    rm -f npc nps
done << EOF
darwin amd64
windows 386
windows amd64
freebsd 386
freebsd amd64
freebsd arm
linux 386
linux amd64
linux arm64
linux arm
linux mips64le
linux mips64
linux mipsle
linux mips
EOF
