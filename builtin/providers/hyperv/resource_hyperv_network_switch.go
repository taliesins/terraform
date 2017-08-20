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

			"net_adapter_names": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"default_flow_minimum_bandwidth_absolute": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"default_flow_minimum_bandwidth_weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"default_queue_vmmq_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"default_queue_vmmq_queue_pairs": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"default_queue_vrss_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	netAdapterNames := (d.Get("net_adapter_names")).([]string)
	defaultFlowMinimumBandwidthAbsolute := (d.Get("default_flow_minimum_bandwidth_absolute")).(int)
	defaultFlowMinimumBandwidthWeight := (d.Get("default_flow_minimum_bandwidth_weight")).(int)
	defaultQueueVmmqEnabled := (d.Get("default_queue_vmmq_enabled")).(bool)
	defaultQueueVmmqQueuePairs := (d.Get("default_queue_vmmq_queue_pairs")).(int)
	defaultQueueVrssEnabled := (d.Get("default_queue_vrss_enabled")).(bool)

	err = c.CreateVMSwitch(switchName, notes, allowManagementOS, embeddedTeamingEnabled, iovEnabled, packetDirectEnabled, bandwidthReservationMode, switchType, netAdapterInterfaceDescriptions, netAdapterNames, defaultFlowMinimumBandwidthAbsolute, defaultFlowMinimumBandwidthWeight, defaultQueueVmmqEnabled, defaultQueueVmmqQueuePairs, defaultQueueVrssEnabled)
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
	d.Set("net_adapter_names", s.NetAdapterNames)
	d.Set("default_flow_minimum_bandwidth_absolute", s.DefaultFlowMinimumBandwidthAbsolute)
	d.Set("default_flow_minimum_bandwidth_weight", s.DefaultFlowMinimumBandwidthWeight)
	d.Set("default_queue_vmmq_enabled", s.DefaultQueueVmmqEnabled)
	d.Set("default_queue_vmmq_queue_pairs", s.DefaultQueueVmmqQueuePairs)
	d.Set("default_queue_vrss_enabled", s.DefaultQueueVrssEnabled)

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
	netAdapterNames := (d.Get("net_adapter_names")).([]string)
	defaultFlowMinimumBandwidthAbsolute := (d.Get("default_flow_minimum_bandwidth_absolute")).(int)
	defaultFlowMinimumBandwidthWeight := (d.Get("default_flow_minimum_bandwidth_weight")).(int)
	defaultQueueVmmqEnabled := (d.Get("default_queue_vmmq_enabled")).(bool)
	defaultQueueVmmqQueuePairs := (d.Get("default_queue_vmmq_queue_pairs")).(int)
	defaultQueueVrssEnabled := (d.Get("default_queue_vrss_enabled")).(bool)

	err = c.UpdateVMSwitch(switchName, notes, allowManagementOS, switchType, netAdapterInterfaceDescriptions, netAdapterNames, defaultFlowMinimumBandwidthAbsolute, defaultFlowMinimumBandwidthWeight, defaultQueueVmmqEnabled, defaultQueueVmmqQueuePairs, defaultQueueVrssEnabled)
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