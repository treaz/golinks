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
	db.Db.Exec("DELETE FROM tags")

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
