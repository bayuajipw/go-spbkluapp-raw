package bssmodel

import (
	"database/sql"
	"spbkluapp/config"
	"spbkluapp/entities"
	"time"

	"github.com/go-sql-driver/mysql"
)

func formatTimestamp(nt mysql.NullTime, layout string) string {
	if nt.Valid {
		return nt.Time.Format(layout)
	}
	return "NULL"
}

func GetAll() []entities.BssRes {
	rows, err := config.DB.Query("Select bss.*, users.email from bss left join users on bss.user_active = users.id")

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
		var createdat, updatedat time.Time
		var lastping mysql.NullTime
		var useractive int8
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
		bss.UserActive = useractive

		var qrRes string
		if qrcode.Valid {
			qrRes = qrcode.String
		} else {
			qrRes = "NULL"
		}
		bss.QrCode = qrRes

		var transactionRes string
		if transactionid.Valid {
			transactionRes = transactionid.String
		} else {
			transactionRes = "NULL"
		}
		bss.TransactionId = transactionRes

		bss.LastPing = formatTimestamp(lastping, "2006-01-02 15:04:05")
		bss.CreatedAt = createdat.Format("2006-01-02 15:04:05")
		bss.UpdatedAt = updatedat.Format("2006-01-02 15:04:05")

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
