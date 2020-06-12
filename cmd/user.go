/*
Copyright © 2020 Axway, Inc. <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/antihax/optional"
	"github.com/skckadiyala/apimanager/apimgr"
	"github.com/skckadiyala/kubecrt-vms/utils"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var (
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "Create User",
		Long: `Create user from a file. JSON format is accepted. 
	
	For example:
	
	# Create an organization using the data in org.json
	apimanager create user -f user.json
	`,
		Run: createUser,
	}

	userDelCmd = &cobra.Command{
		Use:   "user",
		Short: "Delete User",
		Long: `Delete user from a file. JSON format is accepted. 
	
	For example:
	
	# delete user by name
	apimanager delete user -n username
	`,
		Run: deleteUser,
	}

	userListCmd = &cobra.Command{
		Use:   "users",
		Short: "List all users",
		Long: `List all users. 
	
	For example:
	
	# list all the users
	apimanager list users 
	`,
		Run: listUsers,
	}
)

func init() {
	createCmd.AddCommand(userCmd)
	deleteCmd.AddCommand(userDelCmd)
	listCmd.AddCommand(userListCmd)

	userCmd.Flags().StringVarP(&file, "file", "f", "", "The filename of the raw data to be stored")
	userCmd.MarkFlagRequired("file")
	userCmd.Flags().StringVarP(&orgName, "orgName", "o", "", "The name to store Organization name")
	userCmd.MarkFlagRequired("orgName")
	userCmd.Flags().StringVarP(&password, "password", "p", "", "The password for the user")
	userCmd.MarkFlagRequired("password")

	userDelCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the username")
	userDelCmd.MarkFlagRequired("name")
}

func createUser(cmd *cobra.Command, args []string) {

	cfg := getConfig()
	userBody, err := ioutil.ReadFile(file) // pass the file name with path
	if err != nil {
		fmt.Print(err)
	}

	name = orgName
	orgID := getOrganizationByName(args)

	newUser := apimgr.User{}
	newUser.OrganizationId = orgID
	err = json.Unmarshal([]byte(userBody), &newUser)
	if err != nil {
		utils.PrettyPrintErr("Error unmarshaling user json: %v", err)
	}
	utils.PrettyPrintInfo("Creating a new user %v....", newUser.Name)
	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	userVars := &apimgr.UsersPostOpts{}
	userVars.Body = optional.NewInterface(newUser)

	user, _, err := client.UsersApi.UsersPost(context.Background(), userVars)
	if err != nil {
		utils.PrettyPrintErr("Error creating user:%v", err)
	}
	utils.PrettyPrintInfo("New user %v created", user.Name)
	changeUserPassword(user.Id, password, cfg)
	return
}

func getUserByName(args []string) string {

	cfg := getConfig()

	utils.PrettyPrintInfo("Finding User %v ....", name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	getUserVars := &apimgr.UsersGetOpts{}

	getUserVars.Field = optional.NewInterface("name")
	getUserVars.Op = optional.NewInterface("eq")
	getUserVars.Value = optional.NewInterface(name)

	users, _, err := client.UsersApi.UsersGet(context.Background(), getUserVars)
	if err != nil {
		utils.PrettyPrintErr("Error finding the user: %v", err)
		os.Exit(0)
	}
	if len(users) != 0 {
		utils.PrettyPrintInfo("User found: %v", users[0].Name)
		return users[0].Id
	}
	utils.PrettyPrintInfo("User %v not found ", name)
	return "User Not Found"

}

func listUsers(cmd *cobra.Command, args []string) {

	cfg := getConfig()
	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	getUserVars := &apimgr.UsersGetOpts{}

	users, _, err := client.UsersApi.UsersGet(context.Background(), getUserVars)
	if err != nil {
		utils.PrettyPrintErr("Error listing the users: %v", err)
		return
	}
	if len(users) != 0 {
		fmt.Printf("Name \t\t ID \n")
		for _, user := range users {
			fmt.Printf("%v \t %v \n", user.Name, user.Id)
		}
	} else {
		utils.PrettyPrintInfo("No users found ")
		return
	}
}

func deleteUser(cmd *cobra.Command, args []string) {

	cfg := getConfig()

	userID := getUserByName(args)

	utils.PrettyPrintInfo("Deleting User %v ....", name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	_, err := client.UsersApi.UsersIdDelete(context.Background(), userID)
	if err != nil {
		utils.PrettyPrintErr("Unable to delete the user: %v", err)
		return
	}
	return
}

func changeUserPassword(userID, newPassword string, cfg *apimgr.Configuration) {
	utils.PrettyPrintInfo("Change password for the UserId :%v", userID)
	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	_, err := client.UsersApi.UsersIdChangepasswordPost(context.Background(), userID, newPassword)
	if err != nil {
		utils.PrettyPrintErr("Error updating password :%v", err)
		return
	}
	utils.PrettyPrintInfo("Password updated")
	return
}
