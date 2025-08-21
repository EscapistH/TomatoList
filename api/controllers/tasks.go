// controllers/tasks.go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"TomatoList/models"
)

// GetTasks 获取用户的所有任务
// 与Python FastAPI对比：
// @app.get("/tasks")
// async def get_tasks(
//
//	skip: int = 0,
//	limit: int = 100,
//	completed: Optional[bool] = None,
//	db: Session = Depends(get_db),
//	current_user: User = Depends(get_current_user)
//
// ):
//
//	query = db.query(Task).filter(Task.user_id == current_user.id)
//
//	if completed is not None:
//	    query = query.filter(Task.completed == completed)
//
//	tasks = query.offset(skip).limit(limit).all()
//	return tasks
func GetTasks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// 获取查询参数
	completedStr := c.Query("completed")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	// 构建查询
	query := db.Where("user_id = ?", userID)

	// 过滤完成状态
	if completedStr != "" {
		completed, err := strconv.ParseBool(completedStr)
		if err == nil {
			query = query.Where("completed = ?", completed)
		}
	}

	// 分页处理
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var tasks []models.Task
	var total int64

	// 获取总数
	query.Model(&models.Task{}).Count(&total)

	// 获取分页数据
	if result := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		return
	}

	// 返回任务列表和分页信息
	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
			"pages":    (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// GetTask 获取单个任务详情
func GetTask(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// 获取任务ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var task models.Task
	// 查找任务，并确保属于当前用户
	if result := db.Where("id = ? AND user_id = ?", id, userID).First(&task); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask 创建新任务
func CreateTask(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	// 设置任务所属用户
	task.UserID = userID

	// 验证数据
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务标题不能为空"})
		return
	}

	// 设置默认优先级
	if task.Priority == "" {
		task.Priority = "中"
	}

	if result := db.Create(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// 获取任务ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	// 查找现有任务
	var existingTask models.Task
	if result := db.Where("id = ? AND user_id = ?", id, userID).First(&existingTask); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		}
		return
	}

	// 绑定更新数据
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 移除不能更新的字段
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "created_at")

	// 更新任务
	if result := db.Model(&existingTask).Updates(updates); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新任务失败"})
		return
	}

	c.JSON(http.StatusOK, existingTask)
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// 获取任务ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	// 删除任务，并确保属于当前用户
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除任务失败"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}
