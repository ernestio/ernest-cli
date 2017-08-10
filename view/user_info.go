package view

import (
	"fmt"

	"github.com/ernestio/ernest-cli/model"
)

// PrintUserInfo : ...
func PrintUserInfo(u model.User) {
	fmt.Println("Username: ", u.Username)
	fmt.Println("Projects:")
	for _, v := range u.Projects {
		fmt.Println("  " + v)
	}
	fmt.Println("Environments:")
	for _, v := range u.Envs {
		fmt.Println("  " + v)
	}

}
