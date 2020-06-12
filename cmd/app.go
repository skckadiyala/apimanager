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
	"fmt"
	"os"

	"github.com/antihax/optional"
	"github.com/skckadiyala/apimanager/apimgr"
	"github.com/skckadiyala/kubecrt-vms/utils"
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var (
	appCmd = &cobra.Command{
		Use:   "app",
		Short: "Create an application",
		Long: `Create an application by name. 

	For example:
	
	  # Create an application by name
	  apimanager create app -n <appname> -o <orgName>
	`,
		Run: createApplication,
	}
	appDelCmd = &cobra.Command{
		Use:   "app",
		Short: "Delete an application",
		Long: `Delete application from a file. JSON format is accepted. 
	
	For example:
	
	# delete user by name
	apimanager delete app -n <appname> 
	`,
		Run: deleteApplication,
	}

	appListCmd = &cobra.Command{
		Use:   "apps",
		Short: "List all applications",
		Long: `List all applications. 
	
	For example:
	
	# list all the applications 
	apimanager list apps 
	`,
		Run: listApplications,
	}
)

func init() {
	createCmd.AddCommand(appCmd)
	deleteCmd.AddCommand(appDelCmd)
	listCmd.AddCommand(appListCmd)

	appCmd.Flags().StringVarP(&orgName, "orgName", "o", "", "The name to store Organization name")
	appCmd.MarkFlagRequired("orgName")
	appCmd.Flags().StringVarP(&name, "name", "n", "", "The name to store application name")
	appCmd.MarkFlagRequired("name")

	appDelCmd.Flags().StringVarP(&name, "name", "n", "", "The name to store application name")
	appDelCmd.MarkFlagRequired("name")
	// appListCmd.Flags().StringVarP(&name, "name", "n", "", "The name to store application name")
	// appListCmd.MarkFlagRequired("name")

}

func createApplication(cmd *cobra.Command, args []string) {
	cfg := getConfig()

	newApp := apimgr.ApplicationRequest{}
	newApp.Name = name
	newApp.Description = name + ": is application"
	newApp.Phone = "+1 877-564-7700"
	newApp.Email = name + "@postmanpat.com"
	newApp.Apis = []string{}

	name = orgName
	orgID := getOrganizationByName(args)

	newApp.OrganizationId = orgID

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	appVars := &apimgr.ApplicationsPostOpts{}
	appVars.Body = optional.NewInterface(newApp)

	app, _, err := client.ApplicationsApi.ApplicationsPost(context.Background(), appVars)
	if err != nil {
		utils.PrettyPrintErr("Error creating application: %v", err)
		return
	}
	utils.PrettyPrintInfo("Application %v Created", app.Name)
	return
}

func deleteApplication(cmd *cobra.Command, args []string) {
	cfg := getConfig()
	appID := getApplicationByName(args)

	utils.PrettyPrintInfo("Deleting application %v ....", name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	_, err := client.ApplicationsApi.ApplicationsIdDelete(context.Background(), appID)
	if err != nil {
		utils.PrettyPrintErr("Unable to delete the application: %v", err)
		return
	}
	return
}

func getApplicationByName(args []string) string {
	cfg := getConfig()

	utils.PrettyPrintInfo("Finding application %v ....", name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	getAppVars := &apimgr.ApplicationsGetOpts{}

	getAppVars.Field = optional.NewInterface("name")
	getAppVars.Op = optional.NewInterface("eq")
	getAppVars.Value = optional.NewInterface(name)

	apps, _, err := client.ApplicationsApi.ApplicationsGet(context.Background(), getAppVars)
	if err != nil {
		utils.PrettyPrintErr("Error finding the application: %v", err)
		os.Exit(0)
	}
	if len(apps) != 0 {
		utils.PrettyPrintInfo("application found: %v", apps[0].Name)
		return apps[0].Id
	}
	utils.PrettyPrintInfo("application %v not found ", name)
	os.Exit(0)
	return apps[0].Id
}

func listApplications(cmd *cobra.Command, args []string) {
	cfg := getConfig()
	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	getAppVars := &apimgr.ApplicationsGetOpts{}

	apps, _, err := client.ApplicationsApi.ApplicationsGet(context.Background(), getAppVars)
	if err != nil {
		utils.PrettyPrintErr("Error listing the applications: %v", err)
		return
	}
	if len(apps) != 0 {
		fmt.Printf("Name \t\t ID \n")
		for _, app := range apps {
			fmt.Printf("%v \t %v \n", app.Name, app.Id)
		}
	} else {
		utils.PrettyPrintInfo("No application found ")
		return
	}
}
