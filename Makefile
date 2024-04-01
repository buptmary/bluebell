# 虚拟目标：确保即使存在同名文件，目标也会被执行
.PHONY: all build run gotool clean help

# 设置可执行文件的名称为"bluebell"
BINARY="bluebell"

# 默认目标，当直接运行make时，会先执行gotool，然后执行build
all: gotool build

build:
	@echo "build..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}

run:
	@echo "run..."
	@go run ./main.go conf/config.yaml

gotool:
	go fmt ./
	go vet ./

# 清理目标，如果二进制文件存在，则移除它
clean:
	@if [ -f ${BINARY} ]; then rm ${BINARY}; fi

# 帮助目标，显示可用的make命令及其说明
help:
	@echo "make - Format Go code and compile to generate binary file."
	@echo "make build - Compile Go code and generate binary file."
	@echo "make run - Compile and run Go code."
	@echo "make gotool - Format Go code."
	@echo "make clean - Clean up the generated binary file."