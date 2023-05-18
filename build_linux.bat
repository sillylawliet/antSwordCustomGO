go env -w GOARCH=amd64
go env -w GOARM=7
go env -w GOOS=linux
go env -w CGO_ENABLED=0
go build -ldflags "-s -w" -o gos