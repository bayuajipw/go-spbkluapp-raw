package homecontroller

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

var temp *template.Template
var store = sessions.NewCookieStore([]byte("mysession"))
var role_name string

func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")

	// Check if a flash message should be displayed
	if showFlashMessage, ok := session.Values["showFlashMessage"].(bool); ok && showFlashMessage {
		flashMessage, _ := session.Values["flashMessage"].(string)
		role, _ := session.Values["role"].(int8)

		if role == 0 {
			role_name = "Superadmin"
		} else if role == 1 {
			role_name = "Admin"
		} else {
			role_name = "User"
		}

		// Remove the flash message from the session
		delete(session.Values, "flashMessage")
		delete(session.Values, "showFlashMessage")
		session.Save(r, w)

		data := map[string]interface{}{
			"ShowFlashMessage": showFlashMessage,
			"FlashMessage":     flashMessage + " Welcome, " + role_name + "!",
		}

		// Display the "welcome.html" page with the flash message
		temp, err := template.ParseFiles("views/home/index.html", "views/header.html", "views/sidebar.html", "views/navbar.html", "views/footer.html", "views/footerjs.html") // display multiple file

		if err != nil {
			panic(err)
		}

		// temp.Execute(w, nil)
		temp.ExecuteTemplate(w, "index.html", data)
	} else {
		temp, err := template.ParseFiles("views/home/index.html", "views/header.html", "views/sidebar.html", "views/navbar.html", "views/footer.html", "views/footerjs.html") // display multiple file

		if err != nil {
			panic(err)
		}

		// temp.Execute(w, nil)
		temp.ExecuteTemplate(w, "index.html", nil)
	}

}
