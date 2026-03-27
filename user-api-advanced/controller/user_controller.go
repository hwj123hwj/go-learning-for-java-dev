package controller

import (
	"strconv"
	"user-api-advanced/dto"
	"user-api-advanced/model"
	"user-api-advanced/service"
	"user-api-advanced/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		utils.InternalError(ctx, "获取用户列表失败")
		return
	}
	utils.Success(ctx, users)
}

func (c *UserController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "无效的用户ID")
		return
	}
	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		utils.NotFound(ctx, "用户不存在")
		return
	}
	utils.Success(ctx, user)
}

func (c *UserController) Create(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "参数校验失败: "+err.Error())
		return
	}

	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	if err := c.service.CreateUser(user); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, user)
}

func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "参数校验失败: "+err.Error())
		return
	}

	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	if err := c.service.UpdateUser(uint(id), user); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, gin.H{"message": "更新成功"})
}

func (c *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "无效的用户ID")
		return
	}

	if err := c.service.DeleteUser(uint(id)); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.Success(ctx, gin.H{"message": "删除成功"})
}
