package hyperv

import (
	"text/template"
	"encoding/json"
)

type CriticalErrorAction int

const (
	CriticalErrorAction_None CriticalErrorAction = 0
	CriticalErrorAction_Pause CriticalErrorAction = 1
)

type StartAction int

const (
	StartAction_Nothing StartAction = 2
	StartAction_StartIfRunning StartAction = 3
	StartAction_Start StartAction = 4
)

type StopAction int

const (
	StopAction_TurnOff StopAction = 2
	StopAction_Save StopAction = 3
	StopAction_ShutDown StopAction = 4
)

type CheckpointType int

const (
	CheckpointType_Disabled CheckpointType = 2
	CheckpointType_Production CheckpointType = 3
	CheckpointType_ProductionOnly CheckpointType = 4
	CheckpointType_Standard CheckpointType = 5
)

type OnOffState int

const (
	OnOffState_On OnOffState = 0
	OnOffState_Off OnOffState = 1
)

type vm struct {
	Name string
	Generation int
	AllowUnverifiedPaths bool
	AutomaticCriticalErrorAction CriticalErrorAction
	AutomaticCriticalErrorActionTimeout int32
	AutomaticStartAction StartAction
	AutomaticStartDelay int32
	AutomaticStopAction StopAction
	CheckpointType CheckpointType
	DynamicMemory bool
	GuestControlledCacheTypes bool
	HighMemoryMappedIoSpace int64
	LockOnDisconnect OnOffState
	LowMemoryMappedIoSpace int32
	MemoryMaximumBytes int64
	MemoryMinimumBytes int64
	MemoryStartupBytes int64
	Notes string
	ProcessorCount int64
	SmartPagingFilePath string
	SnapshotFileLocation string
	StaticMemory bool
}

type createVMArgs struct {
	VmJson		string
}

var createVMTemplate = template.Must(template.New("CreateVM").Parse(`
$vm = '{{.VmJson}}' | ConvertFrom-Json
$automaticCriticalErrorAction = [Microsoft.HyperV.PowerShell.CriticalErrorAction]$vm.AutomaticCriticalErrorAction
$automaticStartAction = [Microsoft.HyperV.PowerShell.StartAction]$vm.AutomaticStartAction
$automaticStopAction = [Microsoft.HyperV.PowerShell.StopAction]$vm.AutomaticStopAction
$checkpointType = [Microsoft.HyperV.PowerShell.CheckpointType]$vm.CheckpointType
$lockOnDisconnect = [Microsoft.HyperV.PowerShell.OnOffState]$vm.LockOnDisconnect

New-Vm -Name $vm.Name -Generation $vm.Generation -MemoryStartupBytes $vm.MemoryStartupBytes
Set-Vm -Name $vm.Name -GuestControlledCacheTypes:$vm.GuestControlledCacheTypes -LowMemoryMappedIoSpace $vm.LowMemoryMappedIoSpace -HighMemoryMappedIoSpace $vm.HighMemoryMappedIoSpace -ProcessorCount $vm.ProcessorCount -DynamicMemory:$vm.DynamicMemory -StaticMemory:$vm.StaticMemory -MemoryMinimumBytes $vm.MemoryMinimumBytes -MemoryMaximumBytes $vm.MemoryMaximumBytes -MemoryStartupBytes $vm.MemoryStartupBytes -AutomaticStartAction $automaticStartAction -AutomaticStopAction $automaticStopAction -AutomaticStartDelay $vm.AutomaticStartDelay -AutomaticCriticalErrorAction $automaticCriticalErrorAction -AutomaticCriticalErrorActionTimeout $vm.AutomaticCriticalErrorActionTimeout -LockOnDisconnect $lockOnDisconnect -Notes $vm.Notes -SnapshotFileLocation $vm.SnapshotFileLocation -SmartPagingFilePath $vm.SmartPagingFilePath -CheckpointType $checkpointType -AllowUnverifiedPaths $vm.AllowUnverifiedPaths
`))

