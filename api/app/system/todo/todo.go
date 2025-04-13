package todo

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/config"
	"api-server/db/pgdb/system"
	"api-server/db/rdb/systemuser"
)

// TodoResponse 包含待办事项及其关联用户信息的响应结构
type TodoResponse struct {
	system.SystemUserTodo
	CreatorName  string `json:"creator_name"`
	AssigneeName string `json:"assignee_name"`
}

// TodoCommentResponse 包含待办事项评论及其关联用户信息的响应结构
type TodoCommentResponse struct {
	system.SystemUserTodoComments
	UserName string `json:"user_name"`
}

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

	// 从数据库获取原始待办事项列表
	todos, total, err := system.FindTodoList(&todo, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项失败")
		return
	}

	// 构建响应数据，添加用户信息
	todoResponses := make([]TodoResponse, 0, len(todos))
	for _, todo := range todos {
		todoResp := TodoResponse{
			SystemUserTodo: todo,
			CreatorName:    "",
			AssigneeName:   "",
		}

		// 获取创建者信息
		if todo.CreatorUserID > 0 {
			creatorInfo, err := systemuser.GetUserFromCache(todo.CreatorUserID)
			if err == nil && creatorInfo != nil {
				todoResp.CreatorName = creatorInfo.Name
			}
		}

		// 获取被分配者信息
		if todo.AssigneeUserID > 0 {
			assigneeInfo, err := systemuser.GetUserFromCache(todo.AssigneeUserID)
			if err == nil && assigneeInfo != nil {
				todoResp.AssigneeName = assigneeInfo.Name
			}
		}

		todoResponses = append(todoResponses, todoResp)
	}

	response.ReturnOkWithCount(c, int(total), todoResponses)
}

// 查询单个 Todo
func GetTodo(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 获取待办事项基本信息
	todo, err := system.GetTodo(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项失败")
		return
	}

	// 构建带用户信息的响应数据
	todoResp := TodoResponse{
		SystemUserTodo: todo,
		CreatorName:    "",
		AssigneeName:   "",
	}

	// 获取创建者信息
	if todo.CreatorUserID > 0 {
		creatorInfo, err := systemuser.GetUserFromCache(todo.CreatorUserID)
		if err == nil && creatorInfo != nil {
			todoResp.CreatorName = creatorInfo.Name
		}
	}

	// 获取被分配者信息
	if todo.AssigneeUserID > 0 {
		assigneeInfo, err := systemuser.GetUserFromCache(todo.AssigneeUserID)
		if err == nil && assigneeInfo != nil {
			todoResp.AssigneeName = assigneeInfo.Name
		}
	}

	// 获取步骤
	steps, err := system.FindTodoSteps(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项步骤失败")
		return
	}

	// 获取评论 (传入-1,-1表示获取所有评论，不分页)
	comments, _, err := system.FindTodoComments(params.ID, config.CancelPage, config.CancelPageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项评论失败")
		return
	}

	// 构建评论响应，添加用户信息
	commentResponses := make([]TodoCommentResponse, 0, len(comments))
	for _, comment := range comments {
		commentResp := TodoCommentResponse{
			SystemUserTodoComments: comment,
			UserName:               "",
		}

		// 从缓存获取评论用户信息
		if comment.SystemUserID > 0 {
			userInfo, err := systemuser.GetUserFromCache(comment.SystemUserID)
			if err == nil && userInfo != nil {
				commentResp.UserName = userInfo.Name
			}
		}

		commentResponses = append(commentResponses, commentResp)
	}

	// 获取日志
	logs, err := system.FindTodoLogs(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询待办事项日志失败")
		return
	}

	response.ReturnOk(c, gin.H{
		"todo":     todoResp,
		"steps":    steps,
		"comments": commentResponses,
		"logs":     logs,
	})
}

// 新增 Todo
func AddTodo(c *gin.Context) {
	params := &struct {
		Title          string `json:"title" form:"title" binding:"required"`
		Content        string `json:"content" form:"content" binding:"required"`
		Deadline       string `json:"deadline" form:"deadline"`
		Priority       uint   `json:"priority" form:"priority"`
		Status         uint   `json:"status" form:"status"`
		AssigneeUserID uint   `json:"assignee_user_id" form:"assignee_user_id"`
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
		CreatorUserID:  uint(id),
		AssigneeUserID: params.AssigneeUserID,
		Title:          params.Title,
		Content:        params.Content,
		Deadline:       params.Deadline,
		Priority:       params.Priority,
		Status:         params.Status,
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

// 更新 Todo
func UpdateTodo(c *gin.Context) {
	params := &struct {
		ID             uint   `json:"id" form:"id" binding:"required"`
		Title          string `json:"title" form:"title" binding:"required"`
		Content        string `json:"content" form:"content" binding:"required"`
		Deadline       string `json:"deadline" form:"deadline"`
		Priority       uint   `json:"priority" form:"priority"`
		Status         uint   `json:"status" form:"status"`
		AssigneeUserID uint   `json:"assignee_user_id" form:"assignee_user_id"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 验证状态值
	if params.Status != 0 && params.Status != 1 && params.Status != 2 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的状态值，应为 0(默认), 1(未完成) 或 2(已完成)")
		return
	}

	todo := system.SystemUserTodo{
		Model:          gorm.Model{ID: params.ID},
		Title:          params.Title,
		Content:        params.Content,
		Deadline:       params.Deadline,
		Priority:       params.Priority,
		Status:         params.Status,
		AssigneeUserID: params.AssigneeUserID,
	}

	if err := system.UpdateTodo(&todo); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新待办事项失败")
		return
	}

	response.ReturnOk(c, todo)
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

// 查询 Todo 评论列表（带分页）
func FindTodoComments(c *gin.Context) {
	params := &struct {
		TodoID uint `json:"todo_id" form:"todo_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 获取分页参数
	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)

	// 获取原始评论列表
	comments, total, err := system.FindTodoComments(params.TodoID, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询评论失败")
		return
	}

	// 构建评论响应，添加用户信息
	commentResponses := make([]TodoCommentResponse, 0, len(comments))
	for _, comment := range comments {
		commentResp := TodoCommentResponse{
			SystemUserTodoComments: comment,
			UserName:               "",
		}

		// 从缓存获取评论用户信息
		if comment.SystemUserID > 0 {
			userInfo, err := systemuser.GetUserFromCache(comment.SystemUserID)
			if err == nil && userInfo != nil {
				commentResp.UserName = userInfo.Name
			}
		}

		commentResponses = append(commentResponses, commentResp)
	}

	response.ReturnOkWithCount(c, int(total), commentResponses)
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
