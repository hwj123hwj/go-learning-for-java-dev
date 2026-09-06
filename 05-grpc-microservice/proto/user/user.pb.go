package user

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

type CreateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type CreateUserResponse struct {
	User *User `json:"user,omitempty"`
}
