# build go api program
Set-Location -Path "..\api"

# get os type and arch
$os = & go env GOOS
$arch = & go env GOARCH

Write-Host "CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o bin\chatgpt-plus-api main.go"

$env:CGO_ENABLED="0"
$env:GOOS=$os
$env:GOARCH=$arch
go build -o bin\chatgpt-plus-api.exe main.go