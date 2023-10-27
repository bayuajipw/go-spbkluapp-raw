package authmiddleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("mysession"))

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "mysession")
		authenticated := session.Values["authenticated"]

		if authenticated == true {
			h.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	})
}
