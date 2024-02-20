package mainmodel

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func FormatTimestamp(nt mysql.NullTime, layout string) string {
	if nt.Valid {
		return nt.Time.Format(layout)
	}
	return "NULL"
}

func HandleNullInt(nullableInt sql.NullInt16) int8 {
	if nullableInt.Valid {
		return int8(nullableInt.Int16)
	}
	return 0
}

func HandleNullString(nullableString sql.NullString) string {
	if nullableString.Valid {
		return nullableString.String
	}
	return "Null"
}

func HandleNullFloat(nullableFloat sql.NullFloat64) float64 {
	if nullableFloat.Valid {
		return nullableFloat.Float64
	}
	return 0.0
}

// formatCurrency formats a numerical value into a currency format
func FormatCurrency(amount float64) string {
	// Format the amount with comma as thousands separator and two decimal places
	formatted := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	var result string
	if len(integerPart) > 3 {
		result = integerPart[:len(integerPart)%3]
		for i := len(integerPart) % 3; i < len(integerPart); i += 3 {
			result += "." + integerPart[i:i+3]
		}
	} else {
		result = integerPart
	}

	return "Rp" + result + "," + decimalPart
}
