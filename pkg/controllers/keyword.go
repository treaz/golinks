package controllers

import (
	"crhuber/golinks/pkg/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (c *AppController) updateViewCount(link *models.Link) {
	link.Views += 1
	c.db.Db.Save(link)
}

func (c *AppController) GetKeyword(w http.ResponseWriter, r *http.Request) {
	link := models.Link{}
	keyword := chi.URLParam(r, "keyword")
	subkey := chi.URLParam(r, "subkey")
	wildcardPath := chi.URLParam(r, "*")

	// 1. Try exact match first (standard behavior)
	fullKeyword := keyword
	if subkey != "" {
		fullKeyword = fmt.Sprintf("%s/%s", keyword, subkey)
	}

	err := c.db.Db.First(&link, "keyword = ?", fullKeyword).Error
	if err == nil {
		// Found exact match. Check if we have extra path segments to substitute.
		var extraParts []string
		if wildcardPath != "" {
			extraParts = strings.Split(wildcardPath, "/")
		}

		// If we have extra parts, attempt {*} substitution
		if len(extraParts) > 0 {
			destination := link.Destination
			for _, part := range extraParts {
				// Replace the first occurrence of {*} with the part
				// If no {*} exists, this does nothing
				destination = strings.Replace(destination, "{*}", part, 1)
			}
			go c.updateViewCount(&link)
			http.Redirect(w, r, destination, http.StatusTemporaryRedirect)
			return
		}

		// No extra parts - redirect to destination as-is
		go c.updateViewCount(&link)
		http.Redirect(w, r, link.Destination, http.StatusTemporaryRedirect)
		return
	}

	// 2. Base keyword match - check if extra path segments should be treated as variables
	// We look up by the base keyword (the first segment)
	err = c.db.Db.First(&link, "keyword = ?", keyword).Error
	if err == nil {
		// Found base keyword. Collect all extra path segments.
		var extraParts []string
		if subkey != "" {
			extraParts = append(extraParts, subkey)
		}
		if wildcardPath != "" {
			extraParts = append(extraParts, strings.Split(wildcardPath, "/")...)
		}

		// If we have extra parts, attempt {*} substitution
		if len(extraParts) > 0 {
			destination := link.Destination
			for _, part := range extraParts {
				// Replace the first occurrence of {*} with the part
				// If no {*} exists in destination, this does nothing
				destination = strings.Replace(destination, "{*}", part, 1)
			}
			go c.updateViewCount(&link)
			http.Redirect(w, r, destination, http.StatusTemporaryRedirect)
			return
		}

		// No extra parts - redirect to destination as-is
		go c.updateViewCount(&link)
		http.Redirect(w, r, link.Destination, http.StatusTemporaryRedirect)
		return
	}

	// 3. Not found, redirect to search
	http.Redirect(w, r, fmt.Sprintf("/?q=%s", keyword), http.StatusTemporaryRedirect)
}
