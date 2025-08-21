// controllers/auth.go
package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"TomatoList/models"
	"TomatoList/utils"
)

// Register 用户注册控制器
// 与Python FastAPI对比：
// @app.post("/register")
// async def register(user: UserCreate, db: Session = Depends(get_db)):
//
//	# 检查用户是否存在
//	existing_user = db.query(User).filter(User.Email == user.email).first()
//	if existing_user:
//	    raise HTTPException(status_code=400, detail="Email already registered")
//
//	# 密码哈希
//	hashed_password = bcrypt.hashpw(user.password.encode('utf-8'), bcrypt.gensalt())
//
//	# 创建用户
//	db_user = User(email=user.email, password_hash=hashed_password.decode('utf-8'))
//	db.add(db_user)
//	db.commit()
//	db.refresh(db_user)
//
//	return db_user
func Register(c *gin.Context) {
	// 从上下文中获取数据库连接
	db := c.MustGet("db").(*gorm.DB)

	// 定义请求体结构
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	// 检查邮箱是否已存在
	var existingUser models.User
	if result := db.Where("email = ?", request.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	user := models.User{
		Email:    request.Email,
		Password: string(hashedPassword),
		Name:     request.Name,
	}

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回用户信息和令牌（不返回密码）
	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": token,
	})
}

// Login 用户登录控制器
func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 查找用户
	var user models.User
	if result := db.Where("email = ?", request.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	// 更新最后登录时间
	db.Model(&user).Update("last_login", time.Now())

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回用户信息和令牌
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": token,
	})
}
