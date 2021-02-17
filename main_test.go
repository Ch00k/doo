package main

import (
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

func readExpected(filename string) string {
	data, err := ioutil.ReadFile(fmt.Sprintf("test_data/responses/%s", filename))
	if err != nil {
		panic(err)
	}
	return string(data)
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

	return reflect.DeepEqual(j1, j2)
}

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
	assert.True(t, areEqualJSON(readExpected("list.json"), w.Body.String()))
}

func TestGetEntryExisting(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entries/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.True(t, areEqualJSON(readExpected("get.json"), w.Body.String()))
}

func TestGetEntryNonExistent(t *testing.T) {
	prepareTestDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/entries/42", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}
