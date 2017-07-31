package hyperv

import (
	"text/template"
	"encoding/json"
)

type VMSwitchBandwidthMode int

const (
	VMSwitchBandwidthMode_Default VMSwitchBandwidthMode = 0
	VMSwitchBandwidthMode_Weight VMSwitchBandwidthMode = 1
	VMSwitchBandwidthMode_Absolute VMSwitchBandwidthMode = 2
	VMSwitchBandwidthMode_None VMSwitchBandwidthMode = 3
)

type VMSwitchType int

const (
	VMSwitchType_Private VMSwitchType = 0
	VMSwitchType_Internal VMSwitchType = 1
	VMSwitchType_External VMSwitchType = 2
)


type vmSwitch struct {
	Name		string
	Notes		string
	AllowManagementOS bool
	EmbeddedTeamingEnabled bool
	IovEnabled bool
	PacketDirectEnabled bool
	BandwidthReservationMode VMSwitchBandwidthMode
	SwitchType	VMSwitchType
	NetAdapterInterfaceDescriptions []string
	NetAdapterNames []string
	DefaultFlowMinimumBandwidthAbsolute int64
	DefaultFlowMinimumBandwidthWeight int64
	DefaultQueueVmmqEnabled bool
	DefaultQueueVmmqQueuePairs int32
	DefaultQueueVrssEnabled bool
}

type createVMSwitchArgs struct {
	VmSwitchJson		string
}

var createVMSwitchTemplate = template.Must(template.New("CreateVMSwitch").Parse(`
$vmSwitch = '{{.VmSwitchJson}}' | ConvertFrom-Json
$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]$vmSwitch.BandwidthReservationMode
$switchType = [Microsoft.HyperV.PowerShell.VMSwitchType]$vmSwitch.SwitchType
$NetAdapterInterfaceDescriptions = @($vmSwitch.NetAdapterInterfaceDescriptions)
$NetAdapterNames = @($vmSwitch.$NetAdapterNames)
#when EnablePacketDirect=true it seems to throw an exception if EnableIov=true or EnableEmbeddedTeaming=true

if ($NetAdapterInterfaceDescriptions || $NetAdapterNames) {
	New-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -AllowManagementOS $vmSwitch.AllowManagementOS -EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled -EnableIov $vmSwitch.IovEnabled -EnablePacketDirect $vmSwitch.PacketDirectEnabled -MinimumBandwidthMode $minimumBandwidthMode -NetAdapterInterfaceDescription $NetAdapterInterfaceDescriptions -NetAdapterName $NetAdapterNames
} else {
	New-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled -EnableIov $vmSwitch.IovEnabled -EnablePacketDirect $vmSwitch.PacketDirectEnabled -MinimumBandwidthMode $minimumBandwidthMode -SwitchType $switchType

	#not used unless interface is specified
	#-AllowManagementOS $vmSwitch.AllowManagementOS
}

Set-VMSwitch -Name $vmSwitch.Name -DefaultFlowMinimumBandwidthAbsolute $vmSwitch.DefaultFlowMinimumBandwidthAbsolute -DefaultFlowMinimumBandwidthWeight $vmSwitch.DefaultFlowMinimumBandwidthWeight -DefaultQueueVmmqEnabled $vmSwitch.DefaultQueueVmmqEnabled -DefaultQueueVmmqQueuePairs $vmSwitch.DefaultQueueVmmqQueuePairs -DefaultQueueVrssEnabled $vmSwitch.DefaultQueueVrssEnabled
`))

