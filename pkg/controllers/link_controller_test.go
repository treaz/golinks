package controllers

import (
	"bytes"
	"crhuber/golinks/pkg/database"
	"crhuber/golinks/pkg/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *database.DbConnection {
	// Use in-memory SQLite database
	db, err := database.NewConnection("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	err = db.RunMigration()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Clean up database
	db.Db.Exec("DELETE FROM links")

	return db
}

func TestAppController_CreateLink(t *testing.T) {
	db := setupTestDB(t)
	controller := NewAppController(db)
	r := chi.NewRouter()
	r.Post("/links", controller.CreateLink)

	tests := []struct {
		name           string
		input          models.LinkInput
		expectedStatus int
		verify         func(*testing.T, *database.DbConnection, models.LinkInput)
	}{
		{
			name: "Create Valid Link",
			input: models.LinkInput{
				Keyword:     "test",
				Destination: "https://example.com",
				Description: "Test Link",
			},
			expectedStatus: http.StatusOK,
			verify: func(t *testing.T, db *database.DbConnection, input models.LinkInput) {
				var link models.Link
				err := db.Db.Where("keyword = ?", "test").First(&link).Error
				assert.NoError(t, err)
				assert.Equal(t, "test", link.Keyword)
				assert.Equal(t, "https://example.com", link.Destination)
			},
		},
		{
			name: "Create Duplicate Link",
			input: models.LinkInput{
				Keyword:     "test", // Same as above
				Destination: "https://other.com",
			},
			expectedStatus: http.StatusBadRequest,
			verify:         nil,
		},
		{
			name: "Create Link Starting With Slash",
			input: models.LinkInput{
				Keyword:     "/slash",
				Destination: "https://example.com",
			},
			expectedStatus: http.StatusBadRequest,
			verify:         nil,
		},
		{
			name: "Create Link with Invalid Destination",
			input: models.LinkInput{
				Keyword:     "invalid",
				Destination: "not-a-url",
			},
			expectedStatus: http.StatusBadRequest,
			verify:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/links", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.verify != nil {
				tt.verify(t, db, tt.input)
			}
		})
	}
}

func TestAppController_GetLink(t *testing.T) {
	db := setupTestDB(t)
	controller := NewAppController(db)
	r := chi.NewRouter()
	r.Get("/links/{id}", controller.GetLink)

	// Seed DB
	link := models.Link{
		Keyword:     "findme",
		Destination: "https://found.com",
	}
	db.Db.Create(&link)

	tests := []struct {
		name           string
		linkID         string
		expectedStatus int
		expectedDest   string
	}{
		{
			name:           "Get Existing Link",
			linkID:         "1", // ID should be 1 as it's the first inserted
			expectedStatus: http.StatusOK,
			expectedDest:   "https://found.com",
		},
		{
			name:           "Get Non-Existent Link",
			linkID:         "999",
			expectedStatus: http.StatusNotFound,
			expectedDest:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/links/"+tt.linkID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var responseLink models.Link
				json.NewDecoder(w.Body).Decode(&responseLink)
				assert.Equal(t, tt.expectedDest, responseLink.Destination)
			}
		})
	}
}

