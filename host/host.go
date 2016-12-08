package host

import (
	"github.com/zhangfuwen/sshutils/scp"
	"golang.org/x/crypto/ssh"
	"io"
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

func (h *Host) Tailf(fileName string) (io.ReadCloser, error) {
	//dial a new client for this
	client, err := ssh.Dial("tcp", h.IP+":"+h.Port, &ssh.ClientConfig{
		User: h.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
	})
	if err!=nil {
		return nil,err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	type readCloser struct {
		io.Reader
		io.Closer
	}

	reader ,err := session.StdoutPipe();
	if err!=nil {
		return nil, err
	}
	var rc = readCloser{ reader, session }
	if err:= session.Start("tail -f "+fileName);err!=nil {
		rc.Close()
		return nil,err
	}
	return rc,nil
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
