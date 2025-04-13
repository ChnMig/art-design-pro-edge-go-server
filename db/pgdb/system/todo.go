package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
	"api-server/db/pgdb"
)

// FindTodoList 查询Todo列表(带分页)
func FindTodoList(todo *SystemUserTodo, page, pageSize int) ([]SystemUserTodo, int64, error) {
	var todos []SystemUserTodo
	var total int64
	db := pgdb.GetClient()

	// 构建基础查询
	baseQuery := db.Model(&SystemUserTodo{}).Where("deleted_at IS NULL")

	// 构建条件查询
	if todo.CreatorUserID != 0 {
		baseQuery = baseQuery.Where("system_user_id = ?", todo.CreatorUserID)
	}
	if todo.AssigneeUserID != 0 {
		baseQuery = baseQuery.Where("assignee_user_id = ?", todo.AssigneeUserID)
	}
	if todo.Title != "" {
		baseQuery = baseQuery.Where("title LIKE ?", "%"+todo.Title+"%")
	}
	if todo.Status != 0 {
		baseQuery = baseQuery.Where("status = ?", todo.Status)
	}

	// 获取符合条件的总记录数
	baseQuery.Count(&total)

	// 准备查询
	queryOrder := baseQuery.Order("created_at DESC") // 按创建时间倒序排序

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := queryOrder.Find(&todos).Error; err != nil {
			zap.L().Error("failed to find all todo list", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := queryOrder.Offset((page - 1) * pageSize).Limit(pageSize).Find(&todos).Error; err != nil {
			zap.L().Error("failed to find todo list", zap.Error(err))
			return nil, 0, err
		}
	}

	return todos, total, nil
}

// AddTodo 新增Todo
func AddTodo(todo *SystemUserTodo) error {
	tx := pgdb.GetClient().Begin()
	if err := tx.Create(todo).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to add todo", zap.Error(err))
		return err
	}

	// 添加Todo日志
	log := SystemUserTodoLog{
		SystemUserTodoID: todo.ID,
		Content:          "创建了任务",
	}
	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to add todo log", zap.Error(err))
		return err
	}

	return tx.Commit().Error
}

// AddTodoComment 新增Todo评论
func AddTodoComment(comment *SystemUserTodoComments) error {
	tx := pgdb.GetClient().Begin()
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to add todo comment", zap.Error(err))
		return err
	}

	// 添加Todo日志
	log := SystemUserTodoLog{
		SystemUserTodoID: comment.SystemUserTodoID,
		Content:          "添加了评论",
	}
	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to add todo log", zap.Error(err))
		return err
	}

	return tx.Commit().Error
}

// AddTodoStep 新增Todo步骤
func AddTodoStep(step *SystemUserTodoStep) error {
	tx := pgdb.GetClient().Begin()
	if err := tx.Create(step).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to add todo step", zap.Error(err))
		return err
	}

	// 添加Todo日志
	log := SystemUserTodoLog{
		SystemUserTodoID: step.SystemUserTodoID,
		Content:          "添加了步骤",
	}
	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to add todo log", zap.Error(err))
		return err
	}

	return tx.Commit().Error
}

// FindTodoComments 查询Todo评论(带分页)
func FindTodoComments(todoID uint, page, pageSize int) ([]SystemUserTodoComments, int64, error) {
	var comments []SystemUserTodoComments
	var total int64
	db := pgdb.GetClient()

	// 构建基础查询
	baseQuery := db.Model(&SystemUserTodoComments{}).
		Where("system_user_todo_id = ? AND deleted_at IS NULL", todoID)

	// 获取符合条件的总记录数
	if err := baseQuery.Count(&total).Error; err != nil {
		zap.L().Error("failed to count todo comments", zap.Error(err))
		return nil, 0, err
	}

	// 构建排序查询
	queryOrder := baseQuery.Order("created_at ASC") // 按创建时间升序排序

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := queryOrder.Find(&comments).Error; err != nil {
			zap.L().Error("failed to find all todo comments", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := queryOrder.Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments).Error; err != nil {
			zap.L().Error("failed to find todo comments with pagination", zap.Error(err))
			return nil, 0, err
		}
	}

	return comments, total, nil
}

// FindTodoSteps 查询Todo步骤
func FindTodoSteps(todoID uint) ([]SystemUserTodoStep, error) {
	var steps []SystemUserTodoStep

	if err := pgdb.GetClient().
		Where("system_user_todo_id = ?", todoID).
		Order("created_at ASC"). // 按创建时间升序排序
		Find(&steps).Error; err != nil {
		zap.L().Error("failed to find todo steps", zap.Error(err))
		return nil, err
	}

	return steps, nil
}

