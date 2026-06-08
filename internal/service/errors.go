package service

import "errors"

// 輸入驗證錯誤 - 在執行業務邏輯前就能找出的錯誤
var (
	ErrUserIDRequired     = errors.New("user_id is required")
	ErrChatIDRequired     = errors.New("chat_id is required")
	ErrSenderIDRequired   = errors.New("sender_id is required")
	ErrLastReadIDRequired = errors.New("last_read_message_id is required")
	ErrContentEmpty       = errors.New("content is empty")
	ErrNameRequired       = errors.New("name is required")
	ErrCannotChatWithSelf = errors.New("cannot chat with yourself")
)

// 登入/註冊相關錯誤
var (
	ErrInvalidRegisterInput = errors.New("account, username and password are required")
	ErrInvalidLoginInput    = errors.New("account and password are required")
	ErrWeakPassword         = errors.New("password must be at least 6 characters")
	ErrAccountTaken         = errors.New("account already exists")
	ErrInvalidCredentials   = errors.New("invalid account or password")
)

// 業務規則錯誤 - 違反系統規定
var (
	ErrNotChatMember    = errors.New("not a chat member")
	ErrNotGroupChat     = errors.New("chat is not a group")
	ErrAlreadyMember    = errors.New("user is already a member")
	ErrInvalidMessageID = errors.New("invalid message id")
)
