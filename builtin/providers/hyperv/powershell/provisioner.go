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

	command := string(executeCommandTemplateRendered.Bytes())

	commandText, err = generateCommandLineRunner(command)
	if err != nil {
		return "", fmt.Errorf("Error generating command line runner: %s", err)
	}

	return commandText, err
}

func createElevatedCommand(comm *winrm.Communicator, elevatedUser string, elevatedPassword string, vars string, remotePath string) (commandText string, err error) {
	var executeElevatedCommandTemplateRendered bytes.Buffer

	err = executeElevatedCommandTemplate.Execute(&executeElevatedCommandTemplateRendered, executeElevatedCommandTemplateOptions{
		Vars: vars,
		Path: remotePath,
	})

	if err != nil {
		fmt.Printf("Error creating elevated command template: %s", err)
		return "", err
	}

	command := string(executeElevatedCommandTemplateRendered.Bytes())

	// OK so we need an elevated shell runner to wrap our command, this is going to have its own path
	// generate the script and update the command runner in the process
	commandText, err = generateElevatedRunner(comm, elevatedUser, elevatedPassword, command)
	if err != nil {
		return "", fmt.Errorf("Error generating elevated runner: %s", err)
	}

	return commandText, err
}

func generateCommandLineRunner(command string) (commandText string, err error) {
	log.Printf("Building command line for: %s", command)

	base64EncodedCommand, err := powershellEncode(command)
	if err != nil {
		return "", fmt.Errorf("Error encoding command: %s", err)
	}

	commandText = "powershell -executionpolicy bypass -encodedCommand " + base64EncodedCommand

	return commandText, nil
}

func generateElevatedRunner(comm *winrm.Communicator, elevatedUser string, elevatedPassword string, command string) (commandText string, err error) {
	log.Printf("Building elevated command wrapper for: %s", command)

	// generate command
	base64EncodedCommand, err := powershellEncode(command)
	if err != nil {
		return "", fmt.Errorf("Error encoding command: %s", err)
	}

	name := fmt.Sprintf("terraform-%s", timeOrderedUUID())
	fileName := fmt.Sprintf(`elevated-shell-%s.ps1`, name)

	var elevatedCommandTemplateRendered bytes.Buffer
	err = elevatedCommandTemplate.Execute(&elevatedCommandTemplateRendered, elevatedCommandTemplateOptions{
		User:            elevatedUser,
		Password:        elevatedPassword,
		TaskDescription: "Terraform elevated task",
		TaskName:        name,
		EncodedCommand:  base64EncodedCommand,
	})

	if err != nil {
		fmt.Printf("Error creating elevated command template: %s", err)
		return "", err
	}

	elevatedCommand := string(elevatedCommandTemplateRendered.Bytes())

	path, err := uploadScript(comm, fileName, elevatedCommand)
	if err != nil {
		return "", err
	}

	//We have uploaded a wrapper file that we can execute like a standard command.
	//Vars are not needed as it will be provided in the script that has been uploaded
	return createCommand("", path)
}

func uploadScript(comm *winrm.Communicator, fileName string, command string) (path string, err error) {
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

	log.Printf("Uploading shell wrapper for command to [%s] from [%s]", path, tmpFile.Name())
	err = comm.UploadScript(path, f)
	if err != nil {
		return "", fmt.Errorf("Error uploading shell script: %s", err)
	}

	path = fmt.Sprintf(`%s\%s`, `%TEMP%`, fileName)

	return path, nil
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

	var cmd remote.Cmd
	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)
	cmd.Command = command
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	err = comm.Start(&cmd)
	if err != nil {
		return false, 0, "", "", fmt.Errorf("error executing remote command: %s", err)
	}
	cmd.Wait()

	return cmd.Exited, cmd.ExitStatus, stdoutBuffer.String(), stderrBuffer.String(), nil
}