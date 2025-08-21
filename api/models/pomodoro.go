// models/pomodoro.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Pomodoro 番茄钟模型
// 与Python SQLAlchemy对比：
// # class Pomodoro(Base):
// #     __tablename__ = "pomodoros"
// #     id = Column(Integer, primary_key=True, index=True)
// #     task_id = Column(Integer, ForeignKey("tasks.id"))
// #     user_id = Column(Integer, ForeignKey("users.id"))
// #     start_time = Column(DateTime, nullable=False)
// #     end_time = Column(DateTime)
// #     expected_end_time = Column(DateTime, nullable=False)
// #     status = Column(String, default="进行中")  # 进行中, 已完成, 已中断
// #     created_at = Column(DateTime, default=datetime.utcnow)
type Pomodoro struct {
	gorm.Model
	TaskID          uint      `json:"taskId" gorm:"not null"`          // 关联的任务ID
	UserID          uint      `json:"userId" gorm:"not null"`          // 关联的用户ID
	StartTime       time.Time `json:"startTime" gorm:"not null"`       // 开始时间
	EndTime         time.Time `json:"endTime"`                         // 结束时间（实际结束时间）
	ExpectedEndTime time.Time `json:"expectedEndTime" gorm:"not null"` // 预期结束时间
	Status          string    `json:"status" gorm:"default:'进行中'"`     // 状态：进行中、已完成、已中断

	// 关联关系
	Task Task `json:"task,omitempty" gorm:"foreignKey:TaskID"` // 关联的任务
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"` // 关联的用户
}

// TableName 指定表名
func (Pomodoro) TableName() string {
	return "pomodoros"
}

// Duration 计算番茄钟实际持续时间（分钟）
func (p *Pomodoro) Duration() float64 {
	if p.EndTime.IsZero() {
		return 0
	}
	return p.EndTime.Sub(p.StartTime).Minutes()
}

// IsCompleted 检查番茄钟是否已完成
func (p *Pomodoro) IsCompleted() bool {
	return p.Status == "已完成"
}

// IsOvertime 检查番茄钟是否超时
func (p *Pomodoro) IsOvertime() bool {
	if p.EndTime.IsZero() {
		return time.Now().After(p.ExpectedEndTime)
	}
	return p.EndTime.After(p.ExpectedEndTime)
}
