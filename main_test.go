package main

import (
	"app/models"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestHealth(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAuth(t *testing.T) {
	secretKey := createSecureRandomStr()
	db, _, err := getDbMock()
	if err != nil {
		t.Fatal(err)
	}
	PasswordDigest, err := bcrypt.GenerateFromPassword([]byte("test"), 10)
	if err != nil {
		t.Fatal(err)
	}
	testUser := models.User{Email: "test@example.com", PasswordDigest: string(PasswordDigest)}
	db.Create(&testUser)

	_, r := gin.CreateTestContext(httptest.NewRecorder())
	authMiddleware, err := createJwtMiddleware(secretKey, "email", db, time.Hour)
	if err != nil {
		panic("failed to create jwt middleware")
	}
	if err := authMiddleware.MiddlewareInit(); err != nil {
		t.Fatal(err)
	}
	r.Use(authMiddleware.MiddlewareFunc())
	r.GET("/test_auth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	token, _, _ := authMiddleware.TokenGenerator(testUser)
	req, _ := http.NewRequest("GET", "/test_auth", nil)
	req.Header.Add("Authorization", "Bearer "+token)

}

func getDbMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gdb, mock, nil
}

func createSecureRandomStr() string {
	length := 10
	k := make([]byte, length)
	if _, err := rand.Read(k); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", k)
}