func (c *HypervClient) CreateVMSwitch(
	name string,
	notes string,
	allowManagementOS bool,
	embeddedTeamingEnabled bool,
	iovEnabled bool,
	packetDirectEnabled bool,
	bandwidthReservationMode VMSwitchBandwidthMode,
	switchType VMSwitchType,
	netAdapterInterfaceDescriptions []string,
	netAdapterNames []string,
	defaultFlowMinimumBandwidthAbsolute int64,
	defaultFlowMinimumBandwidthWeight int64,
	defaultQueueVmmqEnabled bool,
	defaultQueueVmmqQueuePairs int32,
	defaultQueueVrssEnabled bool,
) (err error) {

	vmSwitchJson, err := json.Marshal(vmSwitch{
		Name:name,
		Notes:notes,
		AllowManagementOS:allowManagementOS,
		EmbeddedTeamingEnabled:embeddedTeamingEnabled,
		IovEnabled:iovEnabled,
		PacketDirectEnabled:packetDirectEnabled,
		BandwidthReservationMode:bandwidthReservationMode,
		SwitchType:switchType,
		NetAdapterInterfaceDescriptions:netAdapterInterfaceDescriptions,
		NetAdapterNames:netAdapterNames,
		DefaultFlowMinimumBandwidthAbsolute:defaultFlowMinimumBandwidthAbsolute,
		DefaultFlowMinimumBandwidthWeight:defaultFlowMinimumBandwidthWeight,
		DefaultQueueVmmqEnabled:defaultQueueVmmqEnabled,
		DefaultQueueVmmqQueuePairs:defaultQueueVmmqQueuePairs,
		DefaultQueueVrssEnabled:defaultQueueVrssEnabled,
	})

	err = c.runFireAndForgetScript(createVMSwitchTemplate, createVMSwitchArgs{
		VmSwitchJson:string(vmSwitchJson),
	});

	return err
}

type getVMSwitchArgs struct {
	Name		string
}

var getVMSwitchTemplate = template.Must(template.New("GetVMSwitch").Parse(`
(Get-VMSwitch -name '{{.Name}}') | %{ @{
	Name=$_.Name;
	Notes=$_.Notes;
	AllowManagementOS=$_.AllowManagementOS;
	EmbeddedTeamingEnabled=$_.EmbeddedTeamingEnabled;
	IovEnabled=$_.IovEnabled;
	PacketDirectEnabled=$_.PacketDirectEnabled;
	BandwidthReservationMode=$_.BandwidthReservationMode;
	SwitchType=$_.SwitchType;
	NetAdapterInterfaceDescriptions=$_.NetAdapterInterfaceDescriptions;
	NetAdapterNames=@(Get-NetAdapter -InterfaceDescription $_.NetAdapterInterfaceDescriptions | %{$_.Name});
	DefaultFlowMinimumBandwidthAbsolute=$_.DefaultFlowMinimumBandwidthAbsolute;
	DefaultFlowMinimumBandwidthWeight=$_.DefaultFlowMinimumBandwidthWeight;
	DefaultQueueVmmqEnabled=$_.DefaultQueueVmmqEnabled;
	DefaultQueueVmmqQueuePairs=$_.DefaultQueueVmmqQueuePairs;
	DefaultQueueVrssEnabled=$_.DefaultQueueVrssEnabled;
}} | ConvertTo-Json
`))

func (c *HypervClient) GetVMSwitch(name string) (result vmSwitch, err error) {
	err = c.runScriptWithResult(getVMSwitchTemplate, getVMSwitchArgs{
		Name:name,
	}, &result);

	return result, err
}

type updateVMSwitchArgs struct {
	VmSwitchJson		string
}

