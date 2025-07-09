package models

import (
	"time"
	"github.com/uptrace/bun"
)

// Conversation represents a chat session between user and AI.
type Conversation struct {
	bun.BaseModel `bun:"table:conversations,alias:c"`
	ID        string    `bun:",pk,type:uuid,default:gen_random_uuid()"`
	UserID    string    `bun:"user_id,notnull,type:uuid"`
	Title     string    `bun:"title,nullzero"`
	CreatedAt time.Time `bun:",null,default:current_timestamp"`
}

// Message represents a single message in a conversation.
type Message struct {
	bun.BaseModel    `bun:"table:messages,alias:m"`
	ID              string    `bun:",pk,type:uuid,default:gen_random_uuid()"`
	ConversationID  string    `bun:"conversation_id,notnull,type:uuid"`
	Sender          string    `bun:"sender,notnull"` // 'user' or 'ai'
	Content         string    `bun:"content,notnull"`
	CreatedAt       time.Time `bun:",null,default:current_timestamp"`
}

