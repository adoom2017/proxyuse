package sftpproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
)

type RemoteCopyInfo struct {
	User       string
	Password   string
	Target     string
	Proxy      string
	sftpClient *sftp.Client
}

func (info *RemoteCopyInfo) SshConnect() error {

	dialer, err := proxy.SOCKS5("tcp", info.Proxy, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		return err
	}

	hostconn, err := dialer.Dial("tcp", info.Target)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to dial to host:", err)
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: info.User,
		Auth: []ssh.AuthMethod{ssh.Password(info.Password)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	ncc, chans, reqs, err := ssh.NewClientConn(hostconn, info.Target, sshConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the ssh host:", err)
		return err
	}

	client := ssh.NewClient(ncc, chans, reqs)

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}

	info.sftpClient = sftpClient

	return nil
}

func (info *RemoteCopyInfo) UploadFile(localFile string, remotePath string) error {

	srcFile, err := os.Open(localFile)
	if err != nil {
		fmt.Println("os.Open error : ", localFile)
		log.Fatal(err)
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(localFile)
	dstFile, err := info.sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)
	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFile)
		log.Fatal(err)
	}
	dstFile.Write(ff)
	fmt.Println(localFile + " copy file to remote server finished!")

	return nil
}

func (info *RemoteCopyInfo) UploadDirectory(localPath string, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("read dir list fail ", err)
	}
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			info.sftpClient.Mkdir(remoteFilePath)
			info.UploadDirectory(localFilePath, remoteFilePath)
		} else {
			info.UploadFile(path.Join(localPath, backupDir.Name()), remotePath)
		}
	}
	fmt.Println(localPath + " copy directory to remote server finished!")
}
