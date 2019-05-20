$version="0.0.1"
& { $env:GOARCH="amd64"; $env:GOOS="linux"; go build -o bin/bom-$version.linux-amd64 main.go }
& { $env:GOARCH="arm"; $env:GOOS="linux"; go build -o bin/bom-$version.linux-arm main.go }
& { $env:GOARCH="amd64"; $env:GOOS="windows"; go build -o bin/bom-$version.windows-amd64.exe main.go }
& { $env:GOARCH="amd64"; $env:GOOS="freebsd"; go build -o bin/bom-$version.freebsd-amd64 main.go }
& { $env:GOARCH="arm"; $env:GOOS="freebsd"; go build -o bin/bom-$version.freebsd-arm main.go }

Remove-Item env:GOOS
Remove-Item env:GOARCH