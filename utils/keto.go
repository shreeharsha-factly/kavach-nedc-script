package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/factly/x/loggerx"
	"github.com/shreeharsha-factly/kavach-nedc-script/config"
)

type KetoSubjectSet struct {
	Namespace string `json:"namespace"`
	Object    string `json:"object"`
	Relation  string `json:"relation"`
}

type KetoRelationTupleWithSubjectID struct {
	KetoSubjectSet
	SubjectID string `json:"subject_id"`
}

func MigrateOrgDatatoKeto() error {
	orgUsers := make([]OrganisationUser, 0)
	err := config.LocalDB.Model(&OrganisationUser{}).Order("id ASC").Find(&orgUsers).Error
	if err != nil {
		return err
	}

	for _, eachOrgUser := range orgUsers {
		tuple := KetoRelationTupleWithSubjectID{
			KetoSubjectSet: KetoSubjectSet{
				Namespace: "organisations",
				Object:    fmt.Sprintf("org:%d", eachOrgUser.OrganisationID),
				Relation:  eachOrgUser.Role,
			},
			SubjectID: fmt.Sprintf("%d", eachOrgUser.UserID),
		}

		err = CreateRelationTupleWithSubjectID(&tuple)
		if err != nil {
			log.Println("error at", eachOrgUser.ID, eachOrgUser.UserID)
			return err
		}
		log.Println(eachOrgUser.ID, eachOrgUser.UserID)
	}
	return nil
}

func CreateRelationTupleWithSubjectID(tuple *KetoRelationTupleWithSubjectID) error {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(&tuple)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", "http://127.0.0.1:7710/admin/relation-tuples", buf)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		responseBody := make(map[string]interface{})
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			return err
		}
		errorMessage, ok := responseBody["message"].(string)
		if ok {
			loggerx.Error(errors.New(errorMessage))
		}
		return errors.New("error in creating the relation tuple")
	}
	return nil
}
