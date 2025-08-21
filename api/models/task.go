// models/task.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务模型
// 与Python SQLAlchemy对比：
// # class Task(Base):
// #     __tablename__ = "tasks"
// #     id = Column(Integer, primary_key=True, index=True)
// #     title = Column(String, nullable=False)
// #     description = Column(Text)
// #     priority = Column(String, default="中")
// #     completed = Column(Boolean, default=False)
// #     due_date = Column(DateTime)
// #     user_id = Column(Integer, ForeignKey("users.id"))
// #     created_at = Column(DateTime, default=datetime.utcnow)
type Task struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`          // 任务标题
	Description string    `json:"description"`                    // 任务描述
	Priority    string    `json:"priority" gorm:"default:'中'"`    // 优先级：高、中、低
	Completed   bool      `json:"completed" gorm:"default:false"` // 是否完成
	DueDate     time.Time `json:"dueDate"`                        // 截止日期
	UserID      uint      `json:"userId" gorm:"not null"`         // 关联的用户ID

	// 关联关系
	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID"`                                  // 关联的用户
	Pomodoros []Pomodoro `json:"pomodoros,omitempty" gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"` // 关联的番茄钟记录
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}

// IsOverdue 检查任务是否过期
func (t *Task) IsOverdue() bool {
	if t.DueDate.IsZero() {
		return false
	}
	return !t.Completed && time.Now().After(t.DueDate)
}

// PomodoroCount 获取任务的番茄钟数量
func (t *Task) PomodoroCount(db *gorm.DB) int64 {
	var count int64
	db.Model(&Pomodoro{}).Where("task_id = ?", t.ID).Count(&count)
	return count
}

// CompletedPomodoroCount 获取任务完成的番茄钟数量
func (t *Task) CompletedPomodoroCount(db *gorm.DB) int64 {
	var count int64
	db.Model(&Pomodoro{}).Where("task_id = ? AND status = ?", t.ID, "已完成").Count(&count)
	return count
}
