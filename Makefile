build:
	env CGO_ENABLED=0 GOOS=linux go build -o bin/param-cache cmd/main.go

travis:
	env CGO_ENABLED=0 GOOS=linux go build -o bin/param-cache cmd/main.go
	cd bin && tar -czf param-cache-linux.${TRAVIS_TAG}.tar.gz ./param-cache && rm param-cache

	env CGO_ENABLED=0 GOOS=darwin go build -o bin/param-cache cmd/main.go
	cd bin && tar -czf param-cache-darwin.${TRAVIS_TAG}.tar.gz ./param-cache && rm param-cache

	env CGO_ENABLED=0 GOOS=windows go build -o bin/param-cache.exe cmd/main.go
	cd bin && zip -9 param-cache-windows.${TRAVIS_TAG}.zip ./param-cache.exe && rm param-cache.exe
