
type User struct {
	ID        uint `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Password  string
	UsmPesos  int
}

type Book struct {
	ID              uint `gorm:"primaryKey"`
	BookName        string
	BookCategory    string
	TransactionType string
	Price           int
	Status          string
	PopularityScore int
	Inventory       Inventory `gorm:"foreignKey:BookID"`
}

type Inventory struct {
	BookID            uint `gorm:"primaryKey"`
	AvailableQuantity int
}

type Loan struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	BookID     uint
	StartDate  string
	DueDate    string
	ReturnDate *string
	Status     string
}

type Sale struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	BookID   uint
	SaleDate string
}

type Transaction struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	BookID *uint
	Type   string
	Date   string
	Amount int
}
