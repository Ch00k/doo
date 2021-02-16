package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var fixtures *testfixtures.Loader
var db *gorm.DB
var r *gin.Engine

var nullTime = gorm.DeletedAt{Valid: false}

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

func parseISO8601(t string) time.Time {
	v, err := time.Parse(time.RFC3339, t)
	if err != nil {
		panic(err)
	}
	return v
}

func TestListEntries(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entries", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	entry1 := Entry{Text: "wubalubadubdub", Comments: []Comment{}}
	entry1.Model = gorm.Model{ID: 1, CreatedAt: parseISO8601("2020-12-31T23:59:59+01:00"), UpdatedAt: parseISO8601("2020-12-31T23:59:59+01:00"), DeletedAt: nullTime}
	entry2 := Entry{Text: "gazorpazorp", CompletedAt: parseISO8601("2021-01-02T23:59:59+01:00"), Comments: []Comment{}}
	entry2.Model = gorm.Model{ID: 2, CreatedAt: parseISO8601("2021-01-01T23:59:59+01:00"), UpdatedAt: parseISO8601("2021-01-01T23:59:59+01:00"), DeletedAt: nullTime}
	entry3 := Entry{Text: "foobarbaz", CompletedAt: parseISO8601("2021-01-02T23:59:59+01:00"), Comments: []Comment{}}
	entry3.Model = gorm.Model{ID: 3, CreatedAt: parseISO8601("2021-01-01T23:59:59+01:00"), UpdatedAt: parseISO8601("2021-01-01T23:59:59+01:00"), DeletedAt: gorm.DeletedAt{Time: parseISO8601("2021-01-03T23:59:59+01:00"), Valid: true}}

	expectedBody := []Entry{entry1, entry2, entry3}
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