// FindTodoLogs 查询Todo日志
func FindTodoLogs(todoID uint) ([]SystemUserTodoLog, error) {
	var logs []SystemUserTodoLog

	if err := pgdb.GetClient().
		Where("system_user_todo_id = ?", todoID).
		Order("created_at DESC"). // 按创建时间倒序排序
		Find(&logs).Error; err != nil {
		zap.L().Error("failed to find todo logs", zap.Error(err))
		return nil, err
	}

	return logs, nil
}

// DeleteTodo 删除Todo
func DeleteTodo(todo *SystemUserTodo) error {
	if err := pgdb.GetClient().Delete(todo).Error; err != nil {
		zap.L().Error("failed to delete todo", zap.Error(err))
		return err
	}
	return nil
}

// DeleteTodoWithRelated 删除Todo并连带删除相关数据
func DeleteTodoWithRelated(todoID uint) error {
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		// 删除相关的评论
		if err := tx.Where("system_user_todo_id = ?", todoID).Delete(&SystemUserTodoComments{}).Error; err != nil {
			zap.L().Error("failed to delete todo comments", zap.Error(err))
			return err
		}

		// 删除相关的步骤
		if err := tx.Where("system_user_todo_id = ?", todoID).Delete(&SystemUserTodoStep{}).Error; err != nil {
			zap.L().Error("failed to delete todo steps", zap.Error(err))
			return err
		}

		// 删除相关的日志
		if err := tx.Where("system_user_todo_id = ?", todoID).Delete(&SystemUserTodoLog{}).Error; err != nil {
			zap.L().Error("failed to delete todo logs", zap.Error(err))
			return err
		}

		// 最后删除Todo本身
		if err := tx.Delete(&SystemUserTodo{Model: gorm.Model{ID: todoID}}).Error; err != nil {
			zap.L().Error("failed to delete todo", zap.Error(err))
			return err
		}

		return nil
	})
}

// UpdateTodoStatus 更新Todo状态
func UpdateTodoStatus(todoID uint, status uint) error {
	tx := pgdb.GetClient().Begin()

	// 更新状态
	if err := tx.Model(&SystemUserTodo{}).
		Where("id = ?", todoID).
		Update("status", status).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to update todo status", zap.Error(err))
		return err
	}

	// 添加日志
	statusText := "未完成"
	if status == 2 {
		statusText = "已完成"
	}

	log := SystemUserTodoLog{
		SystemUserTodoID: todoID,
		Content:          "将任务状态修改为：" + statusText,
	}

	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to add todo status log", zap.Error(err))
		return err
	}

	return tx.Commit().Error
}

// GetTodo 查询单个Todo
func GetTodo(todoID uint) (SystemUserTodo, error) {
	var todo SystemUserTodo

	if err := pgdb.GetClient().
		Where("id = ? AND deleted_at IS NULL", todoID).
		First(&todo).Error; err != nil {
		zap.L().Error("failed to get todo", zap.Error(err))
		return todo, err
	}

	return todo, nil
}

// UpdateTodo 更新Todo全部信息
func UpdateTodo(todo *SystemUserTodo) error {
	tx := pgdb.GetClient().Begin()

	// 首先获取原始数据以记录变更
	var originalTodo SystemUserTodo
	if err := tx.Where("id = ?", todo.ID).First(&originalTodo).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to get original todo", zap.Error(err))
		return err
	}

	// 保存更新前的状态值，用于判断状态是否变更
	originalStatus := originalTodo.Status

	// 更新Todo
	if err := tx.Model(&SystemUserTodo{}).Where("id = ?", todo.ID).Updates(map[string]interface{}{
		"title":            todo.Title,
		"content":          todo.Content,
		"deadline":         todo.Deadline,
		"priority":         todo.Priority,
		"status":           todo.Status,
		"assignee_user_id": todo.AssigneeUserID,
	}).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to update todo", zap.Error(err))
		return err
	}

	// 添加更新日志
	var logContent string

	// 判断状态是否变更，如果变更了则添加状态变更日志
	if originalStatus != todo.Status {
		statusText := "未完成"
		if todo.Status == 2 {
			statusText = "已完成"
		}
		logContent = "将任务状态修改为：" + statusText
	} else {
		logContent = "更新了任务信息"
	}

	log := SystemUserTodoLog{
		SystemUserTodoID: todo.ID,
		Content:          logContent,
	}

	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		zap.L().Error("failed to add todo update log", zap.Error(err))
		return err
	}

	return tx.Commit().Error
}
