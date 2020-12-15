package model

type User struct {
	Model
	Phone   int      `gorm:"not null;unique;type:bigint"`
	Email   *string  `gorm:"unique"`
	Servers []Server `gorm:"ForeignKey:BindUser"`
}
