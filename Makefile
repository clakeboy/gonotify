NAME = gonotify
ARCH = amd64
OS = linux
#linux
all:
	go generate && CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=$(OS)  go build -x -v -ldflags "-w" -o $(NAME) main.go
	upx -9 $(NAME)
.PHONY : clean
clean:
	rm -f $(NAME)