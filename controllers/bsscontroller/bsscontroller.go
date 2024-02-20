package bsscontroller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"spbkluapp/config"
	"spbkluapp/controllers/authcontroller"
	"spbkluapp/controllers/maincontroller"
	"spbkluapp/entities"
	"spbkluapp/models/bssmodel"
	"strconv"
	"time"
)

var temp *template.Template

func Index(w http.ResponseWriter, r *http.Request) {
	// bsslist := bssmodel.GetAll()
	// data := map[string]any{
	// 	"bsslist": bsslist,
	// }
	userData, _ := authcontroller.GetUserDataFromSession(r)

	data := map[string]interface{}{
		"username": userData.Username,
		"role":     userData.Role,
	}

	temp, err := template.ParseFiles("views/bss/index.html", "views/header.html", "views/sidebar.html", "views/navbar.html", "views/footer.html", "views/footerjs.html") // display multiple file

	if err != nil {
		panic(err)
	}

	// temp.Execute(w, nil)
	temp.ExecuteTemplate(w, "index.html", data)
}

func Get(w http.ResponseWriter, r *http.Request) {
	userData, _ := authcontroller.GetUserDataFromSession(r)
	role := userData.Role
	updatedData := bssmodel.GetAll()

	data := map[string]interface{}{
		"data": updatedData,
		"role": role,
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
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
	name := r.FormValue("add-name")
	address := r.FormValue("add-address")
	city := r.FormValue("add-city")
	province := r.FormValue("add-province")
	longitude := r.FormValue("add-longitude")
	latitude := r.FormValue("add-latitude")
	slot := r.FormValue("add-slot")

	// Perform server-side validation here
	if name == "" || address == "" || city == "" || province == "" || longitude == "" || latitude == "" || slot == "" {
		maincontroller.RespondWithError(w, "All fields are required")
		return
	}

	strSlot, err := strconv.Atoi(slot)
	if err != nil {
		maincontroller.RespondWithError(w, "Invalid slot value")
		return
	}

	strLongitude, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		maincontroller.RespondWithError(w, "Invalid longitude value")
		return
	}

	strLatitude, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		maincontroller.RespondWithError(w, "Invalid latitude value")
		return
	}

	input := entities.BssReq{
		Name:      name,
		Address:   address,
		City:      city,
		Province:  province,
		Slot:      int8(strSlot),
		Longitude: strLongitude,
		Latitude:  strLatitude,
	}

	// Check if BSS name already exists
	rowsName, err := config.DB.Query("SELECT * FROM bss WHERE name = ?", input.Name)
	if err != nil {
		maincontroller.RespondWithError(w, "Error querying database: "+err.Error())
		return
	}
	defer rowsName.Close()

	if rowsName.Next() {
		maincontroller.RespondWithError(w, "BSS name is already in use, please choose another name!")
		return
	}

	// Insert into BSS table
	result, err := config.DB.Exec("INSERT INTO bss (name, address, city, province, longitude, latitude, slot, user_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		input.Name, input.Address, input.City, input.Province, input.Longitude, input.Latitude, input.Slot, 0, time.Now(), time.Now())
	if err != nil {
		maincontroller.RespondWithError(w, "Error adding data to the database: "+err.Error())
		return
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		maincontroller.RespondWithError(w, "Error getting last inserted ID: "+err.Error())
		return
	}

	// Insert into locker table within a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		maincontroller.RespondWithError(w, "Error starting transaction: "+err.Error())
		return
	}

	defer func() {
		if err != nil {
			// Rollback the transaction if an error occurred
			fmt.Println("Rolling back transaction:", tx.Rollback())
		}
	}()

	// Loop to insert rows
	for lockerNumber := 1; lockerNumber <= int(input.Slot); lockerNumber++ {
		lockerName := fmt.Sprintf("lk%d", lockerNumber)
		_, err := tx.Exec("INSERT INTO locker (bss_id, name, locker_number, status, swapping_out, soc, soh, temp, stay_hour, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			lastInsertID, lockerName, lockerNumber, 0, 0, 0, 0, 0, 0, time.Now(), time.Now())

		if err != nil {
			maincontroller.RespondWithError(w, "Error inserting into locker table: "+err.Error())
			return
		}
	}

	// Commit the transaction if all inserts succeeded
	if err := tx.Commit(); err != nil {
		maincontroller.RespondWithError(w, "Error committing transaction: "+err.Error())
		return
	}

	// Data added successfully
	maincontroller.RespondWithJSON(w, map[string]interface{}{"success": true})
}

