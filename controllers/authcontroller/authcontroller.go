package authcontroller

import (
	"errors"
	"html/template"
	"net/http"
	"spbkluapp/config"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var temp *template.Template
var store = sessions.NewCookieStore([]byte("mysession"))

func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	authenticated := session.Values["authenticated"]

	if authenticated == true {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	} else {
		// http.Redirect(w, r, "/", http.StatusSeeOther)
		temp, err := template.ParseFiles("views/auth/index.html", "views/header.html", "views/footer.html", "views/footerjs.html") // display multiple file

		if err != nil {
			panic(err)
		}

		temp.ExecuteTemplate(w, "index.html", nil)
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var hashedPassword string
	var role int8
	var name string

	rows, err := config.DB.Query("SELECT password, role, name FROM users WHERE email = ?", email)
	if err != nil {
		data := map[string]interface{}{
			"message": "Login error!",
		}
		temp, _ := template.ParseFiles("views/auth/index.html", "views/header.html", "views/footerjs.html") // display multiple file
		temp.ExecuteTemplate(w, "index.html", data)
	}

	defer rows.Close() // Close the result set when done

	// Iterate through the result set (assuming only one row is expected)
	for rows.Next() {
		err = rows.Scan(&hashedPassword, &role, &name)
		if err != nil {
			data := map[string]interface{}{
				"message": "Login error!",
			}
			temp, _ := template.ParseFiles("views/auth/index.html", "views/header.html", "views/footerjs.html") // display multiple file
			temp.ExecuteTemplate(w, "index.html", data)
		}
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	// if email == "email@email.com" && password == "12345" {
	if err == nil {
		session, _ := store.Get(r, "mysession")
		session.Values["email"] = email
		session.Values["role"] = role
		session.Values["name"] = name
		session.Values["flashMessage"] = "Login Success!"
		session.Values["showFlashMessage"] = true
		session.Values["authenticated"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		data := map[string]interface{}{
			"message": "Login Failed! Username or password doesn't match",
		}
		temp, _ := template.ParseFiles("views/auth/index.html", "views/header.html", "views/footerjs.html") // display multiple file
		temp.ExecuteTemplate(w, "index.html", data)
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type UserData struct {
	Username string
	Email    string
	Role     int8
}

func GetUserDataFromSession(r *http.Request) (UserData, error) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		return UserData{}, err
	}

	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return UserData{}, errors.New("User is not authenticated")
	}

	// Access user data from session values
	username, _ := session.Values["username"].(string)
	email, _ := session.Values["email"].(string)
	role, _ := session.Values["role"].(int8)

	userData := UserData{
		Username: username,
		Email:    email,
		Role:     role,
	}

	return userData, nil
}

// how to use

// func MyHandler(w http.ResponseWriter, r *http.Request) {
//     userData, err := GetUserDataFromSession(r)
//     if err != nil {
//         // Handle the case where the user is not authenticated
//         // You can redirect them to the login page or show an error message
//         http.Redirect(w, r, "/login", http.StatusSeeOther)
//         return
//     }

//     // Use the user data in your handler
//     // Example: userData.Username, userData.Email, userData.Role
// }
