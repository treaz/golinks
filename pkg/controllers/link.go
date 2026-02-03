package controllers

import (
	"crhuber/golinks/pkg/database"
	"crhuber/golinks/pkg/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

type AppController struct {
	db *database.DbConnection // save pointer to gormDB
}

// convienince method to return a pointer to a AppController
func NewAppController(db *database.DbConnection) *AppController {
	// takes in a db
	return &AppController{
		db: db,
	}
}

// anything of type link controller has these functions available
func (c *AppController) GetLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	link := models.Link{}
	err := c.db.Db.First(&link, id).Error
	if err != nil {
		JsonError(w, err, http.StatusNotFound, "Not Found")
		return
	}
	json.NewEncoder(w).Encode(link)
}

func (c *AppController) GetLinks(w http.ResponseWriter, r *http.Request) {
	var order string
	var sort string
	var qs models.QueryString
	inputOrder := r.URL.Query().Get("order")
	inputSort := r.URL.Query().Get("sort")

	if inputOrder == "" {
		order = "desc"
	}

	if inputSort == "" {
		sort = "updated_at"
	}

	qs.Sort = sort
	qs.Order = order

	err := qs.Validate()
	if err != nil {
		slog.Error("Error validating querystring fields")
		JsonError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	links := models.Links{}
	c.db.Db.Order(fmt.Sprintf("%s %s", qs.Sort, qs.Order)).Find(&links)
	JsonResponse(w, links)
}

func (c *AppController) CreateLink(w http.ResponseWriter, r *http.Request) {
	ic := models.LinkInput{}
	json.NewDecoder(r.Body).Decode(&ic)

	err := ic.Validate()
	if err != nil {
		slog.Error("error validating fields", slog.Any("error", err))
		JsonError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	// prevent redirection loop
	// Handle {*} in destination for parsing check
	safeDestination := strings.ReplaceAll(ic.Destination, "{*}", "placeholder")
	u, err := url.Parse(safeDestination)
	if err != nil {
		JsonError(w, err, http.StatusBadRequest, "error parsing destination")
		return
	}
	if u.Host == r.Host {
		JsonError(w, err, http.StatusBadRequest, "cannot create link with destination same as current host")
		return
	}
	// validate keyword
	if strings.HasPrefix(ic.Keyword, "/") {
		JsonError(w, err, http.StatusBadRequest, "cannot create link starting with slash")
		return
	}

	link := ic.ToNative()
	err = c.db.Db.Create(&link).Error
	if err != nil {
		slog.Error("error saving link to db")
		JsonError(w, err, http.StatusBadRequest, "error saving link to db. keywords must be unique")
		return
	}
	JsonResponse(w, link)
}

func (c *AppController) UpdateLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ic := models.LinkInput{}
	json.NewDecoder(r.Body).Decode(&ic)

	err := ic.Validate()
	if err != nil {
		slog.Error("error validating fields")
		JsonError(w, err, http.StatusBadRequest, err.Error())
		return
	}
	// prevent redirection loop
	// Handle {*} in destination for parsing check
	safeDestination := strings.ReplaceAll(ic.Destination, "{*}", "placeholder")
	u, err := url.Parse(safeDestination)
	if err != nil {
		JsonError(w, err, http.StatusBadRequest, "error parsing destination")
		return
	}
	if u.Host == r.Host {
		JsonError(w, err, http.StatusBadRequest, "cannot create link with destination same as current host")
		return
	}
	// validate keyword
	if strings.HasPrefix(ic.Keyword, "/") {
		JsonError(w, err, http.StatusBadRequest, "cannot create link starting with slash")
		return
	}
	oldLink := models.Link{}
	err = c.db.Db.First(&oldLink, id).Error
	if err != nil {
		JsonError(w, err, http.StatusNotFound, "Not Found")
		return
	}

	newLink := ic.ToNative()
	c.db.Db.Model(&oldLink).Updates(newLink)
	JsonResponse(w, oldLink)
}

func (c *AppController) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	link := models.Link{}
	err := c.db.Db.Unscoped().Delete(&link, id).Error
	if err != nil {
		// If we didn't find it, 404
		JsonError(w, err, http.StatusNotFound, "Not Found")
		return
	}
	JsonResponse(w, "Link Deleted")
}

func (c *AppController) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "OK"}); err != nil {
		panic(err)
	}
}

func (c *AppController) SearchLinks(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("qs")
	if input == "" {
		JsonError(w, nil, http.StatusBadRequest, "query string is required")
		return
	}
	links := models.Links{}
	c.db.Db.Limit(10).Where("keyword LIKE ?", fmt.Sprintf("%v%%", input)).Find(&links)
	JsonResponse(w, links)
}
