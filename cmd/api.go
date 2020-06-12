/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

// apiCmd represents the api command
var (
	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Create an api",
		Long: `Create an api from a swagger file. 

For example:

# Create an api using the data in swagger.json
apimanager create api -n <name> -f swagger.json -o <orgName> `,
		Run: createBackendAPI,
	}

	apiListCmd = &cobra.Command{
		Use:   "apis",
		Short: "List all backend apis",
		Long: `List all backend apis from a swagger file. 

For example:

# List all apis using the data 
apimanager list apis `,
		Run: listBackendAPI,
	}

	apiDelCmd = &cobra.Command{
		Use:   "api",
		Short: "Delete an API",
		Long: `Delete an API from a swagger file. 

For example:ß

# Delete an api using the data 
apimanager delete api -n <name> `,
		Run: deleteAPI,
	}
)

func init() {
	createCmd.AddCommand(apiCmd)
	listCmd.AddCommand(apiListCmd)
	deleteCmd.AddCommand(apiDelCmd)

	apiCmd.Flags().StringVarP(&file, "file", "f", "", "The filename of the swagger api to be stored")
	apiCmd.MarkFlagRequired("file")
	apiCmd.Flags().StringVarP(&name, "name", "n", "", "The name to store API name")
	apiCmd.MarkFlagRequired("name")
	apiCmd.Flags().StringVarP(&orgName, "orgName", "o", "", "The name to store Organization name")
	apiCmd.MarkFlagRequired("orgName")

	apiDelCmd.Flags().StringVarP(&name, "name", "n", "", "The name to store API name")
	apiDelCmd.MarkFlagRequired("name")
}

func createBackendAPI(cmd *cobra.Command, args []string) {
	utils.PrettyPrintInfo("Creating backend API")

	apiName := name
	name = orgName
	cfg := getConfig()
	orgID := getOrganizationByName(args)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	file, err := os.Open(file)
	if err != nil {
		utils.PrettyPrintErr("Error Opening file: %v", err)
	}
	beAPI, _, err := client.APIRepositoryApi.ApirepoImportPost(context.Background(), orgID, apiName, "swagger", file)
	if err != nil {
		utils.PrettyPrintErr("Error Creating Backend API: %v", err)
		return
	}
	utils.PrettyPrintInfo("Backend API Name %v with ID: %v created", beAPI.Name, beAPI.Id)
	return
}

func listBackendAPI(cmd *cobra.Command, args []string) {
	cfg := getConfig()
	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)
	apiGetOpts := &apimgr.ApirepoGetOpts{}
	stdout := fmtDisplay()

	apis, _, err := client.APIRepositoryApi.ApirepoGet(context.Background(), apiGetOpts)
	if err != nil {
		utils.PrettyPrintErr("Error Creating Backend API: %v", err)
		return
	}

	if len(apis) != 0 {
		fmt.Fprintf(stdout, "NAME\tID\n")
		for _, api := range apis {
			fmt.Fprintf(stdout, "%v\t%v\n", api.Name, api.Id)
		}
		fmt.Fprint(stdout)
		stdout.Flush()
	} else {
		utils.PrettyPrintInfo("No backend api's found ")
		return
	}
}

func getAPIByName(args []string) string {
	cfg := getConfig()

	utils.PrettyPrintInfo("Finding API %v ....", name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	apiGetOpts := &apimgr.ApirepoGetOpts{}

	apiGetOpts.Field = optional.NewInterface("name")
	apiGetOpts.Op = optional.NewInterface("eq")
	apiGetOpts.Value = optional.NewInterface(name)

	apis, _, err := client.APIRepositoryApi.ApirepoGet(context.Background(), apiGetOpts)
	if err != nil {
		utils.PrettyPrintErr("Error finding backend API: %v", err)
		os.Exit(0)
	}
	if len(apis) != 0 {
		utils.PrettyPrintInfo("Backend API found: %v", apis[0].Name)
		return apis[0].Id
	}
	utils.PrettyPrintInfo("Backend API %v not found ", name)
	os.Exit(0)
	return apis[0].Id
}

func deleteAPI(cmd *cobra.Command, args []string) {
	cfg := getConfig()

	apiID := getAPIByName(args)

	utils.PrettyPrintInfo("Deleting API %v ....", name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	_, err := client.APIRepositoryApi.ApirepoIdDelete(context.Background(), apiID)
	if err != nil {
		utils.PrettyPrintErr("Unable to delete the backend API: %v", err)
		return
	}
	return
}
