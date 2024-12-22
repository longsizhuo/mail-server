package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ContactMessage struct {
	ID      uint   `gorm:"primary_key" json:"-"` // 自增主键
	Name    string `json:"name"`                 // 用户名
	Email   string `json:"email"`                // 邮箱
	Message string `json:"message"`              // 消息
}

var db *gorm.DB

// 初始化数据库
func initDB() {
	// 设置日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// 数据库连接
	dsn := "COMP9900:root@tcp(longsizhuo.com:3306)/resume_website?charset=utf8mb4&parseTime=True&loc=Local" // 使用环境变量管理连接字符串
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库
	if err = db.AutoMigrate(&ContactMessage{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// 使用本地 Postfix 发送邮件
func sendEmailLocal(contact ContactMessage) error {
	message := fmt.Sprintf(
		"From: Contact Form <noreply@localhost>\n"+
			"To: Admin <longsizhuo@gmail.com>\n"+
			"Subject: Contact Form Submission: %s (%s)\n\n"+
			"Name: %s\nEmail: %s\nMessage: %s\n",
		contact.Name, contact.Email, contact.Name, contact.Email, contact.Message,
	)

	cmd := exec.Command("mail", "-s", fmt.Sprintf("Contact Form Submission: %s (%s)", contact.Name, contact.Email), "longsizhuo@gmail.com")

	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Failed to get mail command stdin pipe: %v", err)
		return err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start mail command: %v", err)
		return err
	}
	log.Println("Mail command started successfully")

	// 写入邮件内容
	_, err = pipe.Write([]byte(message))
	if err != nil {
		log.Printf("Failed to write to mail command: %v", err)
		return err
	}

	// 关闭管道
	if err := pipe.Close(); err != nil {
		log.Printf("Failed to close pipe: %v", err)
		return err
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		log.Printf("Mail command execution failed: %v", err)
		return err
	}

	return nil
}

// 处理联系表单
func handleContactForm(c *gin.Context) {
	var contact ContactMessage
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 保存到数据库
	if err := db.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// 发送邮件
	if err := sendEmailLocal(contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message received successfully!"})
}

func main() {
	initDB()

	router := gin.Default()

	// 配置跨域
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 设置路由
	router.POST("/contact", handleContactForm)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	fmt.Println("Server running on port 8181...")
	if err := router.Run("0.0.0.0:8181"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
