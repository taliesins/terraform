package powershell

import (
	"text/template"
)

type executeCommandTemplateOptions struct {
	Vars            string
	Path        	string
}

var executeCommandTemplate = template.Must(template.New("ExecuteCommand").Parse(`if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue'};{{.Vars}}&'{{.Path}}';exit $LastExitCode`))

type elevatedCommandTemplateOptions struct {
	User            		string
	Password        		string
	TaskName        		string
	TaskDescription 		string
	TaskExecutionTimeLimit 	string
	Vars            		string
	ScriptPath  			string
}

var elevatedCommandTemplate = template.Must(template.New("ElevatedCommand").Parse(`
function GetTempFile($fileName) {
  $path = $env:TEMP
  if (!$path){
    $path = 'c:\windows\Temp\'
  }
  return Join-Path -Path $path -ChildPath $fileName
}
function SlurpStdout($outFile, $currentLine) {
  if (Test-Path $outFile) {
    get-content $outFile | select -skip $currentLine | %{
      $currentLine += 1
      Write-Host "$_"
    }
  }
  return $currentLine
}

function SanitizeFileName($fileName) {
    return $fileName.Replace(' ', '_').Replace('&', 'and').Replace('{', '(').Replace('}', ')').Replace('~', '-').Replace('#', '').Replace('%', '')
}

function RunAsScheduledTask($username, $password, $taskName, $taskDescription, $taskExecutionTimeLimit, $vars, $scriptFile)
{
  $stdoutFile = GetTempFile("$(SanitizeFileName($taskName))_stdout.log")
  if (Test-Path $stdoutFile) {
    Remove-Item $stdoutFile | Out-Null
  }
  $taskXml = @'
<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
    <RegistrationInfo>
	    <Description>{taskDescription}</Description>
    </RegistrationInfo>
    <Principals>
        <Principal id="Author">
        <UserId>{username}</UserId>
        <LogonType>Password</LogonType>
        <RunLevel>HighestAvailable</RunLevel>
        </Principal>
    </Principals>
    <Settings>
        <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
        <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
        <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
        <AllowHardTerminate>true</AllowHardTerminate>
        <StartWhenAvailable>false</StartWhenAvailable>
        <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
        <IdleSettings>
        <StopOnIdleEnd>false</StopOnIdleEnd>
        <RestartOnIdle>false</RestartOnIdle>
        </IdleSettings>
        <AllowStartOnDemand>true</AllowStartOnDemand>
        <Enabled>true</Enabled>
        <Hidden>false</Hidden>
        <RunOnlyIfIdle>false</RunOnlyIfIdle>
        <WakeToRun>false</WakeToRun>
        <ExecutionTimeLimit>{taskExecutionTimeLimit}</ExecutionTimeLimit>
        <Priority>4</Priority>
    </Settings>
    <Actions Context="Author">
        <Exec>
        <Command>cmd</Command>
        <Arguments>{arguments}</Arguments>
        </Exec>
    </Actions>
</Task>
'@
  $arguments = '/c powershell "& { if (Test-Path variable:global:ProgressPreference){$ProgressPreference=''SilentlyContinue''};' + $vars+ ';'+$scriptFile+'; exit $LastExitCode }" *> $stdoutFile'
  $taskXml = $taskXml.Replace("{arguments}", $arguments.Replace('&', '&amp;').Replace('<', '&lt;').Replace('>', '&gt;').Replace('"', '&quot;').Replace('''', '&apos;'))
  $taskXml = $taskXml.Replace("{username}", $username.Replace('&', '&amp;').Replace('<', '&lt;').Replace('>', '&gt;').Replace('"', '&quot;').Replace('''', '&apos;'))
  $taskXml = $taskXml.Replace("{taskDescription}", $taskDescription.Replace('&', '&amp;').Replace('<', '&lt;').Replace('>', '&gt;').Replace('"', '&quot;').Replace('''', '&apos;'))
  $taskXml = $taskXml.Replace("{taskExecutionTimeLimit}", $taskExecutionTimeLimit.Replace('&', '&amp;').Replace('<', '&lt;').Replace('>', '&gt;').Replace('"', '&quot;').Replace('''', '&apos;'))

  $schedule = New-Object -ComObject "Schedule.Service"
  $schedule.Connect()
  $task = $schedule.NewTask($null)
  $task.XmlText = $taskXml

  $folder = $schedule.GetFolder('\')
  $folder.RegisterTaskDefinition($taskName, $task, 6, $username, $password, 1, $null) | Out-Null
  $registeredTask = $folder.GetTask("\$taskName")
  $registeredTask.Run($null) | Out-Null
  $timeout = 10
  $sec = 0
  while ((!($registeredTask.state -eq 4)) -and ($sec -lt $timeout)) {
    Start-Sleep -s 1
    $sec++
  }
  $stdoutCurrentLine = 0
  do {
    Start-Sleep -m 100
    $stdoutCurrentLine = SlurpStdout $stdoutFile $stdoutCurrentLine
  } while (!($registeredTask.state -eq 3))
  Start-Sleep -m 100
  $exit_code = $registeredTask.LastTaskResult
  $stdoutCurrentLine = SlurpStdout $stdoutFile $stdoutCurrentLine

  if (Test-Path $stdoutFile) {
    Remove-Item $stdoutFile | Out-Null
  }

  $folder.DeleteTask($taskName, 0)
  [System.Runtime.Interopservices.Marshal]::ReleaseComObject($schedule) | Out-Null

  return $exit_code
}

$username = '{{.User}}'.Replace('\.\\', $env:computername+'\')
$password = '{{.Password}}'
$taskName = '{{.TaskName}}'
$taskDescription = '{{.TaskDescription}}'
$taskExecutionTimeLimit = '{{.TaskExecutionTimeLimit}}'
$vars = '{{.Vars}}'
$scriptPath = '{{.ScriptPath}}'
$exitCode = RunAsScheduledTask -username $username -password $password -taskName $taskName -taskDescription $taskDescription -taskExecutionTimeLimit -vars $vars -scriptPath $scriptPath
exit $exitCode
`))