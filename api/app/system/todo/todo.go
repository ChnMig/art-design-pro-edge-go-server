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

	// 准备更详细的步骤内容
	stepContent := "【创建任务】"

	// 添加标题和内容信息
	stepContent += "\n标题: " + params.Title
	stepContent += "\n内容: " + params.Content

	// 添加截止日期信息（如果有）
	if params.Deadline != "" {
		stepContent += "\n截止日期: " + params.Deadline
	}

	// 添加优先级信息
	priorityText := "中"
	if params.Priority == 1 {
		priorityText = "低"
	} else if params.Priority == 3 {
		priorityText = "高"
	}
	stepContent += "\n优先级: " + priorityText

	// 添加状态信息
	statusText := "未处理"
	if params.Status == 1 {
		statusText = "未处理"
	} else if params.Status == 2 {
		statusText = "处理中"
	} else if params.Status == 3 {
		statusText = "已完成"
	} else if params.Status == 4 {
		statusText = "已取消"
	}
	stepContent += "\n状态: " + statusText

	// 添加负责人信息（如果有）
	if params.AssigneeUserID > 0 {
		assigneeInfo, err := systemuser.GetUserFromCache(params.AssigneeUserID)
		if err == nil && assigneeInfo != nil {
			stepContent += "\n负责人: " + assigneeInfo.Name
		}
	}

	// 自动添加初始步骤
	initialStep := system.SystemUserTodoStep{
		SystemUserTodoID: todo.ID,
		Content:          stepContent,
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

// 获取状态文字描述
func getStatusText(status uint) string {
	switch status {
	case 1:
		return "未处理"
	case 2:
		return "处理中"
	case 3:
		return "已完成"
	case 4:
		return "已取消"
	default:
		return "未知状态"
	}
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
	if params.Status != 1 && params.Status != 2 && params.Status != 3 && params.Status != 4 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的状态值")
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

	// 准备步骤内容，详细记录更新的字段和值
	stepContent := "【更新任务】"

	// 记录标题变更
	if originalTodo.Title != params.Title {
		stepContent += "\n标题: [" + originalTodo.Title + "] -> [" + params.Title + "]"
	}

	// 记录内容变更 (可能较长，只记录是否变更)
	if originalTodo.Content != params.Content {
		stepContent += "\n内容: 已更新"
	}

	// 记录截止日期变更
	if originalTodo.Deadline != params.Deadline {
		oldDeadline := originalTodo.Deadline
		if oldDeadline == "" {
			oldDeadline = "无"
		}
		newDeadline := params.Deadline
		if newDeadline == "" {
			newDeadline = "无"
		}
		stepContent += "\n截止日期: [" + oldDeadline + "] -> [" + newDeadline + "]"
	}

	// 记录优先级变更
	if originalTodo.Priority != params.Priority {
		oldPriority := "中"
		if originalTodo.Priority == 1 {
			oldPriority = "低"
		} else if originalTodo.Priority == 3 {
			oldPriority = "高"
		}

		newPriority := "中"
		if params.Priority == 1 {
			newPriority = "低"
		} else if params.Priority == 3 {
			newPriority = "高"
		}

		stepContent += "\n优先级: [" + oldPriority + "] -> [" + newPriority + "]"
	}

	// 记录状态变更
	if originalTodo.Status != params.Status {
		oldStatus := getStatusText(originalTodo.Status)
		newStatus := getStatusText(params.Status)
		stepContent += "\n状态: [" + oldStatus + "] -> [" + newStatus + "]"
	}

	// 记录负责人变更
	if originalTodo.AssigneeUserID != params.AssigneeUserID {
		var oldAssigneeName, newAssigneeName string

		if originalTodo.AssigneeUserID > 0 {
			oldAssigneeInfo, err := systemuser.GetUserFromCache(originalTodo.AssigneeUserID)
			if err == nil && oldAssigneeInfo != nil {
				oldAssigneeName = oldAssigneeInfo.Name
			} else {
				oldAssigneeName = "ID:" + strconv.FormatUint(uint64(originalTodo.AssigneeUserID), 10)
			}
		} else {
			oldAssigneeName = "无"
		}

		if params.AssigneeUserID > 0 {
			newAssigneeInfo, err := systemuser.GetUserFromCache(params.AssigneeUserID)
			if err == nil && newAssigneeInfo != nil {
				newAssigneeName = newAssigneeInfo.Name
			} else {
				newAssigneeName = "ID:" + strconv.FormatUint(uint64(params.AssigneeUserID), 10)
			}
		} else {
			newAssigneeName = "无"
		}

		stepContent += "\n负责人: [" + oldAssigneeName + "] -> [" + newAssigneeName + "]"
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
		SystemUserTodoID uint   `json:"system_user_todo_id" form:"system_user_todo_id" binding:"required"`
		Content          string `json:"content" form:"content" binding:"required"`
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
		SystemUserTodoID: params.SystemUserTodoID,
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
