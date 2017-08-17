package move

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/commands"
)

// Client is a MOVE client.
type Client struct {
	c *client.Client
}

// NewClient creates a new client.
func NewClient(c *client.Client) *Client {
	return &Client{c: c}
}

// SupportMove checks if the server supports the MOVE extension.
func (c *Client) SupportMove() (bool, error) {
	return c.c.Support(Capability)
}

func (c *Client) move(uid bool, seqset *imap.SeqSet, dest string) error {
	if c.c.State() != imap.SelectedState {
		return client.ErrNoMailboxSelected
	}

	var cmd imap.Commander = &Command{
		SeqSet:  seqset,
		Mailbox: dest,
	}
	if uid {
		cmd = &commands.Uid{Cmd: cmd}
	}

	if status, err := c.c.Execute(cmd, nil); err != nil {
		return err
	} else {
		return status.Err()
	}
}

// Move moves the specified message(s) to the end of the specified destination
// mailbox.
func (c *Client) Move(seqset *imap.SeqSet, dest string) error {
	return c.move(false, seqset, dest)
}

// UidMove is identical to Move, but seqset is interpreted as containing unique
// identifiers instead of message sequence numbers.
func (c *Client) UidMove(seqset *imap.SeqSet, dest string) error {
	return c.move(true, seqset, dest)
}
