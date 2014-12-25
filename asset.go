package goule

import (
	"net/http"
	pathlib "path"
	"regexp"
)

// Asset serves a static admin asset from an assets directory.
func Asset(assets string, w http.ResponseWriter, r *http.Request) {
	// The empty path redirects to /.
	if r.URL.Path == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// Validate the static file request
	path, ok := validatePath(r)
	if !ok {
		// TODO: send a nicer 404 page here
		http.NotFound(w, r)
		return
	}

	// Get the local path and send it
	localPath := pathlib.Join(assets, path)
	http.ServeFile(w, r, localPath)
}

func internalRedirect(r *http.Request) string {
	redirects := map[string]string{"/": "/index.html", "/login": "/login.html"}
	if newPath, ok := redirects[r.URL.Path]; ok {
		return newPath
	}
	return r.URL.Path
}

func validatePath(r *http.Request) (string, bool) {
	path := internalRedirect(r)

	// Valid types of paths: /style/*.css, /scripts/*.js, /images/*.png,
	// /*.html.
	charMatch := "[a-zA-Z0-9\\-_]*"
	htmlMatch := charMatch + "\\.html"
	cssMatch := "style\\/" + charMatch + "\\.css"
	scriptMatch := "scripts\\/" + charMatch + "\\.js"
	imageMatch := "images\\/" + charMatch + "\\.png"
	expr := "^\\/(" + htmlMatch + "|" + cssMatch + "|" + scriptMatch + "|" +
		imageMatch + ")$"
	if ok, _ := regexp.MatchString(expr, path); !ok {
		return "", false
	}
	return path, true
}