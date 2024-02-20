package authmiddleware

import (
	"net/http"
	"spbkluapp/config"
	"time"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("mysession"))

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "mysession")
		authenticated := session.Values["authenticated"]

		if authenticated == true {
			// Check last login
			email, ok := session.Values["email"].(string)
			if !ok {
				session.Options.MaxAge = -1
				session.Save(r, w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			lastLogin, err := getLastLoginFromDatabase(email)
			if err != nil {
				http.Error(w, "Error checking last login", http.StatusInternalServerError)
				return
			}

			// Check if last login is more than 60 minutes ago
			if time.Since(lastLogin).Minutes() > 60 {
				// Clear session
				session.Options.MaxAge = -1
				session.Save(r, w)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Update last login in the session
			session.Values["last_login"] = time.Now()
			session.Save(r, w)

			h.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})
}

func getLastLoginFromDatabase(email string) (time.Time, error) {
	var lastLogin time.Time

	err := config.DB.QueryRow("SELECT last_login FROM users WHERE email = ?", email).Scan(&lastLogin)
	if err != nil {
		return time.Time{}, err
	}

	return lastLogin, nil
}
