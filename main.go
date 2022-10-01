package main

import (
	"github.com/shreeharsha-factly/kavach-nedc-script/config"
	"github.com/shreeharsha-factly/kavach-nedc-script/utils"
)

func main() {
	config.SetupKavachProdDB()
	config.LocalSetupDB()

	utils.GetEmails()
	// utils.CreateUser()
	// utils.CreateMedium()

	// utils.GetOrgs()
	// utils.CreateOrganisation()
	// utils.CreateOrganisationUser()
	//utils.MigrateOrgDatatoKeto()
}
