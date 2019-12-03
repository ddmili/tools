package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	Host     string
	User     string
	Password string
	*ssh.Client
}

func (this *Client) Connect() error {
	conf := ssh.ClientConfig{
		User:            this.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conf.Auth = append(conf.Auth, ssh.Password(this.Password))
	client, err := ssh.Dial("tcp", this.Host, &conf)
	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}

	this.Client = client
	return nil
}

// Close the connection
func (this *Client) Close() {
	this.Client.Close()
}

func CreateTerminalModes() *ssh.TerminalModes {
	return &ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
}
