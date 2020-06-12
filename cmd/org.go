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

// orgCmd represents the org command
var (
	orgCmd = &cobra.Command{
		Use:   "org",
		Short: "Create an organization",
		Long: `Create an organization from a file.  JSON format is accepted. 

For example:

  # Create an organization using the data in org.json
  apimanager create org -f ./org.json

  apimanager create organization -f ./org.json
`,
		Run: createOrganization,
	}

	orgCreateCmd = &cobra.Command{
		Hidden: true,
		Use:    "organization",
		Short:  "Create an organization",
		Long: `Create an organization from a file. JSON format is accepted. 

For example:

  # Create an organization using the data in org.json
  apimanager create org -f ./org.json

  apimanager create organization -f ./org.json
`,
		Run: createOrganization,
	}

	orgDelCmd = &cobra.Command{
		Use:   "org",
		Short: "Delete an organization",
		Long: `Delete an organization by name. 
	
	For example:
	
	  # Create an organization using the data in org.json
	  apimanager delete org -n orgName`,
		Run: deleteOrganization,
	}

	orgListCmd = &cobra.Command{
		Use:   "orgs",
		Short: "List all organizations",
		Long: `lists all organization by name and ID. 
	
	For example:
	
	  # lists all organization using the data in org.json
	  apimanager list orgs `,
		Run: listOrganizations,
	}
)

func init() {
	createCmd.AddCommand(orgCmd)
	createCmd.AddCommand(orgCreateCmd)
	deleteCmd.AddCommand(orgDelCmd)
	listCmd.AddCommand(orgListCmd)

	orgCmd.Flags().StringVarP(&file, "file", "f", "", "The filename of the raw data to be stored")
	orgCmd.MarkFlagRequired("file")
	orgCreateCmd.Flags().StringVarP(&file, "file", "f", "", "The filename of the raw data to be stored")
	orgCreateCmd.MarkFlagRequired("file")
	orgDelCmd.Flags().StringVarP(&name, "name", "n", "", "The name to store Organization name")
	orgDelCmd.MarkFlagRequired("name")
}

func createOrganization(cmd *cobra.Command, args []string) {
	cfg := getConfig()
	orgBody, err := ioutil.ReadFile(file) // pass the file name with path
	if err != nil {
		fmt.Print(err)
	}

	org := apimgr.Organization{}

	err = json.Unmarshal([]byte(orgBody), &org)
	if err != nil {
		utils.PrettyPrintErr("Error unmarshaling org json: %v", err)
	}
	utils.PrettyPrintInfo("Creating a new Organization %v....", org.Name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	orgVars := &apimgr.OrganizationsPostOpts{}
	orgVars.Body = optional.NewInterface(org)

	org, _, err = client.OrganizationsApi.OrganizationsPost(context.Background(), orgVars)
	if err != nil {
		utils.PrettyPrintErr("Error Creating Organization: %v", err)
		return
	}
	utils.PrettyPrintInfo("Organization %v Created with ID %v", org.Name, org.Id)
	return
}

func deleteOrganization(cmd *cobra.Command, args []string) {
	cfg := getConfig()

	orgID := getOrganizationByName(args)

	utils.PrettyPrintInfo("Deleting Organization %v ....", name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	_, err := client.OrganizationsApi.OrganizationsIdDelete(context.Background(), orgID)
	if err != nil {
		utils.PrettyPrintErr("Unable to delete the Organization: %v", err)
		return
	}
	return
}

func getOrganizationByName(args []string) string {
	cfg := getConfig()

	utils.PrettyPrintInfo("Finding Organization %v ....", name)

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	getOrgVars := &apimgr.OrganizationsGetOpts{}

	getOrgVars.Field = optional.NewInterface("name")
	getOrgVars.Op = optional.NewInterface("eq")
	getOrgVars.Value = optional.NewInterface(name)

	orgs, _, err := client.OrganizationsApi.OrganizationsGet(context.Background(), getOrgVars)
	if err != nil {
		utils.PrettyPrintErr("Error finding the organizations: %v", err)
		os.Exit(0)
	}
	if len(orgs) != 0 {
		utils.PrettyPrintInfo("Organization found: %v", orgs[0].Name)
		return orgs[0].Id
	}
	utils.PrettyPrintInfo("Organization %v not found ", name)
	os.Exit(0)
	return orgs[0].Id
}

func listOrganizations(cmd *cobra.Command, args []string) {
	cfg := getConfig()
	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)
	stdout := fmtDisplay()
	getOrgVars := &apimgr.OrganizationsGetOpts{}

	orgs, _, err := client.OrganizationsApi.OrganizationsGet(context.Background(), getOrgVars)
	if err != nil {
		utils.PrettyPrintErr("Error listing the organizations: %v", err)
		return
	}
	if len(orgs) != 0 {
		fmt.Fprintf(stdout, "NAME\tID\n")
		for _, org := range orgs {
			fmt.Fprintf(stdout, "%v\t%v\n", org.Name, org.Id)
		}
		fmt.Fprint(stdout)
		stdout.Flush()
	} else {
		utils.PrettyPrintInfo("No Organizations found ")
		return
	}
}
