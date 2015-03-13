package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/hoisie/mustache"
	"mime"
	"net/http"
	"path"
	"strings"
)

var Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(16),
	securecookie.GenerateRandomKey(16))

// Control is an http.Handler which serves the web control panel.
type Control struct {
	Config *Config
}

// ServeAsset serves a static asset.
func (c Control) ServeAsset(w http.ResponseWriter, r *http.Request) {
	urlPath := path.Clean(r.URL.Path)
	if data, err := Asset(urlPath[1:]); err != nil {
		http.NotFound(w, r)
	} else {
		mimeType := mime.TypeByExtension(path.Ext(urlPath))
		if mimeType == "" {
			mimeType = "text/plain"
		}
		w.Header().Set("Content-Type", mimeType)
		w.Write(data)
	}
}

// ServeHTTP serves the web control panel.
func (c Control) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: serve pages here
	urlPath := path.Clean(r.URL.Path)
	if urlPath == "/login" {
		c.ServeLogin(w, r)
	} else if strings.HasPrefix(urlPath, "/assets/") {
		c.ServeAsset(w, r)
	} else if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	} else if urlPath == "/" {
		c.ServeRoot(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// ServeLogin serves the login page.
func (c Control) ServeLogin(w http.ResponseWriter, r *http.Request) {
	template := map[string]interface{}{"error": false}
	if r.Method == "POST" {
		// Get their submitted hash and the real hash.
		password := r.PostFormValue("password")
		hash := HashPassword(password)
		GlobalConfig.RLock()
		realHash := GlobalConfig.AdminHash
		GlobalConfig.RUnlock()
		// Check if they got the password correct.
		if hash == realHash {
			s, _ := Store.Get(r, "sessid")
			s.Values["authenticated"] = true
			s.Save(r, w)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		template["error"] = true
	}

	// Serve login page with no template.
	data, err := Asset("templates/login.mustache")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	content := mustache.Render(string(data), template)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}

// ServeRoot serves the homepage.
func (c Control) ServeRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("nothing here..."))
}

// HashPassword returns the SHA-256 hash of a string.
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

// IsAuthenticated returns whether or not a request was authenticated.
func IsAuthenticated(r *http.Request) bool {
	s, _ := Store.Get(r, "sessid")
	val, ok := s.Values["authenticated"].(bool)
	return ok && val
}
