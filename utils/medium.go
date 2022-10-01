package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shreeharsha-factly/kavach-nedc-script/config"
)

// Medium model
type Medium struct {
	Base
	Name        string         `gorm:"column:name" json:"name"`
	Slug        string         `gorm:"column:slug" json:"slug"`
	Type        string         `gorm:"column:type" json:"type"`
	Title       string         `gorm:"column:title" json:"title"`
	Description string         `gorm:"column:description" json:"description"`
	Caption     string         `gorm:"column:caption" json:"caption"`
	AltText     string         `gorm:"column:alt_text" json:"alt_text"`
	FileSize    int64          `gorm:"column:file_size" json:"file_size"`
	URL         postgres.Jsonb `gorm:"column:url" json:"url" swaggertype:"primitive,string"`
	Dimensions  string         `gorm:"column:dimensions" json:"dimensions"`
	UserID      uint           `gorm:"column:user_id" json:"user_id"`
}

type medium struct {
	Base
	Name        string         `gorm:"column:name" json:"name"`
	Slug        string         `gorm:"column:slug" json:"slug"`
	Type        string         `gorm:"column:type" json:"type"`
	Title       string         `gorm:"column:title" json:"title"`
	Description string         `gorm:"column:description" json:"description"`
	Caption     string         `gorm:"column:caption" json:"caption"`
	AltText     string         `gorm:"column:alt_text" json:"alt_text"`
	FileSize    int64          `gorm:"column:file_size" json:"file_size"`
	URL         postgres.Jsonb `gorm:"column:url" json:"url" swaggertype:"primitive,string"`
	Dimensions  string         `gorm:"column:dimensions" json:"dimensions"`
	UserID      uint           `gorm:"column:user_id" json:"user_id"`
}

func (Medium) TableName() string {
	return "media"
}
func (medium) TableName() string {
	return "media"
}

func CreateMedium() {

	users := make(map[uint]uint, 0)

	file, _ := ioutil.ReadFile("./utils/users.json")
	json.Unmarshal(file, &users)

	data := make(map[uint]uint, 0)
	media := make([]medium, 0)

	config.ProdDB.Model(&medium{}).Unscoped().Order("id asc").Find(&media)

	for _, md := range media {
		m := Medium{
			Base: Base{
				CreatedAt:   md.CreatedAt,
				UpdatedAt:   md.UpdatedAt,
				DeletedAt:   md.DeletedAt,
				CreatedByID: users[md.Base.CreatedByID],
				UpdatedByID: users[md.Base.UpdatedByID],
			},
			Name:        md.Name,
			Slug:        md.Slug,
			Type:        md.Type,
			Title:       md.Title,
			Description: md.Description,
			Caption:     md.Caption,
			AltText:     md.AltText,
			FileSize:    md.FileSize,
			URL:         md.URL,
			Dimensions:  md.Dimensions,
			UserID:      users[md.UserID],
		}
		err := config.LocalDB.Model(&Medium{}).Create(&m).Error

		if err != nil {
			log.Println("ID", m.ID)
			log.Fatal(err)

			mediaJson, _ := json.Marshal(data)
			err = ioutil.WriteFile("./utils/media.json", mediaJson, 0644)

			log.Println(data)
			if err != nil {
				log.Println("file error", err)
			}
		}
		data[md.ID] = m.ID
	}

	mediaJson, _ := json.Marshal(data)
	err := ioutil.WriteFile("./utils/media.json", mediaJson, 0644)

	log.Println(data)
	if err != nil {
		log.Println("file error", err)
	}
}
