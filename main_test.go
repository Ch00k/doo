package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var fixtures *testfixtures.Loader
var db *gorm.DB
var r *gin.Engine

func TestMain(m *testing.M) {
	db = setupDB()

	var err error

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(sqlDB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test_data/fixtures"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		panic(err)
	}

	r = setupRouter(db)
	os.Exit(m.Run())

}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}

func TestListEntries(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entries", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	entry1 := Entry{ID: 1, CreatedAt: 1609459199000, UpdatedAt: 1609459199000, Text: "wubalubadubdub", Comments: []Comment{}}
	entry2 := Entry{ID: 2, CreatedAt: 1609459299000, UpdatedAt: 1609459299000, CompletedAt: CompletedAt{Int64: 1609459399000, Valid: true}, Text: "gazorpazorp", Comments: []Comment{}}

	expectedBody := []Entry{entry1, entry2}
	var body []Entry
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, expectedBody, body)
}

//func TestGetEntryExisting(t *testing.T) {
//    prepareTestDatabase()

//    w := httptest.NewRecorder()
//    req, _ := http.NewRequest("GET", "/entries/1", nil)
//    r.ServeHTTP(w, req)

//    assert.Equal(t, 200, w.Code)
//    assert.Equal(t, "[]", w.Body.String())
//}

func TestGetEntryNonExistent(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entries/42", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}
