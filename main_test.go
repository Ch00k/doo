package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var fixtures *testfixtures.Loader
var db *gorm.DB
var r *gin.Engine

func readExpectedResponse(filename string) string {
	data, err := ioutil.ReadFile(fmt.Sprintf("test_data/responses/%s", filename))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func readRequestBody(filename string) *bytes.Buffer {
	data, err := ioutil.ReadFile(fmt.Sprintf("test_data/requests/%s", filename))
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(data)
}

func areEqualJSON(s1, s2 string) bool {
	var j1 interface{}
	var j2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &j1)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(s2), &j2)
	if err != nil {
		panic(err)
	}

	areEqual := reflect.DeepEqual(j1, j2)
	if areEqual {
		return true
	}
	fmt.Println("JSON objects are not equal")
	fmt.Println(j1)
	fmt.Println(j2)
	return false
}

func TestMain(m *testing.M) {
	db = SetupDB()

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

	r = SetupRouter(db)
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
	assert.True(t, areEqualJSON(readExpectedResponse("list.json"), w.Body.String()))
}

func TestGetEntryExisting(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entries/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.True(t, areEqualJSON(readExpectedResponse("get.json"), w.Body.String()))
}

func TestGetEntryNonExistent(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entries/42", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestCreateEntry(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/entries", readRequestBody("create_entry.json"))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var entry, dbEntry Entry

	json.Unmarshal(w.Body.Bytes(), &entry)
	res := db.First(&dbEntry, entry.ID)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry, entry)
}

func TestDeleteEntry(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/entries/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	var entries []Entry
	db.Find(&entries)
	assert.Len(t, entries, 1)
}

func TestUpdateEntry(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/entries/1", readRequestBody("update_entry.json"))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var entry, dbEntry Entry

	json.Unmarshal(w.Body.Bytes(), &entry)
	res := db.First(&dbEntry, entry.ID)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry, entry)
}

func TestCreateComment(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/entries/2/comment", readRequestBody("create_comment.json"))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var comment, dbComment Comment

	json.Unmarshal(w.Body.Bytes(), &comment)
	res := db.First(&dbComment, comment.ID)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbComment, comment)
}

func TestDeleteComment(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/entries/1/comments/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	var comments []Comment
	db.Where(&Comment{EntryID: 1}).Find(&comments)
	assert.Len(t, comments, 1)
}
