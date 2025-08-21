// controllers/pomodoros.go
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"TomatoList/models"
)

// StartPomodoro 开始一个番茄钟
// 与Python FastAPI对比：
// @app.post("/pomodoros")
// async def start_pomodoro(
//
//	task_id: int,
//	db: Session = Depends(get_db),
//	current_user: User = Depends(get_current_user)
//
// ):
//
//	# 检查任务是否存在且属于当前用户
//	task = db.query(Task).filter(Task.id == task_id, Task.user_id == current_user.id).first()
//	if not task:
//	    raise HTTPException(status_code=404, detail="任务不存在")
//
//	# 创建番茄钟记录
//	pomodoro = Pomodoro(
//	    task_id=task_id,
//	    user_id=current_user.id,
//	    start_time=datetime.now(),
//	    expected_end_time=datetime.now() + timedelta(minutes=25)
//	)
//
//	db.add(pomodoro)
//	db.commit()
//	db.refresh(pomodoro)
//
//	return pomodoro
func StartPomodoro(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	var request struct {
		TaskID uint `json:"taskId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查任务是否存在且属于当前用户
	var task models.Task
	if result := db.Where("id = ? AND user_id = ?", request.TaskID, userID).First(&task); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		}
		return
	}

	// 创建番茄钟记录
	now := time.Now()
	pomodoro := models.Pomodoro{
		TaskID:          request.TaskID,
		UserID:          userID,
		StartTime:       now,
		ExpectedEndTime: now.Add(25 * time.Minute), // 标准番茄钟25分钟
		Status:          "进行中",
	}

	if result := db.Create(&pomodoro); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建番茄钟失败"})
		return
	}

	c.JSON(http.StatusCreated, pomodoro)
}

// CompletePomodoro 完成一个番茄钟
func CompletePomodoro(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// 获取番茄钟ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的番茄钟ID"})
		return
	}

	// 查找番茄钟记录
	var pomodoro models.Pomodoro
	if result := db.Where("id = ? AND user_id = ?", id, userID).First(&pomodoro); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "番茄钟记录不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取番茄钟失败"})
		}
		return
	}

	// 更新番茄钟状态
	updates := map[string]interface{}{
		"end_time": time.Now(),
		"status":   "已完成",
	}

	if result := db.Model(&pomodoro).Updates(updates); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新番茄钟失败"})
		return
	}

	c.JSON(http.StatusOK, pomodoro)
}

// GetPomodoros 获取用户的番茄钟记录
func GetPomodoros(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// 获取查询参数
	taskIDStr := c.Query("taskId")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	// 构建查询
	query := db.Where("user_id = ?", userID)

	// 过滤特定任务的番茄钟
	if taskIDStr != "" {
		taskID, err := strconv.Atoi(taskIDStr)
		if err == nil {
			query = query.Where("task_id = ?", taskID)
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

	var pomodoros []models.Pomodoro
	var total int64

	// 获取总数
	query.Model(&models.Pomodoro{}).Count(&total)

	// 获取分页数据
	if result := query.Preload("Task").Order("start_time desc").Offset(offset).Limit(pageSize).Find(&pomodoros); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取番茄钟记录失败"})
		return
	}

	// 返回番茄钟记录和分页信息
	c.JSON(http.StatusOK, gin.H{
		"pomodoros": pomodoros,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
			"pages":    (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// GetPomodoroStats 获取番茄钟统计信息
func GetPomodoroStats(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// 获取时间范围参数（默认最近7天）
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 7
	}

	startDate := time.Now().AddDate(0, 0, -days)

	// 查询已完成番茄钟数量
	var completedCount int64
	db.Model(&models.Pomodoro{}).
		Where("user_id = ? AND status = ? AND start_time >= ?", userID, "已完成", startDate).
		Count(&completedCount)

	// 查询总番茄钟时间（分钟）
	var totalMinutes int64
	rows, err := db.Model(&models.Pomodoro{}).
		Select("COALESCE(SUM(EXTRACT(EPOCH FROM (end_time - start_time)) / 60), 0)").
		Where("user_id = ? AND status = ? AND start_time >= ?", userID, "已完成", startDate).
		Rows()

	if err == nil {
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&totalMinutes)
		}
	}

	// 查询每日番茄钟数量
	type DailyStat struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
	var dailyStats []DailyStat

	rows, err = db.Model(&models.Pomodoro{}).
		Select("DATE(start_time) as date, COUNT(*) as count").
		Where("user_id = ? AND status = ? AND start_time >= ?", userID, "已完成", startDate).
		Group("DATE(start_time)").
		Order("date").
		Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date string
			var count int
			rows.Scan(&date, &count)
			dailyStats = append(dailyStats, DailyStat{Date: date, Count: count})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"completedCount": completedCount,
		"totalMinutes":   totalMinutes,
		"dailyStats":     dailyStats,
		"period":         days,
	})
}
