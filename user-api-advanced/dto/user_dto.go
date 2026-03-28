package dto

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
	Age   *int   `json:"age" binding:"omitempty,gte=0,lte=150"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=50"`
	Email string `json:"email" binding:"omitempty,email"`
	Age   *int   `json:"age" binding:"omitempty,gte=0,lte=150"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      *int   `json:"age" binding:"omitempty,gte=0,lte=150"`
}
