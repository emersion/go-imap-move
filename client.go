package move

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/commands"
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
	return c.client.Caps[Capability]
}

func (c *Client) move(uid bool, seqset *imap.SeqSet, dest string) (err error) {
	if c.client.State != imap.SelectedState {
		err = errors.New("No mailbox selected")
		return
	}

	var cmd imap.Commander
	cmd = &Command{
		SeqSet:  seqset,
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
func (c *Client) Move(seqset *imap.SeqSet, dest string) (err error) {
	return c.move(false, seqset, dest)
}

// Identical to Move, but seqset is interpreted as containing unique
// identifiers instead of message sequence numbers.
func (c *Client) UidMove(seqset *imap.SeqSet, dest string) (err error) {
	return c.move(true, seqset, dest)
}
