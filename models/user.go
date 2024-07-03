package models

type User struct {
	ID            string `bson:"_id,omitempty" json:"_id"`
	FirstName     string `bson:"firstName,omitempty" json:"firstName"`
	LastName      string `bson:"lastName,omitempty" json:"lastName"`
	Role          string `bson:"role,omitempty" json:"role"`
	IsAdmin       string `bson:"isAdmin,omitempty" json:"isAdmin"`
	IsActive      string `bson:"isActive,omitempty" json:"isActive"`
	Email         string `bson:"email,omitempty" json:"email"`
	Password      string `bson:"password,omitempty" json:"password"`
	CreatedDt     string `bson:"createdDt,omitempty" json:"createdDt"`
	LastUpdatedDt string `bson:"lastUpdatedDt,omitempty" json:"lastUpdatedDt"`
}
