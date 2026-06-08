package domain

// CreateGroupChatRequest is the request DTO for creating a group chat.
type CreateGroupChatRequest struct {
	Name    string `json:"name"`
	UserIDs []uint `json:"user_ids"`
}

// CreatePrivateChatRequest is the request DTO for creating a private chat.
type CreatePrivateChatRequest struct {
	TargetUserID uint `json:"target_user_id"`
}

// AddGroupMemberRequest is the request DTO for adding a member to a group chat.
type AddGroupMemberRequest struct {
	UserID uint `json:"user_id"`
}

// MarkMessagesReadRequest is the request DTO for marking messages as read.
type MarkMessagesReadRequest struct {
	LastReadMessageID uint `json:"last_read_message_id"`
}

// RegisterRequest is the request DTO for user registration.
type RegisterRequest struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest is the request DTO for user login.
type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
