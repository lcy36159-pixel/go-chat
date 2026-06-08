package service

import "errors"

// Validation errors — returned when required fields are missing or invalid.
var (
	ErrUserIDRequired     = errors.New("user_id is required")
	ErrChatIDRequired     = errors.New("chat_id is required")
	ErrSenderIDRequired   = errors.New("sender_id is required")
	ErrLastReadIDRequired = errors.New("last_read_message_id is required")
	ErrContentEmpty       = errors.New("content is empty")
	ErrNameRequired       = errors.New("name is required")
	ErrCannotChatWithSelf = errors.New("cannot chat with yourself")
)

// Auth errors.
var (
	ErrInvalidRegisterInput = errors.New("account, username and password are required")
	ErrInvalidLoginInput    = errors.New("account and password are required")
	ErrWeakPassword         = errors.New("password must be at least 6 characters")
	ErrAccountTaken         = errors.New("account already exists")
	ErrInvalidCredentials   = errors.New("invalid account or password")
)

// Business-rule errors.
var (
	ErrNotChatMember    = errors.New("not a chat member")
	ErrNotGroupChat     = errors.New("chat is not a group")
	ErrAlreadyMember    = errors.New("user is already a member")
	ErrInvalidMessageID = errors.New("invalid message id")
)
