package client

import (
	"fmt"
	"log"
	"time"

	"github.com/gorcon/rcon"
)

type RCONClient struct {
	address  string
	password string
}

// WaitForRCON wait 10 minutes for connect to RCON
func WaitForRCON(host string, password string) (client *RCONClient, err error) {
	for i := 0; i < 40; i++ {
		client, err = NewRCONClient(host, password)
		if err == nil {
			return
		}
		log.Println("Waiting for RCON connection to host")
		time.Sleep(15 * time.Second)
	}
	return
}

func NewRCONClient(address string, password string) (*RCONClient, error) {
	return &RCONClient{
		address:  address,
		password: password,
	}, nil
}

func (c *RCONClient) execCommand(cmd string) (string, error) {
	conn, err := rcon.Dial(c.address, c.password)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return conn.Execute(cmd)
}

func (c *RCONClient) StopServer() error {
	_, err := c.execCommand("stop")
	return err
}

func (c *RCONClient) ListPlayers() (int, error) {
	resp, err := c.execCommand("list")
	if err != nil {
		return 0, err
	}

	var players int
	_, err = fmt.Sscanf(resp, "There are %d", &players)
	if err != nil {
		return 0, nil
	}

	return players, nil
}
