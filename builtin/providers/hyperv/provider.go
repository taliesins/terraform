package hyperv

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"io/ioutil"
)

const (
	DefaultHost = "127.0.0.1"

	DefaultUseHTTPS = false

	DefaultAllowInsecure = false

	// DefaultUser is used if there is no user given
	DefaultUser = "Administrator"

	// DefaultPort is used if there is no port given
	DefaultPort = 5985

	DefaultCACertFile = ""

	// DefaultScriptPath is used as the path to copy the file to
	// for remote execution if not provided otherwise.
	DefaultScriptPath = "C:/Temp/terraform_%RAND%.cmd"

	// DefaultTimeout is used if there is no timeout given
	DefaultTimeoutString = "30s"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_USER", DefaultUser),
				Description: "The user name for HyperV API operations.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_PASSWORD", nil),
				Description: "The user password for HyperV API operations.",
			},

			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_HOST", DefaultHost),
				Description: "The HyperV server host for HyperV API operations.",
			},

			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_PORT", DefaultPort),
				Description: "The HyperV server port for HyperV API operations.",
			},

			"https": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_HTTPS", DefaultUseHTTPS),
				Description: "Should https communication be used for HyperV API operations.",
			},

			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_INSECURE", DefaultAllowInsecure),
				Description: "Should insecure communication be used for HyperV API operations.",
			},

			"cacert_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_CACERT_PATH", DefaultCACertFile),
				Description: "The ca cert to use for HyperV API operations.",
			},

			"script_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_SCRIPT_PATH", DefaultScriptPath),
				Description: "The script path on host for HyperV API operations.",
			},

			"timeout": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYPERV_TIMEOUT", DefaultTimeoutString),
				Description: "Timeout for HyperV API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"hyperv_network_switch": resourceHyperVNetworkSwitch(),
			"hyperv_machine": resourceHyperVMachine(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var cacert *[]byte = nil
	cacertPath := d.Get("cacert_path").(string)
	if cacertPath != "" {
		if _, err := os.Stat(cacertPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("cacertPath does not exist - %s.", cacertPath)
		}

		cacertBytes, err := ioutil.ReadFile(cacertPath)
		if err != nil {
			return nil, err
		}
		cacert = &cacertBytes
	}

	config := Config{
		User:       d.Get("user").(string),
		Password:   d.Get("password").(string),
		Host: 		d.Get("host").(string),
		Port: 		d.Get("port").(int),
		HTTPS:		d.Get("https").(bool),
		CACert:		cacert,
		Insecure:	d.Get("insecure").(bool),
		ScriptPath:	d.Get("script_path").(string),
		Timeout:	d.Get("timeout").(string),
	}

	return config.Client()
}