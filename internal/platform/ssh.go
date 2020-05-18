package platform

import (
	"net"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func sshAgentAuth(username, addr string) (*ssh.Client, error) {
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
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error {
			return nil
		},
	}

	return ssh.Dial("tcp", addr, config)
}
