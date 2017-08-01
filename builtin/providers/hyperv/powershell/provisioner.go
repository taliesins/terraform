package powershell

import (
	"bytes"
	"fmt"
	"log"
	"github.com/hashicorp/terraform/communicator/remote"
	"github.com/hashicorp/terraform/communicator/winrm"
	"io/ioutil"
	"os"
	"bufio"
	"time"
	"crypto/rand"
)

// Generates a time ordered UUID. Top 32 bits are a timestamp,
// bottom 96 are random.
func timeOrderedUUID() string {
	unix := uint32(time.Now().UTC().Unix())

	b := make([]byte, 12)
	n, err := rand.Read(b)
	if n != len(b) {
		err = fmt.Errorf("Not enough entropy available")
	}
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%04x%08x",
		unix, b[0:2], b[2:4], b[4:6], b[6:8], b[8:])
}

func createCommand(vars string, remotePath string) (commandText string, err error) {
	var executeCommandTemplateRendered bytes.Buffer

	err = executeCommandTemplate.Execute(&executeCommandTemplateRendered, executeCommandTemplateOptions{
		Vars: vars,
		Path: remotePath,
	})

	if err != nil {
		fmt.Printf("Error creating command template: %s", err)
		return "", err
	}

	commandText = string(executeCommandTemplateRendered.Bytes())

	return commandText, err
}

func createElevatedCommand(comm *winrm.Communicator, elevatedUser string, elevatedPassword string, vars string, remotePath string) (commandText string, err error) {
	commandText, err = createCommand(vars, remotePath)
	if err != nil {
		fmt.Printf("Error creating command template: %s", err)
		return "", err
	}

	elevatedRemotePath, err := generateElevatedRunner(comm, elevatedUser, elevatedPassword, remotePath)
	if err != nil {
		return "", fmt.Errorf("Error generating elevated runner: %s", err)
	}

	return createCommand(vars, elevatedRemotePath)
}

func generateElevatedRunner(comm *winrm.Communicator, elevatedUser string, elevatedPassword string, remotePath string) (elevatedRemotePath string, err error) {
	log.Printf("Building elevated command wrapper for: %s", remotePath)

	name := fmt.Sprintf("terraform-%s", timeOrderedUUID())
	fileName := fmt.Sprintf(`elevated-shell-%s.ps1`, name)

	var elevatedCommandTemplateRendered bytes.Buffer
	err = elevatedCommandTemplate.Execute(&elevatedCommandTemplateRendered, elevatedCommandTemplateOptions{
		User:            		elevatedUser,
		Password:        		elevatedPassword,
		TaskDescription: 		"Terraform elevated task",
		TaskName:        		name,
		TaskExecutionTimeLimit: "PT2H",
		ScriptPath: 			remotePath,
	})

	if err != nil {
		fmt.Printf("Error creating elevated command template: %s", err)
		return "", err
	}

	elevatedCommand := string(elevatedCommandTemplateRendered.Bytes())

	elevatedRemotePath, err = uploadScript(comm, fileName, elevatedCommand)
	if err != nil {
		return "", err
	}

	return elevatedRemotePath, nil
}

func uploadScript(comm *winrm.Communicator, fileName string, command string) (remoteAbsolutePath string, err error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), fileName)
	writer := bufio.NewWriter(tmpFile)
	if _, err := writer.WriteString(command); err != nil {
		return "", fmt.Errorf("Error preparing shell script: %s", err)
	}

	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("Error preparing shell script: %s", err)
	}
	tmpFile.Close()
	f, err := os.Open(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("Error opening temporary shell script: %s", err)
	}
	defer f.Close()
	defer os.Remove(tmpFile.Name())

	remotePath := fmt.Sprintf(`%s\%s`, `$env:TEMP`, fileName)

	log.Printf("Uploading shell wrapper for command from [%s] to [%s] ", tmpFile.Name(), remotePath)
	remoteAbsolutePath, err = comm.UploadScript(remotePath, f)
	if err != nil {
		return "", fmt.Errorf("Error uploading shell script: %s", err)
	}

	return remoteAbsolutePath, nil
}

//Run powers
func RunCommand(comm *winrm.Communicator, elevatedUser string, elevatedPassword string, vars string, commandText string) (exited bool, exitStatus int, stdout string, stderr string, err error) {
	name := fmt.Sprintf("terraform-%s", timeOrderedUUID())
	fileName := fmt.Sprintf(`shell-%s.ps1`, name)

	path, err := uploadScript(comm, fileName, commandText)
	if err != nil {
		return false, 0, "", "", err
	}

	var command string

	if elevatedUser == "" {
		command, err = createCommand(vars, path)
	} else {
		command, err = createElevatedCommand(comm, elevatedUser, elevatedPassword, vars, path)
	}

	if err != nil {
		return false, 0, "", "", err
	}

	var executeCommandFromCommandLineTemplateRendered bytes.Buffer
	err = executeCommandFromCommandLineTemplate.Execute(&executeCommandFromCommandLineTemplateRendered, executeCommandFromCommandLineTemplateOptions{
		Powershell: command,
	})

	if err != nil {
		return false, 0, "", "", err
	}

	commandLine := string(executeCommandFromCommandLineTemplateRendered.Bytes())

	var cmd remote.Cmd
	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)
	cmd.Command = commandLine
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	log.Printf("Executing run command shell wrapper with: %s", commandLine)

	err = comm.Start(&cmd)

	if err != nil {
		return false, 0, "", "", fmt.Errorf("error executing remote command: %s", err)
	}

	cmd.Wait()

	return cmd.Exited, cmd.ExitStatus, stdoutBuffer.String(), stderrBuffer.String(), nil
}