package main

import (
	"fmt"
	"os"
	"proxyuse/sftpproxy"
)

func main() {
	copyInfo := &sftpproxy.RemoteCopyInfo{
		User:     "xxxx",
		Password: "xxxxxxxxxxxxx",
		Target:   "xxxxxxxxxxxxxxxx",
		Proxy:    "127.0.0.1:7890",
	}

	err := copyInfo.SshConnect()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create ssh client:", err)
		return
	}

	copyInfo.UploadDirectory("C:\\Users\\shendongchun\\npp", "/home/adoom/npp")
}
