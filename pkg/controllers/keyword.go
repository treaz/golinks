package controllers

import (
	"crhuber/golinks/pkg/models"
	"fmt"
	"log/slog"
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
		// Found exact match. Check if it's parameterized and needs substitution.
		if link.IsParameterized {
			var extraParts []string
			// If we matched "keyword/subkey" exactly, subkey is part of the keyword, so we don't treat it as an arg?
			// PROPOSAL: If the link is "kibana", and we matched "kibana", then any EXTRA parts are args.
			// If we matched "kibana/foo" (fullKeyword), and the LINK matches "kibana/foo", then wildcard is the only arg source.
			// But if we matched "kibana" (because subkey was empty), then wildcard is arg.

			if wildcardPath != "" {
				extraParts = strings.Split(wildcardPath, "/")
			}

			// Be careful: if fullKeyword included subkey, we consumed it.
			// But if fullKeyword == keyword (subkey empty), then we only look at wildcard.
			// What if link is "kibana/foo" and we call "kibana/foo/bar"?
			// keyword="kibana", subkey="foo", wildcard="bar".
			// fullKeyword="kibana/foo". Found.
			// extraParts should be ["bar"].

			// Logic for parameterized links: go/meetwith/tammy -> https://g.co/meet/{*} -> https://g.co/meet/tammy
			// Supports multiple variables: go/gh/repo/123 -> https://github.com/product/{*}/issues/{*} -> https://github.com/product/repo/issues/123
			destination := link.Destination
			for _, part := range extraParts {
				// Replace the first occurrence of {*} with the part
				destination = strings.Replace(destination, "{*}", part, 1)
			}
			go c.updateViewCount(&link)
			http.Redirect(w, r, destination, http.StatusTemporaryRedirect)
			return
		}

		go c.updateViewCount(&link)
		http.Redirect(w, r, link.Destination, http.StatusTemporaryRedirect)
		return
	}

	// 2. Not found, check if it's a parameterized link or a "programmatic" link (legacy wildcard)
	// We look up by the base keyword (the first segment)
	err = c.db.Db.First(&link, "keyword = ?", keyword).Error
	if err == nil {
		if link.IsParameterized {
			// Logic for parameterized links where exact match FAILED (e.g. go/kibana/qa where subkey=qa)
			// Reconstruct the full path from subkey and wildcard
			// "subkey" might be caught by the /{keyword}/{subkey} route or /{keyword}/*
			// We need to gather all parts after the keyword.
			var extraParts []string
			if subkey != "" {
				extraParts = append(extraParts, subkey)
			}
			if wildcardPath != "" {
				// wildcardPath can be "qa/fra"
				extraParts = append(extraParts, strings.Split(wildcardPath, "/")...)
			}

			destination := link.Destination
			for _, part := range extraParts {
				// Replace the first occurrence of {*} with the part
				destination = strings.Replace(destination, "{*}", part, 1)
			}

			go c.updateViewCount(&link)
			http.Redirect(w, r, destination, http.StatusTemporaryRedirect)
			return
		}

		// Fallthrough for non-parameterized links found by base keyword logic
		// Fallthrough: Link found but not parameterized.
		// If the user visited "go/keyword/foo", but "go/keyword" points to google.com, we probably shouldn't do anything special unless we want to support subpath appending?
		// Existing logic seemed to want to support "programmatic" links via `{*}`, let's keep that check below if we want to support `keyword/{*}` stored in DB.
	}

	// 3. Fallback: check for "programmatic" links stored with `/{*}` suffix (Legacy support from existing code)
	// The existing code was looking for `keyword/{*}` in the DB.
	// We need to re-construct what the "keyword" part is if the URL was `go/foo/bar`.
	// The existing code did splits.

	slog.Info("keyword not found in exact match. trying wildcard pattern")
	keywordParts := strings.Split(fullKeyword, "/")
	if len(keywordParts) > 0 {
		err := c.db.Db.First(&link, "keyword = ?", fmt.Sprintf("%s/{*}", keywordParts[0])).Error
		if err == nil {
			// Found a programmatic link `foo/{*}`
			// We need to extract the part that matched `{*}`
			// If matching `foo/{*}`, and we have `go/foo/bar/baz`, keywordParts are [foo, bar, baz]
			// The part replacing {*} is "bar/baz" (joined)

			if len(keywordParts) > 1 {
				replacement := strings.Join(keywordParts[1:], "/")
				programmaticDestination := strings.Replace(link.Destination, "{*}", replacement, 1)
				go c.updateViewCount(&link)
				http.Redirect(w, r, programmaticDestination, http.StatusTemporaryRedirect)
				return
			}
		}
	}

	// 4. Truly not found, redirect to search
	http.Redirect(w, r, fmt.Sprintf("/?q=%s", keyword), http.StatusTemporaryRedirect)
}