func Edit(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve form values
	id := r.FormValue("edit-rowId")
	name := r.FormValue("edit-name")
	address := r.FormValue("edit-address")
	city := r.FormValue("edit-city")
	province := r.FormValue("edit-province")
	longitude := r.FormValue("edit-longitude")
	latitude := r.FormValue("edit-latitude")
	// slot := r.FormValue("slot")
	status := r.FormValue("status")

	// Perform server-side validation here
	if name == "" || address == "" || city == "" || province == "" || longitude == "" || latitude == "" || status == "" {
		maincontroller.RespondWithError(w, "All fields are required")
		return
	}

	// str_slot, _ := strconv.Atoi(slot)
	str_status, _ := strconv.Atoi(status)
	str_longitude, _ := strconv.ParseFloat(longitude, 64)
	str_latitude, _ := strconv.ParseFloat(latitude, 64)

	input := entities.BssReq{
		Name:     name,
		Address:  address,
		City:     city,
		Province: province,
		// Slot:      int8(str_slot),
		Status:    int8(str_status),
		Longitude: str_longitude,
		Latitude:  str_latitude,
	}

	rows, _ := config.DB.Query("SELECT * FROM bss WHERE name = ? AND id <> ?", input.Name, id)

	var rowCount int // Initialize a counter variable

	for rows.Next() {
		rowCount++ // Increment the counter for each row
		// You can optionally process the row data here using rows.Scan()
	}

	if rowCount > 0 {
		maincontroller.RespondWithError(w, "BSS name is already use, please use another name!")
		return
	}

	_, err := config.DB.Exec("UPDATE bss SET name = ?, address = ?, city = ?, province = ?, longitude = ?, latitude = ?, status = ?, updated_at = ? WHERE id = ?", input.Name, input.Address, input.City, input.Province, input.Longitude, input.Latitude, input.Status, time.Now(), id)
	if err != nil {
		maincontroller.RespondWithError(w, "Error Updating data to the database: "+err.Error())
		return
	}

	// Data added successfully
	maincontroller.RespondWithJSON(w, map[string]interface{}{"success": true})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id") // BSS ID

	// Check if the BSS ID is used in the transaction table
	rowsTransaction, err := config.DB.Query("SELECT COUNT(*) FROM transaction WHERE bss_id = ?", id)
	if err != nil {
		maincontroller.RespondWithError(w, "Error checking if BSS is used in transactions: "+err.Error())
		return
	}
	defer rowsTransaction.Close()

	var transactionCount int
	if rowsTransaction.Next() {
		if err := rowsTransaction.Scan(&transactionCount); err != nil {
			maincontroller.RespondWithError(w, "Error scanning database result: "+err.Error())
			return
		}
	}

	if transactionCount > 0 {
		maincontroller.RespondWithError(w, "Cannot delete BSS. It is used in transactions.")
		return
	}

	// Check if any locker associated with the BSS ID is used in the battery table
	rowsLocker, err := config.DB.Query("SELECT COUNT(*) FROM locker WHERE bss_id = ? AND locker.id IN (SELECT locker_id FROM battery)", id)
	if err != nil {
		maincontroller.RespondWithError(w, "Error checking if lockers are used: "+err.Error())
		return
	}
	defer rowsLocker.Close()

	var lockerCount int
	if rowsLocker.Next() {
		if err := rowsLocker.Scan(&lockerCount); err != nil {
			maincontroller.RespondWithError(w, "Error scanning database result: "+err.Error())
			return
		}
	}

	if lockerCount > 0 {
		maincontroller.RespondWithError(w, "Cannot delete BSS. At least one locker is used in the battery table.")
		return
	}

	// Begin a transaction to delete lockers and BSS
	tx, err := config.DB.Begin()
	if err != nil {
		maincontroller.RespondWithError(w, "Error starting transaction:"+err.Error())
		return
	}
	defer tx.Rollback()

	// Delete all lockers associated with the BSS ID
	_, err = tx.Exec("DELETE FROM locker WHERE bss_id = ?", id)
	if err != nil {
		maincontroller.RespondWithError(w, "Error deleting lockers: "+err.Error())
		return
	}

	// Delete the BSS
	_, err = tx.Exec("DELETE FROM bss WHERE id = ?", id)
	if err != nil {
		maincontroller.RespondWithError(w, "Error deleting BSS: "+err.Error())
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		maincontroller.RespondWithError(w, "Error committing transaction: "+err.Error())
		return
	}

	maincontroller.RespondWithJSON(w, map[string]interface{}{"success": true})
}
