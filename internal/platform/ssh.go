package platform

import (
	"net"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	sshAgentEnv = "SSH_AUTH_SOCK"
)

type SSH struct {
	client *ssh.Client
}

// TODO accept more method of connection than only via ssh-agent
func NewSSHClient(username, addr string) (SSH, error) {
	sshClient, err := createAgentAuth(username, addr)
	if err != nil {
		return SSH{}, err
	}
	return SSH{
		client: sshClient,
	}, nil
}

func createAgentAuth(username, addr string) (*ssh.Client, error) {
	s := os.Getenv(sshAgentEnv)
	if s == "" {
		return nil, errors.Errorf("%s environment varible empty", sshAgentEnv)
	}

	agentConn, err := net.Dial("unix", s)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't connect via ssh-agent")
	}

	auth := ssh.PublicKeysCallback(agent.NewClient(agentConn).Signers)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{auth},
	}

	return ssh.Dial("tcp", addr, config)
}
