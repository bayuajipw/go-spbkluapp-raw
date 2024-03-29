package entities

type BssReq struct {
	Id        uint8
	Name      string
	Address   string
	City      string
	Province  string
	Longitude float64
	Latitude  float64
	Slot      int8
	Status    int8
	// LastPing      mysql.NullTime
	// CreatedAt     time.Time
	// UpdatedAt     time.Time
	// UserActive    sql.NullInt16
	// TransactionId sql.NullString
	// QrCode        sql.NullString
	// Email         string
}

type BssRes struct {
	No            int
	Id            uint8
	Name          string
	Address       string
	City          string
	Province      string
	Longitude     float64
	Latitude      float64
	Slot          int8
	Status        int8
	LastPing      string
	CreatedAt     string
	UpdatedAt     string
	UserActive    int8
	TransactionId string
	QrCode        string
	Email         string
}
