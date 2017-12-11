package view

import (
	"fmt"

	"github.com/ernestio/ernest-go-sdk/models"
)

// PrintUserInfo : ...
func PrintUserInfo(u *models.User) {
	fmt.Println("Username: ", u.Username)
	fmt.Println("Type:     ", u.Type)
	fmt.Println("Projects:")
	for _, v := range u.Projects {
		fmt.Println("  " + v)
	}
	fmt.Println("Environments:")
	for _, v := range u.Envs {
		fmt.Println("  " + v)
	}

}
