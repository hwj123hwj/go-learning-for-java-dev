package controller

import (
	"user-api-advanced/dto"
	"user-api-advanced/model"
	"user-api-advanced/service"
	"user-api-advanced/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *service.UserService
}

func NewAuthController(userService *service.UserService) *AuthController {
	return &AuthController{userService: userService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "参数校验失败: "+err.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalError(ctx, "密码加密失败")
		return
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Age:      req.Age,
	}

	if err := c.userService.CreateUser(user); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "注册成功"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "参数校验失败: "+err.Error())
		return
	}

	user, err := c.userService.GetUserByEmail(req.Email)
	if err != nil {
		utils.Unauthorized(ctx, "邮箱或密码错误")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.Unauthorized(ctx, "邮箱或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.InternalError(ctx, "生成Token失败")
		return
	}

	utils.Success(ctx, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
