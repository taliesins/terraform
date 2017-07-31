package hyperv

import (
	"log"
	"github.com/hashicorp/terraform/communicator/winrm"
)

type Config struct {
	User          	string
	Password      	string
	Host  	      	string
	Port	      	int
	HTTPS	      	bool
	Insecure      	bool
	CACert     	*[]byte
	ScriptPath 	string
	Timeout 	string
}

// Client() returns a new client for configuring hyperv.
func (c *Config) Client() (comm *Client, err error) {
	log.Printf("[INFO] HyperV Client configured for HyperV API operations using:\n"+
			"  Host: %s\n"+
			"  Port: %d\n"+
			"  User: %s\n"+
			"  Password: %t\n"+
			"  HTTPS: %t\n"+
			"  Insecure: %t\n"+
			"  CACert: %t\n"+
			"  ScriptPath: %t\n"+
			"  Timeout: %t",
		c.Host,
		c.Port,
		c.User,
		c.Password != "",
		c.HTTPS,
		c.Insecure,
		c.CACert != nil,
		c.ScriptPath,
		c.Timeout,
	)

	return getApiClient(c)
}

// New creates a new communicator implementation over WinRM.
func getWinRMClient(c *Config) (comm *winrm.Communicator, err error) {
	connectionInfo, err := winrm.GetConnectionInfo(c.Host, c.Port, c.User, c.Password, c.HTTPS, c.Insecure, c.Timeout, c.CACert, c.ScriptPath)
	if err != nil {
		return nil, err
	}

	return winrm.GetCommunicator(connectionInfo)
}

func getApiClient(c *Config) (client *Client, err error) {
	client = &Client{
		ElevatedPassword:c.Password,
		ElevatedUser:c.User,
	}

	comm, err := getWinRMClient(c)

	if err != nil {
		return client, err
	}

	client.Communicator = comm

	return client, err
}