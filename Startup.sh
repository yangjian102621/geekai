# 前端

if ! command -v node > /dev/null; then
	printf 'node is not installed.\n'
	exit 1
fi
cd web
npm install
npm run build
cd ..

# 后端
if ! command -v go > /dev/null; then
	printf 'go is not installed.\n'
	exit 1
fi
cd src
go mod tidy
make linux
cd ..

# Docker
if ! command -v docker > /dev/null; then
	printf 'docker is not installed.\n'
	exit 1
fi
docker compose up -d