var updateVMSwitchTemplate = template.Must(template.New("UpdateVMSwitch").Parse(`
$vmSwitch = '{{.VmSwitchJson}}' | ConvertFrom-Json
$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]$vmSwitch.BandwidthReservationMode
$switchType = [Microsoft.HyperV.PowerShell.VMSwitchType]$vmSwitch.SwitchType
$NetAdapterInterfaceDescriptions = @($vmSwitch.NetAdapterInterfaceDescriptions)
$NetAdapterNames = @($vmSwitch.$NetAdapterNames)
#when EnablePacketDirect=true it seems to throw an exception if EnableIov=true or EnableEmbeddedTeaming=true

if ($NetAdapterInterfaceDescriptions || $NetAdapterNames) {
	Set-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -AllowManagementOS $vmSwitch.AllowManagementOS -NetAdapterInterfaceDescription $vmSwitch.NetAdapterInterfaceDescriptions -NetAdapterName $NetAdapterNames -DefaultFlowMinimumBandwidthAbsolute $vmSwitch.DefaultFlowMinimumBandwidthAbsolute -DefaultFlowMinimumBandwidthWeight $vmSwitch.DefaultFlowMinimumBandwidthWeight -DefaultQueueVmmqEnabled $vmSwitch.DefaultQueueVmmqEnabled -DefaultQueueVmmqQueuePairs $vmSwitch.DefaultQueueVmmqQueuePairs -DefaultQueueVrssEnabled $vmSwitch.DefaultQueueVrssEnabled

	#Updates not supported on:
	#-EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled
	#-EnableIov $vmSwitch.IovEnabled
	#-EnablePacketDirect $vmSwitch.PacketDirectEnabled
	#-MinimumBandwidthMode $minimumBandwidthMode

} else {
	Set-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -AllowManagementOS $vmSwitch.AllowManagementOS -SwitchType $switchType -DefaultFlowMinimumBandwidthAbsolute $vmSwitch.DefaultFlowMinimumBandwidthAbsolute -DefaultFlowMinimumBandwidthWeight $vmSwitch.DefaultFlowMinimumBandwidthWeight -DefaultQueueVmmqEnabled $vmSwitch.DefaultQueueVmmqEnabled -DefaultQueueVmmqQueuePairs $vmSwitch.DefaultQueueVmmqQueuePairs -DefaultQueueVrssEnabled $vmSwitch.DefaultQueueVrssEnabled

	#Updates not supported on:
	#-EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled
	#-EnableIov $vmSwitch.IovEnabled
	#-EnablePacketDirect $vmSwitch.PacketDirectEnabled
	#-MinimumBandwidthMode $minimumBandwidthMode
}
`))

func (c *HypervClient) UpdateVMSwitch(
	name string,
	notes string,
	allowManagementOS bool,
//embeddedTeamingEnabled bool,
//iovEnabled bool,
//packetDirectEnabled bool,
//bandwidthReservationMode VMSwitchBandwidthMode,
	switchType VMSwitchType,
	netAdapterInterfaceDescriptions []string,
	netAdapterNames []string,
	defaultFlowMinimumBandwidthAbsolute int64,
	defaultFlowMinimumBandwidthWeight int64,
	defaultQueueVmmqEnabled bool,
	defaultQueueVmmqQueuePairs int32,
	defaultQueueVrssEnabled bool,
) (err error) {

	vmSwitchJson, err := json.Marshal(vmSwitch{
		Name:name,
		Notes:notes,
		AllowManagementOS:allowManagementOS,
		//EmbeddedTeamingEnabled:embeddedTeamingEnabled,
		//IovEnabled:iovEnabled,
		//PacketDirectEnabled:packetDirectEnabled,
		//BandwidthReservationMode:bandwidthReservationMode,
		SwitchType:switchType,
		NetAdapterInterfaceDescriptions:netAdapterInterfaceDescriptions,
		NetAdapterNames:netAdapterNames,
		DefaultFlowMinimumBandwidthAbsolute:defaultFlowMinimumBandwidthAbsolute,
		DefaultFlowMinimumBandwidthWeight:defaultFlowMinimumBandwidthWeight,
		DefaultQueueVmmqEnabled:defaultQueueVmmqEnabled,
		DefaultQueueVmmqQueuePairs:defaultQueueVmmqQueuePairs,
		DefaultQueueVrssEnabled:defaultQueueVrssEnabled,
	})

	err = c.runFireAndForgetScript(updateVMSwitchTemplate, updateVMSwitchArgs{
		VmSwitchJson:string(vmSwitchJson),
	});

	return err
}

type deleteVMSwitchArgs struct {
	Name		string
}

var deleteVMSwitchTemplate = template.Must(template.New("DeleteVMSwitch").Parse(`
Get-VMSwitch | ?{$_.Name -eq '{{.Name}}'} | Remove-VMSwitch
`))

func (c *HypervClient) DeleteVMSwitch(name string) (err error) {
	err = c.runFireAndForgetScript(deleteVMSwitchTemplate, deleteVMSwitchArgs{
		Name:name,
	});

	return err
}
