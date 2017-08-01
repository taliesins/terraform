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

type HypervClient struct {
	Communicator 		*winrm.Communicator
	ElevatedUser        string
	ElevatedPassword	string
	Vars				string
}

func (c *HypervClient) runFireAndForgetScript(script  *template.Template, args interface{})(error){
	var scriptRendered bytes.Buffer
	err := script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := string(scriptRendered.Bytes())

	exited, exitStatus, stdout, stderr, err := powershell.RunCommand(c.Communicator, c.ElevatedUser, c.ElevatedPassword, c.Vars, command)

	if err != nil {
		return err
	}

	if !exited {
		return fmt.Errorf("Command did not execute completly: \nstderr:\n%s\nstdout:\n%s\nvars:\n%s\ncommand:\n%s", stderr, stdout, c.Vars, command)
	}

	if exitStatus != 0 {
		return fmt.Errorf("Command exit code not expected: %d\nstderr:\n%s\nstdout:\n%s\nvars:\n%s\ncommand:\n%s", exitStatus, stderr, stdout, c.Vars, command)
	}

	stderr = strings.TrimSpace(stderr)

	if len(stderr) > 0 {
		return fmt.Errorf("Command stderr: \nstderr:\n%s\nstdout:\n%s\nvars:\n%s\ncommand:\n%s", stderr, stdout, c.Vars, command)
	}

	return nil
}

func (c *HypervClient) runScriptWithResult(script  *template.Template, args interface{}, result interface{})(err error){
	var scriptRendered bytes.Buffer
	err = script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := string(scriptRendered.Bytes())

	exited, exitStatus, stdout, stderr, err := powershell.RunCommand(c.Communicator, c.ElevatedUser, c.ElevatedPassword, c.Vars, command)

	if err != nil {
		return err
	}

	if !exited {
		return fmt.Errorf("Command did not execute completly: \nstderr:\n%s\nstdout:\n%s\nvars:\n%s\ncommand:\n%s", stderr, stdout, c.Vars, command)
	}

	if exitStatus != 0 {
		return fmt.Errorf("Command exit code not expected: %d\nstderr:\n%s\nstdout:\n%s\nvars:\n%s\ncommand:\n%s", exitStatus, stderr, stdout, c.Vars, command)
	}

	stderr = strings.TrimSpace(stderr)

	if len(stderr) > 0 {
		return fmt.Errorf("Command stderr: \nstderr:\n%s\nstdout:\n%s\nvars:\n%s\ncommand:\n%s", stderr, stdout, c.Vars, command)
	}

	stdout = strings.TrimSpace(stdout)

	return json.Unmarshal([]byte(stdout), &result)
}

