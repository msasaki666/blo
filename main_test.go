package main

import (
	"app/models"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var r *gin.Engine

// func init() {
// 	db, mock, _ = getDbMock()
// 	db.AutoMigrate(models.MigrationTargets()...)
// 	insertTestData(db)
// }

// テストの前後に共通した処理を入れたいときに使う
// １度のテスト実行ごとに一回実行される
func TestMain(m *testing.M) {
	testDbName := "test.db"
	db = setupDb(sqlite.Open(testDbName), models.MigrationTargets()...)
	insertTestData(db)
	r = setupRouter(db)
	m.Run()
	os.Remove(testDbName)
}

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestUserLogin(t *testing.T) {
	w := httptest.NewRecorder()
	body := map[string]string{"email": "test@example.com", "password": "test"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestUserList(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/", nil)
	r.ServeHTTP(w, req)
	fmt.Println(w.Body)

	assert.Equal(t, 200, w.Code)
}

func TestUserShow(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	r.ServeHTTP(w, req)
	fmt.Println(w.Body)

	assert.Equal(t, 200, w.Code)
}

// func TestAuth(t *testing.T) {
// 	secretKey := createSecureRandomStr()
// 	db, _, err := getDbMock()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	PasswordDigest, err := bcrypt.GenerateFromPassword([]byte("test"), 10)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	testUser := models.User{Email: "test@example.com", PasswordDigest: string(PasswordDigest)}
// 	db.Create(&testUser)

// 	_, r := gin.CreateTestContext(httptest.NewRecorder())
// 	authMiddleware, err := createJwtMiddleware(secretKey, "email", db, time.Hour)
// 	if err != nil {
// 		panic("failed to create jwt middleware")
// 	}
// 	if err := authMiddleware.MiddlewareInit(); err != nil {
// 		t.Fatal(err)
// 	}
// 	r.Use(authMiddleware.MiddlewareFunc())
// 	r.GET("/test_auth", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"status": "ok",
// 		})
// 	})
// 	token, _, _ := authMiddleware.TokenGenerator(testUser)
// 	req, _ := http.NewRequest("GET", "/test_auth", nil)
// 	req.Header.Add("Authorization", "Bearer "+token)

// }

// func getDbMock() (*gorm.DB, sqlmock.Sqlmock, error) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer db.Close()

// 	gdb, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn: db,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return gdb, mock, nil
// }

func createSecureRandomStr() string {
	length := 10
	k := make([]byte, length)
	if _, err := rand.Read(k); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", k)
}

func createLoginToken(m *jwt.GinJWTMiddleware, user models.UserLogin) string {
	token, _, _ := m.TokenGenerator(user)
	return token
}

func insertTestData(db *gorm.DB) {
	user := models.User{Email: "test@example.com", PasswordDigest: hashedPassword("test")}
	db.Create(&user)
	post := models.Post{Title: "タイトル", Content: "本文の内容"}
	db.Create(&post)
	db.Model(&post).Association("Tags").Append(&models.Tag{Name: "タグ名"})
}

func hashedPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashed)
}
