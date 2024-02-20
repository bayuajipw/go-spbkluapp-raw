package bssmodel

import (
	"database/sql"
	"spbkluapp/config"
	"spbkluapp/entities"
	"spbkluapp/models/mainmodel"

	"github.com/go-sql-driver/mysql"
)

func GetAll() []entities.BssRes {
	rows, err := config.DB.Query("SELECT bss.id, bss.name, bss.address, bss.city, bss.province, bss.longitude, bss.latitude, bss.slot, bss.status, bss.last_ping, bss.created_at, bss.updated_at, bss.user_active, bss.transaction_id, bss.qrcode, users.email FROM bss LEFT JOIN users ON bss.user_active = users.id")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var bsslist []entities.BssRes
	var no = 0

	for rows.Next() {
		no = no + 1
		var bss entities.BssRes

		var id uint8
		var name, address, city, province string
		var longitude, latitude float64
		var slot, status int8
		// var createdat, updatedat time.Time
		var createdat, updatedat, lastping mysql.NullTime
		var useractive sql.NullInt16
		var transactionid sql.NullString
		var qrcode sql.NullString
		var email sql.NullString

		// err := rows.Scan(&bss.Id, &bss.Name, &bss.Address, &bss.City, &bss.Province, &bss.Longitude, &bss.Latitude, &bss.Slot, &bss.Status, &bss.LastPing, &bss.CreatedAt, &bss.UpdatedAt, &bss.UserActive, &bss.TransactionId)
		err := rows.Scan(&id, &name, &address, &city, &province, &longitude, &latitude, &slot, &status, &lastping, &createdat, &updatedat, &useractive, &transactionid, &qrcode, &email)

		bss.No = no
		bss.Id = id
		bss.Name = name
		bss.Address = address
		bss.City = city
		bss.Province = province
		bss.Longitude = longitude
		bss.Latitude = latitude
		bss.Slot = slot
		bss.Status = status
		bss.UserActive = mainmodel.HandleNullInt(useractive)
		bss.QrCode = mainmodel.HandleNullString(qrcode)
		bss.TransactionId = mainmodel.HandleNullString(transactionid)
		bss.LastPing = mainmodel.FormatTimestamp(lastping, "2006-01-02 15:04:05")
		bss.CreatedAt = mainmodel.FormatTimestamp(createdat, "2006-01-02 15:04:05")
		bss.UpdatedAt = mainmodel.FormatTimestamp(updatedat, "2006-01-02 15:04:05")

		var emailRes string
		if email.Valid {
			emailRes = email.String
		} else {
			emailRes = "NULL"
		}
		bss.Email = emailRes

		if err != nil {
			panic(err)
		}

		bsslist = append(bsslist, bss)
	}

	return bsslist
}
