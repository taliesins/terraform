package powershell

import (
	"testing"
	"bytes"
)

func TestEscapeQuotesOfCommandLineTemplate(t *testing.T) {
	command := `& { if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue'};;&"C:/Windows/Temp/Test.ps1";exit $LastExitCode }`

	var executeCommandFromCommandLineTemplateRendered bytes.Buffer
	err := executeCommandFromCommandLineTemplate.Execute(&executeCommandFromCommandLineTemplateRendered, executeCommandFromCommandLineTemplateOptions{
		Powershell: command,
	})

	if err != nil {
		t.Errorf("Unable to render command line template: %s", err.Error())
	}

	commandLine := string(executeCommandFromCommandLineTemplateRendered.Bytes())

	if commandLine != `powershell "& { if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue'};;&^"C:/Windows/Temp/Test.ps1^";exit $LastExitCode }"` {
		t.Errorf("Command line template output not as expected: %s", err.Error())
	}
}