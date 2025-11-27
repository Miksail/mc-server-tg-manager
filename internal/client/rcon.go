package client

import (
	"fmt"
	"time"

	"github.com/gorcon/rcon"
)

type RCONClient struct {
	conn *rcon.Conn
}

// WaitForRCON wait 10 minutes for connect to RCON
func WaitForRCON(host string, password string) (client *RCONClient, err error) {
	for i := 0; i < 40; i++ {
		client, err = NewRCONClient(host, password)
		if err == nil {
			return
		}
		time.Sleep(15 * time.Second)
	}
	return
}

func NewRCONClient(address string, password string) (*RCONClient, error) {
	conn, err := rcon.Dial(address, password)
	if err != nil {
		return nil, err
	}
	return &RCONClient{
		conn: conn,
	}, nil
}

func (c *RCONClient) Close() error {
	return c.conn.Close()
}

func (c *RCONClient) StopServer() error {
	_, err := c.conn.Execute("stop")
	return err
}

func (c *RCONClient) ListPlayers() (int, error) {
	resp, err := c.conn.Execute("list")
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
