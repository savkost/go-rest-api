package models

type LicenseCategory struct {
	ID                     string `bson:"_id,omitempty" json:"_id"`
	Title                  string `bson:"title,omitempty" json:"title"`
	Description            string `bson:"description,omitempty" json:"description"`
	PriceEurosMonthly      int64  `bson:"priceEurosMonthly,omitempty" json:"priceEurosMonthly"`
	PriceEurosThreeMonths  int64  `bson:"priceEurosThreeMonths,omitempty" json:"priceEurosThreeMonths"`
	PriceEurosSixMonths    int64  `bson:"priceEurosSixMonths,omitempty" json:"priceEurosSixMonths"`
	PriceEurosTwelveMonths int64  `bson:"priceEurosTwelveMonths,omitempty" json:"priceEurosTwelveMonths"`
	CategoryType           string `bson:"categoryType,omitempty" json:"categoryType"`
	Comments               string `bson:"comments,omitempty" json:"comments"`
	IsActive               string `bson:"isActive,omitempty" json:"isActive"`
	TextsQntAllowed        int64  `bson:"textsQntAllowed,omitempty" json:"textsQntAllowed"`
	ImagesQntAllowed       int64  `bson:"imagesQntAllowed,omitempty" json:"imagesQntAllowed"`
	CreatedDt              string `bson:"createdDt,omitempty" json:"createdDt"`
	LastUpdatedDt          string `bson:"lastUpdatedDt,omitempty" json:"lastUpdatedDt"`
}
