package demo

type User struct {
	Id    int32  `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type GetUserRequest struct {
	Id int32 `json:"id,omitempty"`
}

type GetUserResponse struct {
	User *User `json:"user,omitempty"`
}

type ListUsersRequest struct {
	PageSize int32 `json:"page_size,omitempty"`
}

type ListUsersResponse struct {
	User *User `json:"user,omitempty"`
	Page int32 `json:"page,omitempty"`
}

type UploadUsersRequest struct {
	User *User `json:"user,omitempty"`
}

type UploadUsersResponse struct {
	Count   int32  `json:"count,omitempty"`
	Message string `json:"message,omitempty"`
}

type ChatMessage struct {
	From      string `json:"from,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
