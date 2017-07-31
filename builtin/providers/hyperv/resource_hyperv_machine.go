package hyperv

import (
	"fmt"
	"log"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHyperVMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceHyperVMachineCreate,
		Read:   resourceHyperVMachineRead,
		Update: resourceHyperVMachineUpdate,
		Delete: resourceHyperVMachineDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"generation": {
				Type:     schema.TypeInt,
				Optional: true,
				Default: 1,
				ForceNew: true,
			},

			"allow_unverified_paths": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"automatic_critical_error_action": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_critical_error_action_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_start_action": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_start_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_stop_action": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"checkpoint_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"dynamic_memory": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"guest_controlled_cache_types": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"high_memory_mapped_io_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"lock_on_disconnect": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"low_memory_mapped_io_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"memory_maximum_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"memory_minimum_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"memory_startup_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"processor_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"smart_paging_file_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"snapshot_file_location": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"static_memory": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceHyperVMachineCreate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] creating hyperv machine: %#v", d)
	c := meta.(* Client)

	name := ""

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	generation := (d.Get("generation")).(int)
	allowUnverifiedPaths := (d.Get("allow_unverified_paths")).(bool)
	automaticCriticalErrorAction := (d.Get("automatic_critical_error_action")).(CriticalErrorAction)
	automaticCriticalErrorActionTimeout := (d.Get("automatic_critical_error_action_timeout")).(int32)
	automaticStartAction := (d.Get("automatic_start_action")).(StartAction)
	automaticStartDelay := (d.Get("automatic_start_delay")).(int32)
	automaticStopAction := (d.Get("automatic_stop_action")).(StopAction)
	checkpointType := (d.Get("checkpoint_type")).(CheckpointType)
	dynamicMemory := (d.Get("dynamic_memory")).(bool)
	guestControlledCacheTypes := (d.Get("guest_controlled_cache_types")).(bool)
	highMemoryMappedIoSpace := (d.Get("high_memory_mapped_io_space")).(int64)
	lockOnDisconnect := (d.Get("lock_on_disconnect")).(OnOffState)
	lowMemoryMappedIoSpace := (d.Get("low_memory_mapped_io_space")).(int32)
	memoryMaximumBytes := (d.Get("memory_maximum_bytes")).(int64)
	memoryMinimumBytes := (d.Get("memory_minimum_bytes")).(int64)
	memoryStartupBytes := (d.Get("memory_startup_bytes")).(int64)
	notes := (d.Get("notes")).(string)
	processorCount := (d.Get("processor_count")).(int64)
	smartPagingFilePath := (d.Get("smart_paging_file_path")).(string)
	snapshotFileLocation := (d.Get("snapshot_file_location")).(string)
	staticMemory := (d.Get("static_memory")).(bool)

	err = c.CreateVM(name, generation, allowUnverifiedPaths, automaticCriticalErrorAction, automaticCriticalErrorActionTimeout, automaticStartAction, automaticStartDelay, automaticStopAction, checkpointType, dynamicMemory, guestControlledCacheTypes, highMemoryMappedIoSpace, lockOnDisconnect, lowMemoryMappedIoSpace, memoryMaximumBytes, memoryMinimumBytes, memoryStartupBytes, notes, processorCount, smartPagingFilePath, snapshotFileLocation, staticMemory)
	return err
}

func resourceHyperVMachineRead(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] reading hyperv machine: %#v", d)
	c := meta.(* Client)

	name := ""

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	s, err := c.GetVM(name)

	if err != nil {
		return err
	}

	d.Set("generation", s.Generation)
	d.Set("allow_unverified_paths", s.AllowUnverifiedPaths)
	d.Set("automatic_critical_error_action", s.AutomaticCriticalErrorAction)
	d.Set("automatic_critical_error_action_timeout", s.AutomaticCriticalErrorActionTimeout)
	d.Set("automatic_start_action", s.AutomaticStartAction)
	d.Set("automatic_start_delay", s.AutomaticStartDelay)
	d.Set("automatic_stop_action", s.AutomaticStopAction)
	d.Set("checkpoint_type", s.CheckpointType)
	d.Set("dynamic_memory", s.DynamicMemory)
	d.Set("guest_controlled_cache_types", s.GuestControlledCacheTypes)
	d.Set("high_memory_mapped_io_space", s.HighMemoryMappedIoSpace)
	d.Set("lock_on_disconnect", s.LockOnDisconnect)
	d.Set("low_memory_mapped_io_space", s.LowMemoryMappedIoSpace)
	d.Set("memory_maximum_bytes", s.MemoryMaximumBytes)
	d.Set("memory_minimum_bytes", s.MemoryMinimumBytes)
	d.Set("memory_startup_bytes", s.MemoryStartupBytes)
	d.Set("notes", s.Notes)
	d.Set("processor_count", s.ProcessorCount)
	d.Set("smart_paging_file_path", s.SmartPagingFilePath)
	d.Set("snapshot_file_location", s.SnapshotFileLocation)
	d.Set("static_memory", s.StaticMemory)

	return nil
}

func resourceHyperVMachineUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] updating hyperv machine: %#v", d)
	c := meta.(* Client)

	name := ""

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	//generation := (d.Get("generation")).(int)
	allowUnverifiedPaths := (d.Get("allow_unverified_paths")).(bool)
	automaticCriticalErrorAction := (d.Get("automatic_critical_error_action")).(CriticalErrorAction)
	automaticCriticalErrorActionTimeout := (d.Get("automatic_critical_error_action_timeout")).(int32)
	automaticStartAction := (d.Get("automatic_start_action")).(StartAction)
	automaticStartDelay := (d.Get("automatic_start_delay")).(int32)
	automaticStopAction := (d.Get("automatic_stop_action")).(StopAction)
	checkpointType := (d.Get("checkpoint_type")).(CheckpointType)
	dynamicMemory := (d.Get("dynamic_memory")).(bool)
	guestControlledCacheTypes := (d.Get("guest_controlled_cache_types")).(bool)
	highMemoryMappedIoSpace := (d.Get("high_memory_mapped_io_space")).(int64)
	lockOnDisconnect := (d.Get("lock_on_disconnect")).(OnOffState)
	lowMemoryMappedIoSpace := (d.Get("low_memory_mapped_io_space")).(int32)
	memoryMaximumBytes := (d.Get("memory_maximum_bytes")).(int64)
	memoryMinimumBytes := (d.Get("memory_minimum_bytes")).(int64)
	memoryStartupBytes := (d.Get("memory_startup_bytes")).(int64)
	notes := (d.Get("notes")).(string)
	processorCount := (d.Get("processor_count")).(int64)
	smartPagingFilePath := (d.Get("smart_paging_file_path")).(string)
	snapshotFileLocation := (d.Get("snapshot_file_location")).(string)
	staticMemory := (d.Get("static_memory")).(bool)

	err = c.UpdateVM(name, allowUnverifiedPaths, automaticCriticalErrorAction, automaticCriticalErrorActionTimeout, automaticStartAction, automaticStartDelay, automaticStopAction, checkpointType, dynamicMemory, guestControlledCacheTypes, highMemoryMappedIoSpace, lockOnDisconnect, lowMemoryMappedIoSpace, memoryMaximumBytes, memoryMinimumBytes, memoryStartupBytes, notes, processorCount, smartPagingFilePath, snapshotFileLocation, staticMemory)

	return err
}

func resourceHyperVMachineDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] deleting hyperv machine: %#v", d)

	c := meta.(* Client)

	name := ""

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	err = c.DeleteVM(name)

	return err
}