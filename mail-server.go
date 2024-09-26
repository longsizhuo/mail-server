package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ContactMessage struct {
	ID      uint   `gorm:"primary_key" json:"-"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

var db *gorm.DB

// 初始化数据库连接
func initDB() {
	// dsn := "root
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // Slow SQL threshold
			LogLevel:      logger.Error, // Log level
			Colorful:      true,
		},
	)

	dsn := "COMP9900:root@tcp(longsizhuo.com:3306)/resume_website?charset=utf8mb4&parseTime=True&loc=Local"

	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&ContactMessage{})
	if err != nil {
		panic("failed to migrate database")
	}
}

// sendEmail S end email
func sendEmail(messaged ContactMessage) error {
	m := gomail.NewMessage()
	name := messaged.Name
	email := messaged.Email
	content := messaged.Message

	// SMTP Server settings
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Set up authentication information.
	from := smtpUsername
	to := os.Getenv("SMTP_USERNAME")
	subject := name + email + "sent you an Email"
	body := content

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Contact Form: %s", subject))
	m.SetBody("text/plain", fmt.Sprintf("Name: %s\nEmail: %s\n\n%s", name, email, body))

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	d := gomail.NewDialer(smtpHost, port, smtpUsername, smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

// 处理联系表单的提交
func handleContactForm(c *gin.Context) {
	var contact ContactMessage
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := db.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	if err := sendEmail(contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message received successfully!"})
}

func main() {
	initDB()

	router := gin.Default()

	// 允许跨域请求
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.POST("/contact", handleContactForm)

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	fmt.Println("Server running on port 8181...")
	if err := router.Run(":8181"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
