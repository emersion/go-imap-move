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

func (c *Client) moveWithFallback(uid bool, seqset *imap.SeqSet, dest string) error {
	if ok, err := c.SupportMove(); err != nil {
		return err
	} else if ok {
		return c.move(uid, seqset, dest)
	}

	if c.c.State() != imap.SelectedState {
		return client.ErrNoMailboxSelected
	}

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if uid {
		if err := c.c.UidCopy(seqset, dest); err != nil {
			return err
		}

		if err := c.c.UidStore(seqset, item, flags, nil); err != nil {
			return err
		}
	} else {
		if err := c.c.Copy(seqset, dest); err != nil {
			return err
		}

		if err := c.c.Store(seqset, item, flags, nil); err != nil {
			return err
		}
	}

	return c.c.Expunge(nil)
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

// MoveWithFallback tries to move if the server supports it. If it doesn't, it
// falls back to copy, store and expunge, as defined in RFC 6851 section 3.3.
func (c *Client) MoveWithFallback(seqset *imap.SeqSet, dest string) error {
	return c.moveWithFallback(false, seqset, dest)
}

// UidMoveWithFallback is identical to MoveWithFallback, but seqset is
// interpreted as containing unique identifiers instead of message sequence
// numbers.
func (c *Client) UidMoveWithFallback(seqset *imap.SeqSet, dest string) error {
	return c.moveWithFallback(true, seqset, dest)
}
