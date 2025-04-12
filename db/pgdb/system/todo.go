// filepath: /Users/chenming/work/art-design-pro-edge-go-server/db/pgdb/system/todo.go
package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

// TodoWithUser 带用户信息的Todo
type TodoWithUser struct {
	SystemUserTodo
	UserName string `json:"user_name"`
}

// TodoStepWithTodo 带Todo信息的步骤
type TodoStepWithTodo struct {
	SystemUserTodoStep
	TodoTitle string `json:"todo_title"`
}

// TodoCommentWithUser 带用户信息的评论
type TodoCommentWithUser struct {
	SystemUserTodoComments
	UserName string `json:"user_name"`
}

// FindTodoList 查询Todo列表(带分页)
func FindTodoList(todo *SystemUserTodo, page, pageSize int) ([]TodoWithUser, int64, error) {
	var todosWithUser []TodoWithUser
	var total int64
	db := pgdb.GetClient()

	// 构建基础查询
	baseQuery := db.Table("system_user_todos").
		Joins("left join system_users on system_user_todos.system_user_id = system_users.id").
		Where("system_user_todos.deleted_at IS NULL")

	// 构建条件查询
	if todo.SystemUserID != 0 {
		baseQuery = baseQuery.Where("system_user_todos.system_user_id = ?", todo.SystemUserID)
	}
	if todo.Title != "" {
		baseQuery = baseQuery.Where("system_user_todos.title LIKE ?", "%"+todo.Title+"%")
	}
	if todo.Status != 0 {
		baseQuery = baseQuery.Where("system_user_todos.status = ?", todo.Status)
	}

	// 获取符合条件的总记录数
	baseQuery.Count(&total)

	// 应用分页并获取数据
	query := baseQuery.
		Select("system_user_todos.*, system_users.name as user_name").
		Order("system_user_todos.created_at DESC") // 按创建时间倒序排序

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&todosWithUser).Error; err != nil {
		zap.L().Error("failed to find todo list", zap.Error(err))
		return nil, 0, err
	}

	return todosWithUser, total, nil
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

// FindTodoComments 查询Todo评论(不带分页)
func FindTodoComments(todoID uint) ([]TodoCommentWithUser, error) {
	var comments []TodoCommentWithUser

	query := pgdb.GetClient().Table("system_user_todo_comments").
		Joins("left join system_users on system_user_todo_comments.system_user_id = system_users.id").
		Where("system_user_todo_comments.system_user_todo_id = ? AND system_user_todo_comments.deleted_at IS NULL", todoID).
		Select("system_user_todo_comments.*, system_users.name as user_name").
		Order("system_user_todo_comments.created_at ASC") // 按创建时间升序排序

	if err := query.Find(&comments).Error; err != nil {
		zap.L().Error("failed to find todo comments", zap.Error(err))
		return nil, err
	}

	return comments, nil
}

// FindTodoSteps 查询Todo步骤(不带分页)
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

// FindTodoLogs 查询Todo日志(不带分页)
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
func GetTodo(todoID uint) (TodoWithUser, error) {
	var todoWithUser TodoWithUser

	query := pgdb.GetClient().Table("system_user_todos").
		Joins("left join system_users on system_user_todos.system_user_id = system_users.id").
		Where("system_user_todos.id = ? AND system_user_todos.deleted_at IS NULL", todoID).
		Select("system_user_todos.*, system_users.name as user_name")

	if err := query.First(&todoWithUser).Error; err != nil {
		zap.L().Error("failed to get todo", zap.Error(err))
		return todoWithUser, err
	}

	return todoWithUser, nil
}
