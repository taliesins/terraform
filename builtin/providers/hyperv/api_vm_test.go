package hyperv

import (
	"testing"
	"encoding/json"
)

func TestSerializeVm(t *testing.T) {
	vmJson, err := json.Marshal(vm{
		Name:"test",
		Notes:"test notes",
	})

	if err != nil {
		t.Errorf("Unable to deserialize vm: %s", err.Error())
	}

	vmJsonString := string(vmJson)

	if vmJsonString == "" {
		t.Errorf("Unable to deserialize vm: %s", err.Error())
	}
}

func TestDeserializeVm(t *testing.T){
	var vmJson = `
{
    "BandwidthReservationMode":  2,
    "NetAdapterInterfaceDescriptions":  [
                                            "Dell Wireless 1830 802.11ac"
                                        ],
    "Notes":  "test notes",
    "AllowManagementOS":  true,
    "Name":  "test",
    "SwitchType":  2,
    "IovEnabled":  false,
    "EmbeddedTeamingEnabled":  false,
    "PacketDirectEnabled":  false
}
`

	var vm vm
	err := json.Unmarshal([]byte(vmJson), &vm)

	if err != nil {
		t.Errorf("Unable to deserialize vm: %s", err.Error())
	}
}