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
	"encoding/base64"
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
		Use:     "org",
		Aliases: []string{"organization"},
		Short:   "Create an organization",
		Long: `Create an organization from a file.  JSON format is accepted. 

For example:

  # Create an organization using the data in org.json
  apimanager create org -f ./org.json

  apimanager create organization -f ./org.json
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if file == "" {
				cmd.MarkFlagRequired("name")
			}
		},
		Run: createOrganization,
	}

	orgDelCmd = &cobra.Command{
		Use:     "org",
		Aliases: []string{"organization"},
		Short:   "Delete an organization",
		Long: `Delete an organization by name. 
	
	For example:
	
	  # Create an organization using the data in org.json
	  apimanager delete org -n orgName`,
		Run: deleteOrganization,
	}

	orgListCmd = &cobra.Command{
		Use:     "orgs",
		Aliases: []string{"organizations", "org"},
		Short:   "List all organizations",
		Long: `lists all organization by name and ID. 
	
	For example:
	
	  # lists all organization using the data in org.json
	  apimanager list orgs `,
		Run: listOrganizations,
	}
)

func init() {
	createCmd.AddCommand(orgCmd)
	deleteCmd.AddCommand(orgDelCmd)
	listCmd.AddCommand(orgListCmd)

	orgCmd.Flags().StringVarP(&file, "file", "f", "", "filename used to create the organization resource")
	orgCmd.Flags().StringVarP(&orgName, "name", "n", "", "organization name")
	orgCmd.Flags().BoolVarP(&enabled, "enabled", "e", false, "enable the organization")
	orgCmd.Flags().BoolVarP(&development, "development", "d", false, "enable organization for development")
	orgCmd.Flags().StringVarP(&image, "image", "i", "", "filename of the image to be used")

	orgDelCmd.Flags().StringVarP(&orgName, "name", "n", "", "The name to store Organization name")
	orgDelCmd.MarkFlagRequired("name")
}

func createOrganization(cmd *cobra.Command, args []string) {
	cfg := getConfig()
	org := apimgr.Organization{}

	if file != "" {
		orgBody, err := ioutil.ReadFile(file) // pass the file name with path
		if err != nil {
			fmt.Print(err)
		}
		err = json.Unmarshal([]byte(orgBody), &org)
		if err != nil {
			utils.PrettyPrintErr("Error unmarshaling org json: %v", err)
		}
		if org.Name == "" && orgName == "" {
			fmt.Printf("Organization name is not provided\n")
			return
		}
		org.Name = orgName
	} else {
		if image != "" {
			bImage, err := ioutil.ReadFile(image) // pass the file name with path
			if err != nil {
				fmt.Print(err)
			}
			org.Image = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(bImage)
		}

		org.Name = orgName
		org.Description = orgName + "Organization"
		org.Phone = "+1 877-564-7700"
		org.Email = orgName + "@apimanager.com"
		org.Enabled = enabled
		org.Development = development
		org.VirtualHost = ""
	}

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	orgVars := &apimgr.OrganizationsPostOpts{}
	orgVars.Body = optional.NewInterface(org)

	org, _, err := client.OrganizationsApi.OrganizationsPost(context.Background(), orgVars)
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

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	_, err := client.OrganizationsApi.OrganizationsIdDelete(context.Background(), orgID)
	if err != nil {
		utils.PrettyPrintErr("Unable to delete the Organization: %v", err)
		return
	}
	utils.PrettyPrintInfo("Organization %v deleted", orgName)
	return
}

func getOrganizationByName(args []string) string {
	cfg := getConfig()

	client := &apimgr.APIClient{}
	client = apimgr.NewAPIClient(cfg)

	getOrgVars := &apimgr.OrganizationsGetOpts{}

	getOrgVars.Field = optional.NewInterface("name")
	getOrgVars.Op = optional.NewInterface("eq")
	getOrgVars.Value = optional.NewInterface(orgName)

	orgs, _, err := client.OrganizationsApi.OrganizationsGet(context.Background(), getOrgVars)
	if err != nil {
		utils.PrettyPrintErr("Error finding the organizations: %v", err)
		os.Exit(0)
	}
	if len(orgs) != 0 {
		return orgs[0].Id
	}
	utils.PrettyPrintInfo("Organization %v not found ", orgName)
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
		fmt.Fprintf(stdout, "ID\tNAME\tDESCRIPTION\tCONTACT\n")
		for _, org := range orgs {
			fmt.Fprintf(stdout, "%v\t%v\t%v\t%v\n", org.Id, org.Name, org.Description, org.Email)
		}
		fmt.Fprint(stdout)
		stdout.Flush()
	} else {
		utils.PrettyPrintInfo("No Organizations found ")
		return
	}
}
