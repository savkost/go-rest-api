package models

type License struct {
	ID                string `bson:"_id,omitempty" json:"_id"`
	LicenseKey        string `bson:"licenseKey,omitempty" json:"licenseKey"`
	Begin_dt          string `bson:"begin_dt,omitempty" json:"begin_dt"`
	Expiration_dt     string `bson:"expiration_dt,omitempty" json:"expiration_dt"`
	UserHolderId      string `bson:"userHolderId,omitempty" json:"userHolderId"`
	UserFullName      string `bson:"userFullName,omitempty" json:"userFullName"`
	CategoryId        string `bson:"categoryId,omitempty" json:"categoryId"`
	CategoryType      string `bson:"categoryType,omitempty" json:"categoryType"`
	CategoryTitle     string `bson:"categoryTitle,omitempty" json:"categoryTitle"`
	ActivatedOnDevice string `bson:"activatedOnDevice,omitempty" json:"activatedOnDevice"`
	TimeSpanType      int64  `bson:"timeSpanType,omitempty" json:"timeSpanType"`
	Comments          string `bson:"comments,omitempty" json:"comments"`
	IsActive          string `bson:"isActive,omitempty" json:"isActive"`
	IsExpired         string `bson:"isExpired,omitempty" json:"isExpired"`
	CreatedDt         string `bson:"createdDt,omitempty" json:"createdDt"`
	LastUpdatedDt     string `bson:"lastUpdatedDt,omitempty" json:"lastUpdatedDt"`
}
