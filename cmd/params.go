package cmd

import (
	"crypto/tls"
	"net/http"
	"os"
	"text/tabwriter"

	"git.ecd.axway.int/apigov/kubecrt-vms/utils"
	"github.com/skckadiyala/apimanager/apimgr"
	"github.com/spf13/viper"
)

var (
	file         string // blob file location
	config       string // config file location
	orgName      string
	appName      string
	keyID        string
	name         string
	password     string
	apiName      string
	security     string
	resourcePath string
	certPath     string
	proxyState   string
	proxyVersion string
)

type configAPI struct {
	APIManagerHost string `yaml:"apiManagerHost"`
	APIManagerPort string `yaml:"apiManagerPort"`
	Authorization  string `yaml:"authorization"`
}

func getConfig() *apimgr.Configuration {
	if viper.GetString("apimanagerhost") == "" || viper.GetString("apimanagerport") == "" {
		utils.PrettyPrintErr("Please login to API Manager, use 'login' command")
		os.Exit(0)
	}

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	cfg := apimgr.NewConfiguration()
	cfg.Host = viper.GetString("apimanagerhost") + ":" + viper.GetString("apimanagerport")
	cfg.Scheme = "https"
	cfg.AddDefaultHeader("Authorization", "Basic "+viper.GetString("authorization"))
	cfg.HTTPClient = &http.Client{Transport: transCfg}
	return cfg
}

func fmtDisplay() *tabwriter.Writer {
	writeTab := new(tabwriter.Writer)
	writeTab.Init(os.Stdout, 0, 8, 0, '\t', 0)
	return writeTab
}
