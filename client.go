package move

import (
	"errors"

	"github.com/emersion/go-imap/common"
	"github.com/emersion/go-imap/commands"
	"github.com/emersion/go-imap/client"
)

type Client struct {
	client *client.Client
}

// Create a new client.
func NewClient(c *client.Client) *Client {
	return &Client{client: c}
}

// Check if the server supports the MOVE extension.
func (c *Client) SupportsMove() bool {
	return c.client.Caps[CommandName]
}

func (c *Client) move(uid bool, seqset *common.SeqSet, dest string) (err error) {
	if c.client.State != common.SelectedState {
		err = errors.New("No mailbox selected")
		return
	}

	var cmd common.Commander
	cmd = &Move{
		SeqSet: seqset,
		Mailbox: dest,
	}
	if uid {
		cmd = &commands.Uid{Cmd: cmd}
	}

	status, err := c.client.Execute(cmd, nil)
	if err != nil {
		return
	}

	err = status.Err()
	return
}

// Moves the specified message(s) to the end of the specified destination
// mailbox.
func (c *Client) Move(seqset *common.SeqSet, dest string) (err error) {
	return c.move(false, seqset, dest)
}

// Identical to Move, but seqset is interpreted as containing unique
// identifiers instead of message sequence numbers.
func (c *Client) UidMove(seqset *common.SeqSet, dest string) (err error) {
	return c.move(true, seqset, dest)
}
