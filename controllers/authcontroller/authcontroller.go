package authcontroller

import (
	"errors"
	"html/template"
	"net/http"
	"spbkluapp/config"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var temp *template.Template
var store = sessions.NewCookieStore([]byte("mysession"))
var role_name string

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

	var id int8
	var name, hashedPassword string
	var role int8

	rows, err := config.DB.Query("SELECT id, password, role, name FROM users WHERE email = ? && account_active = 1 && role < 4", email)
	if err != nil {
		data := map[string]interface{}{
			"message": "Login error! " + err.Error(),
		}
		temp, _ := template.ParseFiles("views/auth/index.html", "views/header.html", "views/footerjs.html") // display multiple file
		temp.ExecuteTemplate(w, "index.html", data)
	}

	defer rows.Close() // Close the result set when done

	// Iterate through the result set (assuming only one row is expected)
	if rows.Next() {
		err = rows.Scan(&id, &hashedPassword, &role, &name)
		if err != nil {
			data := map[string]interface{}{
				"message": "Login error! " + err.Error(),
			}
			temp, _ := template.ParseFiles("views/auth/index.html", "views/header.html", "views/footerjs.html") // display multiple file
			temp.ExecuteTemplate(w, "index.html", data)
		} else {
			// Compare the hashed password with the provided password
			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

			// if email == "email@email.com" && password == "12345" {
			if err == nil {

				rows, err_login := config.DB.Query("UPDATE users SET last_login = ? WHERE email = ?", time.Now(), email)
				if err_login == nil {
					session, _ := store.Get(r, "mysession")
					session.Values["id"] = id
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
						"message": "Login Failed! Can't connect to database",
					}
					temp, _ := template.ParseFiles("views/auth/index.html", "views/header.html", "views/footerjs.html") // display multiple file
					temp.ExecuteTemplate(w, "index.html", data)
				}
				defer rows.Close()
			} else {
				data := map[string]interface{}{
					"message": "Login Failed! Username or password doesn't match",
				}
				temp, _ := template.ParseFiles("views/auth/index.html", "views/header.html", "views/footerjs.html") // display multiple file
				temp.ExecuteTemplate(w, "index.html", data)
			}
		}
	} else {
		data := map[string]interface{}{
			"message": "User not found!",
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
	Id       int8
	Username string
	Email    string
	Role     int8
}

func GetUserDataFromSession(r *http.Request) (UserData, error) {
	session, err := store.Get(r, "mysession")
	if err != nil {
		return UserData{}, err
	}

	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return UserData{}, errors.New("User is not authenticated")
	}

	// Access user data from session values
	id, _ := session.Values["id"].(int8)
	username, _ := session.Values["name"].(string)
	email, _ := session.Values["email"].(string)
	role, _ := session.Values["role"].(int8)

	userData := UserData{
		Id:       id,
		Username: username,
		Email:    email,
		Role:     role,
	}

	return userData, nil
}

func IndexHome(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	role, _ := session.Values["role"].(int8)
	username, _ := session.Values["name"].(string)

	// Check if a flash message should be displayed
	if showFlashMessage, ok := session.Values["showFlashMessage"].(bool); ok && showFlashMessage {
		flashMessage, _ := session.Values["flashMessage"].(string)

		if role == 0 {
			role_name = "Superadmin"
		} else if role == 1 {
			role_name = "Operator"
		} else if role == 2 {
			role_name = "Admin Mobile"
		} else if role == 3 {
			role_name = "Guest"
		} else {
			role_name = "Customer"
		}

		// Remove the flash message from the session
		delete(session.Values, "flashMessage")
		delete(session.Values, "showFlashMessage")
		session.Save(r, w)

		data := map[string]interface{}{
			"ShowFlashMessage": showFlashMessage,
			"FlashMessage":     flashMessage + " Welcome, " + role_name + "!",
			"username":         username,
			"role":             role,
		}

		// Display the "welcome.html" page with the flash message
		temp, err := template.ParseFiles("views/home/index.html", "views/header.html", "views/sidebar.html", "views/navbar.html", "views/footer.html", "views/footerjs.html") // display multiple file

		if err != nil {
			panic(err)
		}

		// temp.Execute(w, nil)
		temp.ExecuteTemplate(w, "index.html", data)
	} else {
		data := map[string]interface{}{
			"username": username,
			"role":     role,
		}
		temp, err := template.ParseFiles("views/home/index.html", "views/header.html", "views/sidebar.html", "views/navbar.html", "views/footer.html", "views/footerjs.html") // display multiple file

		if err != nil {
			panic(err)
		}

		// temp.Execute(w, nil)
		temp.ExecuteTemplate(w, "index.html", data)
	}

}
