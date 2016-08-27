package move

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/server"
)

// A mailbox supporting the MOVE extension.
type Mailbox interface {
	// Move the specified message(s) to the end of the specified destination
	// mailbox. This means that a new message is created in the target mailbox
	// with a new UID, the original message is removed from the source mailbox,
	// and it appears to the client as a single action.
	//
	// If the destination mailbox does not exist, a server SHOULD return an error.
	// It SHOULD NOT automatically create the mailbox.
	MoveMessages(uid bool, seqset *imap.SeqSet, dest string) error
}

type handler struct {
	Command
}

func (h *handler) handle(uid bool, conn server.Conn) error {
	mailbox := conn.Context().Mailbox
	if mailbox == nil {
		return server.ErrNoMailboxSelected
	}

	if m, ok := mailbox.(Mailbox); ok {
		return m.MoveMessages(uid, h.SeqSet, h.Mailbox)
	}
	return errors.New("MOVE extension not supported")
}

func (h *handler) Handle(conn server.Conn) error {
	return h.handle(false, conn)
}

func (h *handler) UidHandle(conn server.Conn) error {
	return h.handle(true, conn)
}

type extension struct{}

func (ext *extension) Capabilities(c server.Conn) []string {
	if c.Context().State&imap.AuthenticatedState != 0 {
		return []string{Capability}
	}
	return nil
}

func (ext *extension) Command(name string) server.HandlerFactory {
	if name != commandName {
		return nil
	}

	return func() server.Handler {
		return &handler{}
	}
}

func NewExtension() server.Extension {
	return &extension{}
}
