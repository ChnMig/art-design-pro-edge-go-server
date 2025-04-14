package todo

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
	"api-server/db/rdb/systemuser"
)

// TodoResponse 包含待办事项及其关联用户信息的响应结构
type TodoResponse struct {
	system.SystemUserTodo
	CreatorName  string `json:"creator_name"`
	AssigneeName string `json:"assignee_name"`
}

// TodoStepResponse 包含待办事项步骤及其操作人信息的响应结构
type TodoStepResponse struct {
	system.SystemUserTodoStep
	OperatorName string `json:"operator_name"` // 操作人用户名
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

	// 自动添加初始步骤
	initialStep := system.SystemUserTodoStep{
		SystemUserTodoID: todo.ID,
		Content:          "任务已创建，开始处理",
		SystemUserID:     uint(id), // 设置操作人ID
	}

	if err := system.AddTodoStep(&initialStep); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加初始步骤失败")
		return
	}

	response.ReturnOk(c, gin.H{
		"todo": todo,
		"step": initialStep,
	})
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

	if err := system.UpdateTodoStatus(params.ID, params.Status); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新待办事项状态失败")
		return
	}

	// 自动添加状态更新步骤
	statusText := "未完成"
	if params.Status == 2 {
		statusText = "已完成"
	}

	step := system.SystemUserTodoStep{
		SystemUserTodoID: params.ID,
		Content:          "任务状态已更新为：" + statusText,
		SystemUserID:     uint(id), // 设置操作人ID
	}

	if err := system.AddTodoStep(&step); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加状态更新步骤失败")
		return
	}

	response.ReturnOk(c, gin.H{
		"status": params.Status,
		"step":   step,
	})
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

	// 获取原始数据以确定变更内容
	originalTodo, err := system.GetTodo(params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取原始待办事项失败")
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

	// 生成更新步骤的内容
	var stepContent string
	if originalTodo.Status != params.Status {
		statusText := "未完成"
		if params.Status == 2 {
			statusText = "已完成"
		}
		stepContent = "任务状态已更新为：" + statusText
	} else if originalTodo.AssigneeUserID != params.AssigneeUserID {
		var assigneeName string
		if params.AssigneeUserID > 0 {
			assigneeInfo, err := systemuser.GetUserFromCache(params.AssigneeUserID)
			if err == nil && assigneeInfo != nil {
				assigneeName = assigneeInfo.Name
			} else {
				assigneeName = "新负责人"
			}
		} else {
			assigneeName = "无负责人"
		}
		stepContent = "任务负责人已更新为：" + assigneeName
	} else {
		stepContent = "任务信息已更新"
	}

	// 添加更新步骤
	step := system.SystemUserTodoStep{
		SystemUserTodoID: params.ID,
		Content:          stepContent,
		SystemUserID:     uint(id), // 设置操作人ID
	}

	if err := system.AddTodoStep(&step); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加更新步骤失败")
		return
	}

	response.ReturnOk(c, gin.H{
		"todo": todo,
		"step": step,
	})
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

// 新增 Todo 步骤
func AddTodoStep(c *gin.Context) {
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

	step := system.SystemUserTodoStep{
		SystemUserTodoID: params.TodoID,
		Content:          params.Content,
		SystemUserID:     uint(id), // 设置操作人ID
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

	// 构建响应数据，添加操作人用户名信息
	stepResponses := make([]TodoStepResponse, 0, len(steps))
	for _, step := range steps {
		stepResp := TodoStepResponse{
			SystemUserTodoStep: step,
			OperatorName:       "",
		}

		// 获取操作人信息
		if step.SystemUserID > 0 {
			operatorInfo, err := systemuser.GetUserFromCache(step.SystemUserID)
			if err == nil && operatorInfo != nil {
				stepResp.OperatorName = operatorInfo.Name
			}
		}

		stepResponses = append(stepResponses, stepResp)
	}

	response.ReturnOk(c, stepResponses)
}