func (c *HypervClient) CreateVM(
	name string,
	generation int,
	allowUnverifiedPaths bool,
	automaticCriticalErrorAction CriticalErrorAction,
	automaticCriticalErrorActionTimeout int32,
	automaticStartAction StartAction,
	automaticStartDelay int32,
	automaticStopAction StopAction,
	checkpointType CheckpointType,
	dynamicMemory bool,
	guestControlledCacheTypes bool,
	highMemoryMappedIoSpace int64,
	lockOnDisconnect OnOffState,
	lowMemoryMappedIoSpace int32,
	memoryMaximumBytes int64,
	memoryMinimumBytes int64,
	memoryStartupBytes int64,
	notes string,
	processorCount int64,
	smartPagingFilePath string,
	snapshotFileLocation string,
	staticMemory bool,
) (err error) {

	vmJson, err := json.Marshal(vm{
		Name:name,
		Generation:generation,
		AllowUnverifiedPaths: allowUnverifiedPaths,
		AutomaticCriticalErrorAction: automaticCriticalErrorAction,
		AutomaticCriticalErrorActionTimeout: automaticCriticalErrorActionTimeout,
		AutomaticStartAction: automaticStartAction,
		AutomaticStartDelay: automaticStartDelay,
		AutomaticStopAction: automaticStopAction,
		CheckpointType: checkpointType,
		DynamicMemory: dynamicMemory,
		GuestControlledCacheTypes: guestControlledCacheTypes,
		HighMemoryMappedIoSpace: highMemoryMappedIoSpace,
		LockOnDisconnect: lockOnDisconnect,
		LowMemoryMappedIoSpace: lowMemoryMappedIoSpace,
		MemoryMaximumBytes: memoryMaximumBytes,
		MemoryMinimumBytes: memoryMinimumBytes,
		MemoryStartupBytes: memoryStartupBytes,
		Notes: notes,
		ProcessorCount: processorCount,
		SmartPagingFilePath: smartPagingFilePath,
		SnapshotFileLocation: snapshotFileLocation,
		StaticMemory: staticMemory,
	})

	err = c.runFireAndForgetScript(createVMTemplate, createVMArgs{
		VmJson:string(vmJson),
	});

	return err
}

type getVMArgs struct {
	Name		string
}

var getVMTemplate = template.Must(template.New("GetVM").Parse(`
(Get-VM -name '{{.Name}}') | %{ @{
	Name=$_.Name;
	Generation=$_.Generation;
	AllowUnverifiedPaths=$_.AllowUnverifiedPaths;
	AutomaticCriticalErrorAction=$_.AutomaticCriticalErrorAction;
	AutomaticCriticalErrorActionTimeout=$_.AutomaticCriticalErrorActionTimeout;
	AutomaticStartAction=$_.AutomaticStartAction;
	AutomaticStartDelay=$_.AutomaticStartDelay;
	AutomaticStopAction=$_.AutomaticStopAction;
	CheckpointType=$_.CheckpointType;
	DynamicMemory=$_.DynamicMemory;
	GuestControlledCacheTypes=$_.GuestControlledCacheTypes;
	HighMemoryMappedIoSpace=$_.HighMemoryMappedIoSpace;
	LockOnDisconnect=$_.LockOnDisconnect;
	LowMemoryMappedIoSpace=$_.LowMemoryMappedIoSpace;
	MemoryMaximumBytes=$_.MemoryMaximumBytes;
	MemoryMinimumBytes=$_.MemoryMinimumBytes;
	MemoryStartupBytes=$_.MemoryStartupBytes;
	Notes=$_.Notes;
	ProcessorCount=$_.ProcessorCount;
	SmartPagingFilePath=$_.SmartPagingFilePath;
	SnapshotFileLocation=$_.SnapshotFileLocation;
	StaticMemory=$_.StaticMemory;
}} | ConvertTo-Json
`))

func (c *HypervClient) GetVM(name string) (result vm, err error) {
	err = c.runScriptWithResult(getVMTemplate, getVMArgs{
		Name:name,
	}, &result);

	return result, err
}

type updateVMArgs struct {
	VmJson		string
}

