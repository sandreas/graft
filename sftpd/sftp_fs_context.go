package sftpd

import (
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
)

type SftpFsContext struct {
	sshc   *ssh.Client
	sshcfg *ssh.ClientConfig
	Sftpc  *sftp.Client
}
func (ctx *SftpFsContext) Disconnect() error {
	ctx.Sftpc.Close()
	ctx.sshc.Close()
	return nil
}

func NewSftpFsContext(user, password, host string) (*SftpFsContext, error) {
	sshcfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	sshc, err := ssh.Dial("tcp", host, sshcfg)
	if err != nil {
		return nil,err
	}

	sftpc, err := sftp.NewClient(sshc)
	if err != nil {
		return nil,err
	}

	ctx := &SftpFsContext{
		sshc:   sshc,
		sshcfg: sshcfg,
		Sftpc:  sftpc,
	}

	return ctx,nil
}
