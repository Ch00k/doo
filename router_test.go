package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	"gorm.io/gorm/clause"
)

var (
	fixtures *testfixtures.Loader
	db       *gorm.DB
	r        *gin.Engine
)

func makeRequest(method, path string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)
	r.ServeHTTP(w, req)
	return w
}

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
	w := makeRequest("GET", "/entries", nil)
	assert.Equal(t, 200, w.Code)
	assert.True(t, areEqualJSON(readExpectedResponse("list_entries.json"), w.Body.String()))
}

func TestGetEntryExisting(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("GET", "/entries/2", nil)
	assert.Equal(t, 200, w.Code)
	assert.True(t, areEqualJSON(readExpectedResponse("get_entry.json"), w.Body.String()))
}

func TestGetEntryNonExistent(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("GET", "/entries/42", nil)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestCreateEntry(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries", readRequestBody("create_entry.json"))
	assert.Equal(t, 201, w.Code)

	var entry, dbEntry Entry

	json.Unmarshal(w.Body.Bytes(), &entry)
	res := db.First(&dbEntry, entry.ID)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry, entry)
}

func TestCreateEntryWithNewTags(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries", readRequestBody("create_entry_with_new_tags.json"))
	assert.Equal(t, 201, w.Code)

	var entry, dbEntry Entry

	json.Unmarshal(w.Body.Bytes(), &entry)
	res := db.Preload("Tags").First(&dbEntry, entry.ID)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry, entry)
}

func TestCreateEntryWithExistingTags(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries", readRequestBody("create_entry_with_existing_tags.json"))
	assert.Equal(t, 201, w.Code)

	var entry, dbEntry Entry

	json.Unmarshal(w.Body.Bytes(), &entry)
	res := db.Preload("Tags").First(&dbEntry, entry.ID)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry, entry)
}

func TestCreateEntryTextAbsent(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries", bytes.NewBuffer([]byte("{}")))
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"Key: 'Entry.Text' Error:Field validation for 'Text' failed on the 'required' tag\"}", w.Body.String())
}

func TestDeleteEntryExisting(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("DELETE", "/entries/2", nil)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	var entries []Entry
	db.Find(&entries)
	assert.Len(t, entries, 2)
}

func TestDeleteEntryNonExistent(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("DELETE", "/entries/42", nil)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestUpdateEntryExisting(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("PUT", "/entries/1", readRequestBody("update_entry.json"))
	assert.Equal(t, 200, w.Code)

	var entry, dbEntry Entry

	json.Unmarshal(w.Body.Bytes(), &entry)
	res := db.First(&dbEntry, entry.ID)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry, entry)
}

func TestUpdateEntryError(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("PUT", "/entries/1", bytes.NewBuffer([]byte("{}")))
	// This behaviour is debatable. The endpoint could also return a 200. See the TODO in models.go
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"Key: 'Entry.Text' Error:Field validation for 'Text' failed on the 'required' tag\"}", w.Body.String())
}

func TestUpdateEntryNonExistent(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("PUT", "/entries/42", readRequestBody("update_entry.json"))
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestCompleteEntryExisting(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/1/complete", nil)
	assert.Equal(t, 200, w.Code)

	var entry, dbEntry Entry

	json.Unmarshal(w.Body.Bytes(), &entry)
	res := db.Preload(clause.Associations).First(&dbEntry, 1)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry, entry)
	assert.NotNil(t, dbEntry.CompletedAt)
}

func TestCompleteEntryNonExistent(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/42/complete", nil)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestCreateCommentExistingEntry(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/2/comments", readRequestBody("create_comment.json"))
	assert.Equal(t, 201, w.Code)

	var comment, dbComment Comment

	json.Unmarshal(w.Body.Bytes(), &comment)
	res := db.First(&dbComment, comment.ID)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbComment, comment)
}

func TestCreateCommentNonExistentEntry(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/42/comments", readRequestBody("create_comment.json"))
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestCreateCommentTextAbsent(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/2/comments", bytes.NewBuffer([]byte("{}")))
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"Key: 'Comment.Text' Error:Field validation for 'Text' failed on the 'required' tag\"}", w.Body.String())
}

func TestDeleteCommentExisting(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("DELETE", "/entries/1/comments/2", nil)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	var comments []Comment
	db.Where(&Comment{EntryID: 1}).Find(&comments)
	assert.Len(t, comments, 1)
}

func TestDeleteCommentNonExistent(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("DELETE", "/entries/1/comments/42", nil)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestDeleteCommentNonExistentEntry(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("DELETE", "/entries/42/comments/1", nil)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestListTags(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("GET", "/tags", nil)
	assert.Equal(t, 200, w.Code)
	assert.True(t, areEqualJSON(readExpectedResponse("list_tags.json"), w.Body.String()))
}

func TestGetTagExisting(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("GET", "/tags/food", nil)
	assert.Equal(t, 200, w.Code)
	assert.True(t, areEqualJSON(readExpectedResponse("get_tag.json"), w.Body.String()))
}

func TestGetTagNonExistent(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("GET", "/tags/space", nil)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestTagEntryWithExistingTag(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/2/tag", readRequestBody("tag_entry_with_existing_tag.json"))
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	var dbEntry Entry

	res := db.Preload("Tags").First(&dbEntry, 2)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry.Tags, []Tag{{Name: "family"}})
}

func TestTagEntryWithNewTag(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/2/tag", readRequestBody("tag_entry_with_new_tag.json"))
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	var dbEntry Entry

	res := db.Preload("Tags").First(&dbEntry, 2)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry.Tags, []Tag{{Name: "C137"}})
}

func TestTagEntryNonExistentEntry(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/42/tag", readRequestBody("tag_entry_with_new_tag.json"))
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestTagEntryError(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("POST", "/entries/42/tag", bytes.NewBuffer([]byte("{}")))
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"json: cannot unmarshal object into Go value of type []main.Tag\"}", w.Body.String())
}

func TestUntagEntry(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("DELETE", "/entries/3/tags/food", nil)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	var dbEntry Entry

	res := db.Preload("Tags").First(&dbEntry, 2)
	if res.RowsAffected < 1 {
		panic("Entry not found in db")
	}
	assert.Equal(t, dbEntry.Tags, []Tag{})
}

func TestUntagEntryNonExistentEntry(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("DELETE", "/entries/42/tags/food", nil)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestUntagEntryNonExistentTag(t *testing.T) {
	prepareTestDatabase()
	w := makeRequest("DELETE", "/entries/3/tags/blah", nil)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}
