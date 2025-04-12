package todo

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

// 查询 Todo 列表（带分页）
func FindTodoList(c *gin.Context) {
	params := &struct {
		Title  string `json:"title" form:"title"`
		Status uint   `json:"status" form:"status"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)
	todo := system.SystemUserTodo{
		Title:  params.Title,
		Status: params.Status,
	}

	todos, total, err := system.FindTodoList(&todo, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项失败")
		return
	}

	response.ReturnOkWithCount(c, int(total), todos)
}

// 查询单个 Todo
func GetTodo(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	todo, err := system.GetTodo(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项失败")
		return
	}

	// 获取步骤
	steps, err := system.FindTodoSteps(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项步骤失败")
		return
	}

	// 获取评论
	comments, err := system.FindTodoComments(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项评论失败")
		return
	}

	// 获取日志
	logs, err := system.FindTodoLogs(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项日志失败")
		return
	}

	response.ReturnOk(c, gin.H{
		"todo":     todo,
		"steps":    steps,
		"comments": comments,
		"logs":     logs,
	})
}

// 新增 Todo
func AddTodo(c *gin.Context) {
	params := &struct {
		Title   string `json:"title" form:"title" binding:"required"`
		Content string `json:"content" form:"content" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 从Token中获取用户ID
	uID := c.GetString(middleware.JWTDataKey)
	if uID == "" {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}
	id, err := strconv.ParseUint(uID, 10, 64)
	if err != nil {
		response.ReturnError(c, response.UNAUTHENTICATED, "无效的用户ID")
		return
	}

	todo := system.SystemUserTodo{
		CreatorUserID: uint(id),
		Title:         params.Title,
		Content:       params.Content,
		Status:        1, // 默认未完成
	}

	if err := system.AddTodo(&todo); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加待办事项失败")
		return
	}

	response.ReturnOk(c, todo)
}

// 更新 Todo 状态
func UpdateTodoStatus(c *gin.Context) {
	params := &struct {
		ID     uint `json:"id" form:"id" binding:"required"`
		Status uint `json:"status" form:"status" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 验证状态值
	if params.Status != 1 && params.Status != 2 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的状态值，应为 1(未完成) 或 2(已完成)")
		return
	}

	if err := system.UpdateTodoStatus(params.ID, params.Status); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新待办事项状态失败")
		return
	}

	response.ReturnOk(c, nil)
}

// 删除 Todo
func DeleteTodo(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	if err := system.DeleteTodoWithRelated(params.ID); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除待办事项失败")
		return
	}

	response.ReturnOk(c, nil)
}

// 新增 Todo 评论
func AddTodoComment(c *gin.Context) {
	params := &struct {
		TodoID  uint   `json:"todo_id" form:"todo_id" binding:"required"`
		Content string `json:"content" form:"content" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 从Token中获取用户ID
	uID := c.GetString(middleware.JWTDataKey)
	if uID == "" {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}
	id, err := strconv.ParseUint(uID, 10, 64)
	if err != nil {
		response.ReturnError(c, response.UNAUTHENTICATED, "无效的用户ID")
		return
	}

	comment := system.SystemUserTodoComments{
		SystemUserTodoID: params.TodoID,
		SystemUserID:     uint(id),
		Content:          params.Content,
	}

	if err := system.AddTodoComment(&comment); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加评论失败")
		return
	}

	response.ReturnOk(c, comment)
}

// 查询 Todo 评论列表（不带分页）
func FindTodoComments(c *gin.Context) {
	params := &struct {
		TodoID uint `json:"todo_id" form:"todo_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	comments, err := system.FindTodoComments(params.TodoID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询评论失败")
		return
	}

	response.ReturnOk(c, comments)
}

// 新增 Todo 步骤
func AddTodoStep(c *gin.Context) {
	params := &struct {
		TodoID  uint   `json:"todo_id" form:"todo_id" binding:"required"`
		Content string `json:"content" form:"content" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	step := system.SystemUserTodoStep{
		SystemUserTodoID: params.TodoID,
		Content:          params.Content,
	}

	if err := system.AddTodoStep(&step); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加步骤失败")
		return
	}

	response.ReturnOk(c, step)
}

// 查询 Todo 步骤列表（不带分页）
func FindTodoSteps(c *gin.Context) {
	params := &struct {
		TodoID uint `json:"todo_id" form:"todo_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	steps, err := system.FindTodoSteps(params.TodoID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询步骤失败")
		return
	}

	response.ReturnOk(c, steps)
}

// 查询 Todo 日志列表（不带分页）
func FindTodoLogs(c *gin.Context) {
	params := &struct {
		TodoID uint `json:"todo_id" form:"todo_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	logs, err := system.FindTodoLogs(params.TodoID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询日志失败")
		return
	}

	response.ReturnOk(c, logs)
}
