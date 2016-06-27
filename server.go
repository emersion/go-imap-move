package move

import (
	"errors"

	"github.com/emersion/go-imap/common"
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
	MoveMessages(uid bool, seqset *common.SeqSet, dest string) error
}

type Handler struct {
	Command
}

func (h *Handler) handle(uid bool, conn *server.Conn) error {
	if conn.Mailbox == nil {
		return server.ErrNoMailboxSelected
	}

	if m, ok := conn.Mailbox.(Mailbox); ok {
		return m.MoveMessages(uid, h.SeqSet, h.Mailbox)
	}
	return errors.New("MOVE extension not supported")
}

func (h *Handler) Handle(conn *server.Conn) error {
	return h.handle(false, conn)
}

func (h *Handler) UidHandle(conn *server.Conn) error {
	return h.handle(true, conn)
}

// Enable the MOVE extension for a server.
func NewServer(s *server.Server) {
	s.RegisterCapability(CommandName, common.SelectedState)

	s.RegisterCommand(CommandName, func() server.Handler {
		return &Handler{}
	})
}
