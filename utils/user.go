package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shreeharsha-factly/kavach-nedc-script/config"
)

type User struct {
	Base
	Email            string         `gorm:"column:email;uniqueIndex" json:"email"`
	KID              string         `gorm:"column:kid;" json:"kid"`
	FirstName        string         `gorm:"column:first_name" json:"first_name"`
	LastName         string         `gorm:"column:last_name" json:"last_name"`
	Slug             string         `gorm:"column:slug" json:"slug"`
	DisplayName      string         `gorm:"column:display_name" json:"display_name"`
	BirthDate        *time.Time     `gorm:"column:birth_date; default:NULL" json:"birth_date"`
	Gender           string         `gorm:"column:gender" json:"gender"`
	FeaturedMediumID *uint          `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	Medium           *Medium        `gorm:"foreignKey:featured_medium_id" json:"medium"`
	SocialMediaURLs  postgres.Jsonb `gorm:"column:social_media_urls" json:"social_media_urls" swaggertype:"primitive,string"`
	Description      string         `gorm:"column:description" json:"description"`
	Meta             postgres.Jsonb `gorm:"column:meta" json:"meta" swaggertype:"primitive,string"`
	IsActive         bool           `gorm:"column:is_active" json:"is_active" `
	Organisations    []Organisation `gorm:"many2many:organisation_users;" json:"organisations"`
}

type user struct {
	Base
	Email            string         `gorm:"column:email;uniqueIndex" json:"email"`
	KID              string         `gorm:"column:kid;" json:"kid"`
	FirstName        string         `gorm:"column:first_name" json:"first_name"`
	LastName         string         `gorm:"column:last_name" json:"last_name"`
	Slug             string         `gorm:"column:slug" json:"slug"`
	DisplayName      string         `gorm:"column:display_name" json:"display_name"`
	BirthDate        *time.Time     `gorm:"column:birth_date; default:NULL" json:"birth_date"`
	Gender           string         `gorm:"column:gender" json:"gender"`
	FeaturedMediumID *uint          `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	Medium           *Medium        `gorm:"foreignKey:featured_medium_id" json:"medium"`
	SocialMediaURLs  postgres.Jsonb `gorm:"column:social_media_urls" json:"social_media_urls" swaggertype:"primitive,string"`
	Description      string         `gorm:"column:description" json:"description"`
	Meta             postgres.Jsonb `gorm:"column:meta" json:"meta" swaggertype:"primitive,string"`
	IsActive         bool           `gorm:"column:is_active" json:"is_active" `
}

// TODO associate medium
func CreateUser() {
	data := make(map[uint]uint, 0)
	users := make([]user, 0)

	uids := make([]uint, 0)

	file, _ := ioutil.ReadFile("./uids.json")
	json.Unmarshal(file, &uids)

	config.ProdDB.Model(&user{}).Order("id asc").Find(&users)

	for _, ur := range users {
		if contains(uids, ur.ID) {
			u := User{
				Base: Base{
					CreatedAt: ur.CreatedAt,
					UpdatedAt: ur.UpdatedAt,
				},
				Email:           ur.Email,
				KID:             ur.KID,
				FirstName:       ur.FirstName,
				LastName:        ur.LastName,
				Slug:            ur.Slug,
				DisplayName:     ur.DisplayName,
				BirthDate:       ur.BirthDate,
				Gender:          ur.Gender,
				SocialMediaURLs: ur.SocialMediaURLs,
				Description:     ur.Description,
				Meta:            ur.Meta,
				IsActive:        ur.IsActive,
			}
			err := config.LocalDB.Model(&User{}).Order("id asc").Create(&u).Error

			log.Println("err")

			if err != nil {
				log.Println("ID", u.ID)
				log.Fatal(err)

				usersJson, _ := json.Marshal(data)
				err = ioutil.WriteFile("./users.json", usersJson, 0644)

				log.Println(data)
				if err != nil {
					log.Println("file error", err)
				}
			}
			data[ur.ID] = u.ID
		}
	}

	usersJson, _ := json.Marshal(data)
	err := ioutil.WriteFile("./users.json", usersJson, 0644)

	log.Println(data)
	if err != nil {
		log.Println("file error", err)
	}

}

func contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetEmails() {
	data := make(map[uint]string, 0)
	users := make([]user, 0)

	config.ProdDB.Model(&user{}).Unscoped().Order("id asc").Find(&users)

	for _, ur := range users {

		data[ur.ID] = ur.Email
	}

	usersJson, _ := json.Marshal(data)
	err := ioutil.WriteFile("./emails.json", usersJson, 0644)

	log.Println(data)
	if err != nil {
		log.Println("file error", err)
	}
	GetIDs()
}

func GetIDs() {
	data := make([]uint, 0)
	users := make(map[uint]string, 0)

	file, _ := ioutil.ReadFile("./emails.json")
	json.Unmarshal(file, &users)

	for id, _ := range users {

		data = append(data, id)
	}

	usersJson, _ := json.Marshal(data)
	err := ioutil.WriteFile("./uids.json", usersJson, 0644)

	log.Println(data)
	if err != nil {
		log.Println("file error", err)
	}
}
