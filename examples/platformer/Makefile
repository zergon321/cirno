run:
	go run beholder.go bullets.go platform.go player.go main.go --vsync --draw-wires

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/demo.exe .

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/demo .