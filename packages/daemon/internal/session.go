package internal

import "github.com/google/uuid"

type WithSessionId interface {
	GetSessionId() string
}

type WithSessionCommand struct {
	SessionId sessionId `long:"sess-id" default:"generate" description:"(internal) 会话ID" hidden:"yes"`
}

func (c *WithSessionCommand) GetSessionId() string {
	return c.SessionId.String()
}

type sessionId string

func (s sessionId) String() string {
	return string(s)
}

func (s *sessionId) UnmarshalFlag(value string) error {
	if value == "generate" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		} else {
			value = id.String()
		}
	} else if _, err := uuid.Parse(value); err != nil {
		return err
	}
	*s = sessionId(value)
	return nil
}
