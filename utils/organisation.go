package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/shreeharsha-factly/kavach-nedc-script/config"
)

type Organisation struct {
	Base
	Title             string             `gorm:"column:title" json:"title"`
	Slug              string             `gorm:"column:slug" json:"slug"`
	Description       string             `gorm:"column:description" json:"description"`
	FeaturedMediumID  *uint              `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	Medium            *Medium            `gorm:"foreignKey:featured_medium_id" json:"medium"`
	OrganisationUsers []OrganisationUser `gorm:"foreignKey:organisation_id" json:"organisation_users"`
}

// OrganisationUser model definition
type OrganisationUser struct {
	Base
	UserID         uint          `gorm:"column:user_id" json:"user_id"`
	User           *User         `gorm:"foreignKey:user_id" json:"user"`
	OrganisationID uint          `gorm:"column:organisation_id" json:"organisation_id"`
	Organisation   *Organisation `gorm:"foreignKey:organisation_id" json:"organisation"`
	Role           string        `gorm:"column:role" json:"role"`
}

type OrganisationToken struct {
	Base
	Name           string        `gorm:"column:name" json:"name"`
	Description    string        `gorm:"column:description" json:"description"`
	OrganisationID uint          `gorm:"column:organisation_id" json:"organisation_id"`
	Organisation   *Organisation `gorm:"foreignKey:organisation_id" json:"organisation"`
	Token          string        `gorm:"column:token" json:"token"`
}

type organisation struct {
	Base
	Title            string  `gorm:"column:title" json:"title"`
	Slug             string  `gorm:"column:slug" json:"slug"`
	Description      string  `gorm:"column:description" json:"description"`
	FeaturedMediumID *uint   `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	Medium           *medium `gorm:"foreignKey:featured_medium_id" json:"medium"`
}

type organisationUser struct {
	Base
	UserID         uint          `gorm:"column:user_id" json:"user_id"`
	User           *user         `json:"user"`
	OrganisationID uint          `gorm:"column:organisation_id" json:"organisation_id"`
	Organisation   *organisation `json:"organisation"`
	Role           string        `gorm:"column:role" json:"role"`
}

func CreateOrganisation() {

	users := make(map[uint]uint, 0)
	media := make(map[uint]uint, 0)

	file, _ := ioutil.ReadFile("./users.json")
	json.Unmarshal(file, &users)

	file, _ = ioutil.ReadFile("./media.json")
	json.Unmarshal(file, &media)

	oids := make([]uint, 0)

	file, _ = ioutil.ReadFile("./oids.json")
	json.Unmarshal(file, &oids)

	data := make(map[uint]uint, 0)
	orgs := make([]organisation, 0)

	config.ProdDB.Model(&organisation{}).Unscoped().Order("id asc").Find(&orgs)

	for _, org := range orgs {

		if contains(oids, org.ID) {
			o := Organisation{
				Base: Base{
					CreatedAt:   org.CreatedAt,
					UpdatedAt:   org.UpdatedAt,
					DeletedAt:   org.DeletedAt,
					CreatedByID: users[org.Base.CreatedByID],
					UpdatedByID: users[org.Base.UpdatedByID],
				},
				Title:       org.Title,
				Slug:        org.Slug,
				Description: org.Description,
			}

			if org.FeaturedMediumID != nil {
				mediaID := media[*org.FeaturedMediumID]
				o.FeaturedMediumID = &mediaID
			}
			err := config.LocalDB.Model(&Organisation{}).Create(&o).Error

			if err != nil {
				log.Println("ID", o.ID)
				log.Fatal(err)

				orgsJson, _ := json.Marshal(data)
				err = ioutil.WriteFile("./organisations.json", orgsJson, 0644)

				log.Println(data)
				if err != nil {
					log.Println("file error", err)
				}
			}
			data[org.ID] = o.ID
		}
	}

	orgsJson, _ := json.Marshal(data)
	err := ioutil.WriteFile("./organisations.json", orgsJson, 0644)

	log.Println(data)
	if err != nil {
		log.Println("file error", err)
	}
}

func CreateOrganisationUser() {

	users := make(map[uint]uint, 0)
	organisations := make(map[uint]uint, 0)

	file, _ := ioutil.ReadFile("./users.json")
	json.Unmarshal(file, &users)

	file, _ = ioutil.ReadFile("./organisations.json")
	json.Unmarshal(file, &organisations)

	orgUsers := make([]organisationUser, 0)

	config.ProdDB.Model(&organisationUser{}).Order("id asc").Find(&orgUsers)

	for _, org := range orgUsers {

		if users[org.UserID] > 0 && organisations[org.OrganisationID] > 0 {
			o := OrganisationUser{
				Base: Base{
					CreatedAt:   org.CreatedAt,
					UpdatedAt:   org.UpdatedAt,
					DeletedAt:   org.DeletedAt,
					CreatedByID: users[org.Base.CreatedByID],
					UpdatedByID: users[org.Base.UpdatedByID],
				},
				UserID:         users[org.UserID],
				OrganisationID: organisations[org.OrganisationID],
				Role:           org.Role,
			}

			err := config.LocalDB.Model(&organisationUser{}).Create(&o).Error

			if err != nil {
				log.Println("ID", o.ID)
				log.Fatal(err)

			}
			log.Println("done", org.ID)
		}
	}

}

func GetOrgs() {
	data := make(map[uint]string, 0)
	orgs := make([]organisation, 0)

	config.ProdDB.Model(&organisation{}).Unscoped().Order("id asc").Find(&orgs)

	for _, org := range orgs {
		data[org.ID] = org.Title
	}

	usersJson, _ := json.Marshal(data)
	err := ioutil.WriteFile("./orgs.json", usersJson, 0644)

	log.Println(data)
	if err != nil {
		log.Println("file error", err)
	}
}
