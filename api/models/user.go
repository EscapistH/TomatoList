// models/user.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// // 与Python SQLAlchemy对比：
// # class User(Base):
// #     __tablename__ = "users"
// #     id = Column(Integer, primary_key=True, index=True)
// #     email = Column(String, unique=True, index=True, nullable=False)
// #     password = Column(String, nullable=False)
// #     name = Column(String, nullable=False)
// #     created_at = Column(DateTime, default=datetime.utcnow)
// #     last_login = Column(DateTime)
type User struct {
	gorm.Model           // 内嵌gorm.Model，包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Email      string    `json:"email" gorm:"uniqueIndex;not null"` // 邮箱，唯一索引
	Password   string    `json:"-" gorm:"not null"`                 // 密码，不序列化到JSON
	Name       string    `json:"name" gorm:"not null"`              // 用户名
	LastLogin  time.Time `json:"lastLogin"`                         // 最后登录时间

	// 关联关系
	Tasks     []Task     `json:"tasks,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`     // 用户的任务
	Pomodoros []Pomodoro `json:"pomodoros,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // 用户的番茄钟记录
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