func TestAppController_GetLinks(t *testing.T) {
	db := setupTestDB(t)
	controller := NewAppController(db)
	r := chi.NewRouter()
	r.Get("/links", controller.GetLinks)

	// Seed DB
	db.Db.Create(&models.Link{Keyword: "a", Destination: "https://a.com"})
	db.Db.Create(&models.Link{Keyword: "b", Destination: "https://b.com"})

	req, _ := http.NewRequest("GET", "/links", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var links []models.Link
	json.NewDecoder(w.Body).Decode(&links)
	assert.Len(t, links, 2)
}

func TestAppController_UpdateLink(t *testing.T) {
	db := setupTestDB(t)
	controller := NewAppController(db)
	r := chi.NewRouter()
	r.Patch("/links/{id}", controller.UpdateLink)

	// Seed DB
	link := models.Link{
		Keyword:     "update-me",
		Destination: "https://old.com",
		Description: "Old Description",
	}
	db.Db.Create(&link)

	tests := []struct {
		name           string
		linkID         string
		input          models.LinkInput
		expectedStatus int
		verify         func(*testing.T, *database.DbConnection)
	}{
		{
			name:   "Update Link Successfully",
			linkID: "1",
			input: models.LinkInput{
				Keyword:     "update-me",
				Destination: "https://new.com",
				Description: "New Description",
			},
			expectedStatus: http.StatusOK,
			verify: func(t *testing.T, db *database.DbConnection) {
				var updated models.Link
				err := db.Db.Where("keyword = ?", "update-me").First(&updated).Error
				assert.NoError(t, err)
				assert.Equal(t, "https://new.com", updated.Destination)
				assert.Equal(t, "New Description", updated.Description)
			},
		},
		{
			name:   "Update Non-Existent Link",
			linkID: "999",
			input: models.LinkInput{
				Keyword:     "test",
				Destination: "https://example.com",
			},
			expectedStatus: http.StatusNotFound,
			verify:         nil,
		},
		{
			name:   "Update Link with Invalid Destination",
			linkID: "1",
			input: models.LinkInput{
				Keyword:     "update-me",
				Destination: "not-a-url",
			},
			expectedStatus: http.StatusBadRequest,
			verify:         nil,
		},
		{
			name:   "Update Link with Redirect Loop",
			linkID: "1",
			input: models.LinkInput{
				Keyword:     "update-me",
				Destination: "http://localhost:8998/redirect",
			},
			expectedStatus: http.StatusBadRequest,
			verify:         nil,
		},
		{
			name:   "Update Link Starting with Slash",
			linkID: "1",
			input: models.LinkInput{
				Keyword:     "/slash",
				Destination: "https://example.com",
			},
			expectedStatus: http.StatusBadRequest,
			verify:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("PATCH", "/links/"+tt.linkID, bytes.NewBuffer(body))
			// Set Host header for redirect loop test
			req.Host = "localhost:8998"
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.verify != nil {
				tt.verify(t, db)
			}
		})
	}
}

func TestAppController_DeleteLink(t *testing.T) {
	db := setupTestDB(t)
	controller := NewAppController(db)
	r := chi.NewRouter()
	r.Delete("/links/{id}", controller.DeleteLink)

	tests := []struct {
		name           string
		linkID         string
		expectedStatus int
		verify         func(*testing.T, *database.DbConnection)
	}{
		{
			name:           "Delete Existing Link",
			linkID:         "1",
			expectedStatus: http.StatusOK,
			verify: func(t *testing.T, db *database.DbConnection) {
				var link models.Link
				// Even with Unscoped(), the record should be permanently deleted
				err := db.Db.Unscoped().Where("id = ?", 1).First(&link).Error
				assert.Error(t, err) // Should not find the deleted link
			},
		},
		// Note: GORM's Delete doesn't return error for non-existent records,
		// it just affects 0 rows. Testing this would require checking RowsAffected.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Seed DB for each test
			if tt.linkID == "1" {
				db.Db.Exec("DELETE FROM links") // Clean first
				db.Db.Create(&models.Link{
					Keyword:     "delete-me",
					Destination: "https://delete.com",
				})
			}

			req, _ := http.NewRequest("DELETE", "/links/"+tt.linkID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.verify != nil {
				tt.verify(t, db)
			}
		})
	}
}

func TestAppController_SearchLinks(t *testing.T) {
	db := setupTestDB(t)
	controller := NewAppController(db)
	r := chi.NewRouter()
	r.Get("/api/v1/search", controller.SearchLinks)

	// Seed DB
	db.Db.Create(&models.Link{Keyword: "github", Destination: "https://github.com"})
	db.Db.Create(&models.Link{Keyword: "gitlab", Destination: "https://gitlab.com"})
	db.Db.Create(&models.Link{Keyword: "google", Destination: "https://google.com"})
	db.Db.Create(&models.Link{Keyword: "go-docs", Destination: "https://golang.org"})

	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		expectedCount  int
		checkKeywords  []string
	}{
		{
			name:           "Search with Prefix 'git'",
			queryParam:     "git",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			checkKeywords:  []string{"github", "gitlab"},
		},
		{
			name:           "Search with Prefix 'go'",
			queryParam:     "go",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			checkKeywords:  []string{"google", "go-docs"},
		},
		{
			name:           "Search with No Matches",
			queryParam:     "xyz",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
			checkKeywords:  []string{},
		},
		{
			name:           "Search with Empty Query",
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			checkKeywords:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/search"
			if tt.queryParam != "" {
				url += "?qs=" + tt.queryParam
			} else {
				// Test without query param at all
				url += "?other=value"
			}

			req, _ := http.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var links []models.Link
				json.NewDecoder(w.Body).Decode(&links)
				assert.Len(t, links, tt.expectedCount)

				// Verify keywords if specified
				if len(tt.checkKeywords) > 0 {
					keywords := make([]string, len(links))
					for i, link := range links {
						keywords[i] = link.Keyword
					}
					for _, expected := range tt.checkKeywords {
						assert.Contains(t, keywords, expected)
					}
				}
			}
		})
	}
}
