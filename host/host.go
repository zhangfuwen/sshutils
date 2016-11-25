package host

import (
	"github.com/zhangfuwen/sshutils/scp"
	"golang.org/x/crypto/ssh"
)

type Host struct {
	IP       string
	Port     string
	UserName string
	Password string
	client   *ssh.Client
}

func NewHost(ip, user, pass string) Host {
	return Host{
		ip, "22", user, pass, nil,
	}
}

func NewHostWithPort(ip, port, user, pass string) Host {
	return Host {
		ip,port,user,pass,nil,
	}
}

func (h *Host) dial() error {
	var err error
	h.client, err = ssh.Dial("tcp", h.IP+":"+h.Port, &ssh.ClientConfig{
		User: h.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
	})
	return err
}

func (h *Host) Get(remotePath, localDir string) error {
	if h.client == nil {
		if err := h.dial(); err != nil {
			return err
		}
	}
	return scp.CopyRemoteFileToLocalPath(remotePath, localDir, h.client)
}

func (h *Host) Put(localPath, remoteDir string) error {
	if h.client == nil {
		if err := h.dial(); err != nil {
			return err
		}
	}
	return scp.CopyLocalFileToRemotePath(localPath, remoteDir, h.client)
}

func (h *Host) Run(cmd string) error {
	if h.client == nil {
		if err := h.dial(); err != nil {
			return err
		}
	}
	session, err := h.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return session.Run(cmd)
}

func (h *Host) Output(cmd string) (string, error) {
	if h.client == nil {
		if err := h.dial(); err != nil {
			return "", err
		}
	}
	session, err := h.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	bs, err := session.Output(cmd)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
func (h *Host) Send(cmd string, expect string) {

}
func (h *Host) Disconnect() {
	h.client.Close()
	h.client = nil
}
