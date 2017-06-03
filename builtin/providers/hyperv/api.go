package hyperv

import (
	"text/template"
	"encoding/json"
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform/communicator/winrm"
	"github.com/hashicorp/terraform/builtin/providers/hyperv/powershell"
	"strings"
)

type client struct {
	Communicator 		*winrm.Communicator
	ElevatedUser            string
	ElevatedPassword	string
}

func (c *client) runFireAndForgetScript(script  *template.Template, args interface{})(error){
	var scriptRendered bytes.Buffer
	err := script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := string(scriptRendered.Bytes())

	exited, exitStatus, _, stderr, err := powershell.RunCommand(c.Communicator, c.ElevatedUser, c.ElevatedPassword, "", command)

	if err != nil {
		return err
	}

	if !exited {
		return fmt.Errorf("Command did not execute completly")
	}

	if exitStatus != 0 {
		return fmt.Errorf("Command exit code not expected: %s", exitStatus)
	}

	stderr = strings.TrimSpace(stderr)

	if len(stderr) > 0 {
		return fmt.Errorf("Command stderr: %s", stderr)
	}

	return nil
}

func (c *client) runScriptWithResult(script  *template.Template, args interface{}, result interface{})(err error){
	var scriptRendered bytes.Buffer
	err = script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := string(scriptRendered.Bytes())

	exited, exitStatus, stdout, stderr, err := powershell.RunCommand(c.Communicator, c.ElevatedUser, c.ElevatedPassword, "", command)

	if err != nil {
		return err
	}

	if !exited {
		return fmt.Errorf("Command did not execute completly")
	}

	if exitStatus != 0 {
		return fmt.Errorf("Command exit code not expected: %s", exitStatus)
	}

	stderr = strings.TrimSpace(stderr)

	if len(stderr) > 0 {
		return fmt.Errorf("Command stderr: %s", stderr)
	}

	stdout = strings.TrimSpace(stdout)

	return json.Unmarshal([]byte(stdout), &result)
}

type vmSwitch struct {
	Name		string
	Notes		string
	AllowManagementOS bool
	EmbeddedTeamingEnabled bool
	IovEnabled bool
	PacketDirectEnabled bool
	BandwidthReservationMode int
	SwitchType	int
	NetAdapterInterfaceDescriptions []string
}

type createVMSwitchArgs struct {
	VmSwitchJson		string
}

var createVMSwitchTemplate = template.Must(template.New("CreateVMSwitch").Parse(`
$vmSwitch = '{{.VmSwitchJson}}' | ConvertFrom-Json
$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]$vmSwitch.BandwidthReservationMode
$switchType = [Microsoft.HyperV.PowerShell.VMSwitchType]$vmSwitch.SwitchType

if ($vmSwitch.NetAdapterInterfaceDescriptions) {
	New-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -AllowManagementOS $vmSwitch.AllowManagementOS -EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled -EnableIov $vmSwitch.IovEnabled -EnablePacketDirect $vmSwitch.PacketDirectEnabled -MinimumBandwidthMode $minimumBandwidthMode -NetAdapterInterfaceDescription $vmSwitch.NetAdapterInterfaceDescriptions
} else {
	New-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -AllowManagementOS $vmSwitch.AllowManagementOS -EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled -EnableIov $vmSwitch.IovEnabled -EnablePacketDirect $vmSwitch.PacketDirectEnabled -MinimumBandwidthMode $minimumBandwidthMode -SwitchType $switchType
}
`))

func (c *client) CreateVMSwitch(name string,
	notes string,
	allowManagementOS bool,
	embeddedTeamingEnabled bool,
	iovEnabled bool,
	packetDirectEnabled bool,
	bandwidthReservationMode int,
	switchType int,
	netAdapterInterfaceDescriptions []string) (err error) {

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
(Get-VMSwitch -name '{{.SwitchName}}') | %{ @{Name=$_.Name;Notes=$_.Notes;AllowManagementOS=$_.AllowManagementOS;EmbeddedTeamingEnabled=$_.EmbeddedTeamingEnabled;IovEnabled=$_.IovEnabled;PacketDirectEnabled=$_.PacketDirectEnabled;BandwidthReservationMode=$_.BandwidthReservationMode;SwitchType=$_.SwitchType;NetAdapterInterfaceDescriptions=$_.NetAdapterInterfaceDescriptions}} | ConvertTo-Json
`))

func (c *client) GetVMSwitch(name string) (result vmSwitch, err error) {
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

if ($vmSwitch.NetAdapterInterfaceDescriptions) {
	Set-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -AllowManagementOS $vmSwitch.AllowManagementOS -NetAdapterInterfaceDescription $vmSwitch.NetAdapterInterfaceDescriptions

	#Updates not supported on:
	#-EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled
	#-EnableIov $vmSwitch.IovEnabled
	#-EnablePacketDirect $vmSwitch.PacketDirectEnabled
	#-MinimumBandwidthMode $minimumBandwidthMode

} else {
	Set-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -AllowManagementOS $vmSwitch.AllowManagementOS -SwitchType $switchType

	#Updates not supported on:
	#-EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled
	#-EnableIov $vmSwitch.IovEnabled
	#-EnablePacketDirect $vmSwitch.PacketDirectEnabled
	#-MinimumBandwidthMode $minimumBandwidthMode
}
`))

func (c *client) UpdateVMSwitch(name string,
	notes string,
	allowManagementOS bool,
	//embeddedTeamingEnabled bool,
	//iovEnabled bool,
	//packetDirectEnabled bool,
	//bandwidthReservationMode int,
	switchType int,
	netAdapterInterfaceDescriptions []string) (err error) {

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
Get-VMSwitch | ?{$_.Name -eq '{{.SwitchName}}'} | Remove-VMSwitch
`))

func (c *client) DeleteVMSwitch(name string) (err error) {
	err = c.runFireAndForgetScript(deleteVMSwitchTemplate, deleteVMSwitchArgs{
		Name:name,
	});

	return err
}