var updateVMTemplate = template.Must(template.New("UpdateVM").Parse(`
$vm = '{{.VmJson}}' | ConvertFrom-Json
$automaticCriticalErrorAction = [Microsoft.HyperV.PowerShell.CriticalErrorAction]$vm.AutomaticCriticalErrorAction
$automaticStartAction = [Microsoft.HyperV.PowerShell.StartAction]$vm.AutomaticStartAction
$automaticStopAction = [Microsoft.HyperV.PowerShell.StopAction]$vm.AutomaticStopAction
$checkpointType = [Microsoft.HyperV.PowerShell.CheckpointType]$vm.CheckpointType
$lockOnDisconnect = [Microsoft.HyperV.PowerShell.OnOffState]$vm.LockOnDisconnect

Set-Vm -Name $vm.Name -GuestControlledCacheTypes:$vm.GuestControlledCacheTypes -LowMemoryMappedIoSpace $vm.LowMemoryMappedIoSpace -HighMemoryMappedIoSpace $vm.HighMemoryMappedIoSpace -ProcessorCount $vm.ProcessorCount -DynamicMemory:$vm.DynamicMemory -StaticMemory:$vm.StaticMemory -MemoryMinimumBytes $vm.MemoryMinimumBytes -MemoryMaximumBytes $vm.MemoryMaximumBytes -MemoryStartupBytes $vm.MemoryStartupBytes -AutomaticStartAction $automaticStartAction -AutomaticStopAction $automaticStopAction -AutomaticStartDelay $vm.AutomaticStartDelay -AutomaticCriticalErrorAction $automaticCriticalErrorAction -AutomaticCriticalErrorActionTimeout $vm.AutomaticCriticalErrorActionTimeout -LockOnDisconnect $lockOnDisconnect -Notes $vm.Notes -SnapshotFileLocation $vm.SnapshotFileLocation -SmartPagingFilePath $vm.SmartPagingFilePath -CheckpointType $checkpointType -AllowUnverifiedPaths $vm.AllowUnverifiedPaths
`))

func (c *HypervClient) UpdateVM(
	name string,
//	generation int,
	allowUnverifiedPaths bool,
	automaticCriticalErrorAction CriticalErrorAction,
	automaticCriticalErrorActionTimeout int32,
	automaticStartAction StartAction,
	automaticStartDelay int32,
	automaticStopAction StopAction,
	checkpointType CheckpointType,
	dynamicMemory bool,
	guestControlledCacheTypes bool,
	highMemoryMappedIoSpace int64,
	lockOnDisconnect OnOffState,
	lowMemoryMappedIoSpace int32,
	memoryMaximumBytes int64,
	memoryMinimumBytes int64,
	memoryStartupBytes int64,
	notes string,
	processorCount int64,
	smartPagingFilePath string,
	snapshotFileLocation string,
	staticMemory bool,
) (err error) {

	vmJson, err := json.Marshal(vm{
		Name:name,
		//Generation:generation,
		AllowUnverifiedPaths: allowUnverifiedPaths,
		AutomaticCriticalErrorAction: automaticCriticalErrorAction,
		AutomaticCriticalErrorActionTimeout: automaticCriticalErrorActionTimeout,
		AutomaticStartAction: automaticStartAction,
		AutomaticStartDelay: automaticStartDelay,
		AutomaticStopAction: automaticStopAction,
		CheckpointType: checkpointType,
		DynamicMemory: dynamicMemory,
		GuestControlledCacheTypes: guestControlledCacheTypes,
		HighMemoryMappedIoSpace: highMemoryMappedIoSpace,
		LockOnDisconnect: lockOnDisconnect,
		LowMemoryMappedIoSpace: lowMemoryMappedIoSpace,
		MemoryMaximumBytes: memoryMaximumBytes,
		MemoryMinimumBytes: memoryMinimumBytes,
		MemoryStartupBytes: memoryStartupBytes,
		Notes: notes,
		ProcessorCount: processorCount,
		SmartPagingFilePath: smartPagingFilePath,
		SnapshotFileLocation: snapshotFileLocation,
		StaticMemory: staticMemory,
	})

	err = c.runFireAndForgetScript(updateVMTemplate, updateVMArgs{
		VmJson:string(vmJson),
	});

	return err
}

type deleteVMArgs struct {
	Name		string
}

var deleteVMTemplate = template.Must(template.New("DeleteVM").Parse(`
Get-VM | ?{$_.Name -eq '{{.Name}}'} | Remove-VM
`))

func (c *HypervClient) DeleteVM(name string) (err error) {
	err = c.runFireAndForgetScript(deleteVMTemplate, deleteVMArgs{
		Name:name,
	});

	return err
}