default: watch

watch:
	watchexec --restart -w src -w in "go run src/main.go && go fmt out/main.go && go run out/main.go"
