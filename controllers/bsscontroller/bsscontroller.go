package bsscontroller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"spbkluapp/config"
	"spbkluapp/models/bssmodel"
	"time"
)

var temp *template.Template

func respondWithError(w http.ResponseWriter, message string) {
	respondWithJSON(w, map[string]interface{}{"success": false, "message": message})
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func Index(w http.ResponseWriter, r *http.Request) {
	bsslist := bssmodel.GetAll()
	data := map[string]any{
		"bsslist": bsslist,
	}

	temp, err := template.ParseFiles("views/bss/index.html", "views/header.html", "views/sidebar.html", "views/navbar.html", "views/footer.html", "views/footerjs.html") // display multiple file

	if err != nil {
		panic(err)
	}

	// temp.Execute(w, nil)
	temp.ExecuteTemplate(w, "index.html", data)
}

func Get(w http.ResponseWriter, r *http.Request) {
	// Generate or fetch updated data (replace with your logic)
	// updatedData := generateUpdatedData()

	updatedData := bssmodel.GetAll()

	// Convert the data to JSON
	jsonData, err := json.Marshal(updatedData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON data", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(jsonData)
}

func Add(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve form values
	name := r.FormValue("name")
	address := r.FormValue("address")
	city := r.FormValue("city")
	province := r.FormValue("province")
	longitude := r.FormValue("longitude")
	latitude := r.FormValue("latitude")
	slot := r.FormValue("slot")

	// Perform server-side validation here
	if name == "" || address == "" || city == "" || province == "" || longitude == "" || latitude == "" || slot == "" {
		respondWithError(w, "All fields are required")
		return
	}

	rows, _ := config.DB.Query("Select * from bss where name = ?", name)

	var rowCount int // Initialize a counter variable

	for rows.Next() {
		rowCount++ // Increment the counter for each row
		// You can optionally process the row data here using rows.Scan()
	}

	if rowCount > 0 {
		respondWithError(w, "BSS name is already use, please use another name!")
		return
	}

	_, err := config.DB.Exec("INSERT INTO bss (name, address, city, province, longitude, latitude, slot, user_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", name, address, city, province, longitude, latitude, slot, 0, time.Now(), time.Now())
	if err != nil {
		respondWithError(w, "Error adding data to the database")
		return
	}

	// Data added successfully
	respondWithJSON(w, map[string]interface{}{"success": true})
}

func Edit(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve form values
	id := r.FormValue("rowId")
	name := r.FormValue("name")
	address := r.FormValue("address")
	city := r.FormValue("city")
	province := r.FormValue("province")
	longitude := r.FormValue("longitude")
	latitude := r.FormValue("latitude")
	slot := r.FormValue("slot")
	status := r.FormValue("status")

	// Perform server-side validation here
	if name == "" || address == "" || city == "" || province == "" || longitude == "" || latitude == "" || slot == "" || status == "" {
		respondWithError(w, "All fields are required")
		return
	}

	rows, _ := config.DB.Query("Select * from bss where name = ? and id <> ?", name, id)

	var rowCount int // Initialize a counter variable

	for rows.Next() {
		rowCount++ // Increment the counter for each row
		// You can optionally process the row data here using rows.Scan()
	}

	if rowCount > 0 {
		respondWithError(w, "BSS name is already use, please use another name!")
		return
	}

	_, err := config.DB.Exec("Update bss set name = ?, address = ?, city = ?, province = ?, longitude = ?, latitude = ?, status = ?, slot = ?, updated_at = ? where id = ?", name, address, city, province, longitude, latitude, status, slot, time.Now(), id)
	if err != nil {
		respondWithError(w, "Error Updating data to the database")
		return
	}

	// Data added successfully
	respondWithJSON(w, map[string]interface{}{"success": true})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")

	// Execute the DELETE query
	_, err := config.DB.Exec("DELETE FROM bss WHERE id = ?", id)

	if err != nil {
		respondWithError(w, "Error deleting data !")
		return
	}

	respondWithJSON(w, map[string]interface{}{"success": true})
}
