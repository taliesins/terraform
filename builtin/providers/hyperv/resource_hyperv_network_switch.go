package hyperv

import (
	"fmt"
	"log"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHyperVNetworkSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceHyperVNetworkSwitchCreate,
		Read:   resourceHyperVNetworkSwitchRead,
		Update: resourceHyperVNetworkSwitchUpdate,
		Delete: resourceHyperVNetworkSwitchDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"allow_management_os": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enable_embedded_teaming": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"enable_iov": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"enable_packet_direct": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"minimum_bandwidth_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},

			"switch_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"net_adapter_interface_descriptions": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceHyperVNetworkSwitchCreate(d *schema.ResourceData, meta interface{}) (err error) {

	log.Printf("[DEBUG] creating hyperv switch: %#v", d)
	c := meta.(* client)

	switchName := ""

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	notes := (d.Get("notes")).(string)
	allowManagementOS := (d.Get("allow_management_os")).(bool)
	embeddedTeamingEnabled := (d.Get("enable_embedded_teaming")).(bool)
	iovEnabled := (d.Get("enable_iov")).(bool)
	packetDirectEnabled := (d.Get("enable_packet_direct")).(bool)
	bandwidthReservationMode := (d.Get("minimum_bandwidth_mode")).(int)
	switchType := (d.Get("switch_type")).(int)
	netAdapterInterfaceDescriptions := (d.Get("net_adapter_interface_descriptions")).([]string)

	err = c.CreateVMSwitch(switchName, notes, allowManagementOS, embeddedTeamingEnabled, iovEnabled, packetDirectEnabled, bandwidthReservationMode, switchType, netAdapterInterfaceDescriptions)
	return err
}

func resourceHyperVNetworkSwitchRead(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] reading hyperv switch: %#v", d)
	c := meta.(* client)

	switchName := ""

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	s, err := c.GetVMSwitch(switchName)

	if err != nil {
		return err
	}

	d.Set("notes", s.Notes)
	d.Set("allow_management_os", s.AllowManagementOS)
	d.Set("enable_embedded_teaming", s.EmbeddedTeamingEnabled)
	d.Set("enable_iov", s.IovEnabled)
	d.Set("enable_packet_direct", s.PacketDirectEnabled)
	d.Set("minimum_bandwidth_mode", s.BandwidthReservationMode)
	d.Set("switch_type", s.SwitchType)
	d.Set("net_adapter_interface_descriptions", s.NetAdapterInterfaceDescriptions)

	return nil
}

func resourceHyperVNetworkSwitchUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] updating hyperv switch: %#v", d)
	c := meta.(* client)

	switchName := ""

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	notes := (d.Get("notes")).(string)
	allowManagementOS := (d.Get("allow_management_os")).(bool)
	//embeddedTeamingEnabled := (d.Get("enable_embedded_teaming")).(bool)
	//iovEnabled := (d.Get("enable_iov")).(bool)
	//packetDirectEnabled := (d.Get("enable_packet_direct")).(bool)
	//bandwidthReservationMode := (d.Get("minimum_bandwidth_mode")).(int)
	switchType := (d.Get("switch_type")).(int)
	netAdapterInterfaceDescriptions := (d.Get("net_adapter_interface_descriptions")).([]string)

	err = c.UpdateVMSwitch(switchName, notes, allowManagementOS, switchType, netAdapterInterfaceDescriptions)
	return err
}

func resourceHyperVNetworkSwitchDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] deleting hyperv switch: %#v", d)

	c := meta.(* client)

	switchName := ""

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	err = c.DeleteVMSwitch(switchName)

	return err
}