package powershell

import (
	"text/template"
	"strings"
)

type executeCommandFromCommandLineTemplateOptions struct {
	Powershell            string
}

var executeCommandFromCommandLineTemplate = template.Must(template.New("ExecuteCommandFromCommandLine").Funcs(template.FuncMap{
	"escapeDoubleQuotes": func(textToEscape string) string {
		return strings.Replace(textToEscape, `"`, `""`, -1)
	},
}).Parse(`powershell "{{escapeDoubleQuotes .Powershell}}"`))

type executeCommandTemplateOptions struct {
	Vars            string
	Path        	string
}

var executeCommandTemplate = template.Must(template.New("ExecuteCommand").Parse(`& { if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue'};{{.Vars}};&"{{.Path}}";exit $LastExitCode }`))

type elevatedCommandTemplateOptions struct {
	User            		string
	Password        		string
	TaskName        		string
	TaskDescription 		string
	TaskExecutionTimeLimit 	string
	Vars            		string
	ScriptPath  			string
}

var elevatedCommandTemplate = template.Must(template.New("ElevatedCommand").Funcs(template.FuncMap{
	"escapeSingleQuotes": func(textToEscape string) string {
		return strings.Replace(textToEscape, `'`, `''`, -1)
	},
}).Parse(`
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

function RunAsScheduledTask($username, $password, $taskName, $taskDescription, $taskExecutionTimeLimit, $vars, $scriptPath)
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
  $powershellToExecute = '& { if (Test-Path variable:global:ProgressPreference){$ProgressPreference=''SilentlyContinue''};' + $vars + ';&"' + $scriptPath + '";exit $LastExitCode }'
  $powershellToExecute = $powershellToExecute.Replace('"', '^"')

  $arguments = '/C powershell "' + $powershellToExecute + '" *> "' + $stdoutFile + '"'
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

$username = '{{escapeSingleQuotes .User}}'.Replace('\.\\', $env:computername+'\')
$password = '{{escapeSingleQuotes .Password}}'
$taskName = '{{escapeSingleQuotes .TaskName}}'
$taskDescription = '{{escapeSingleQuotes .TaskDescription}}'
$taskExecutionTimeLimit = '{{escapeSingleQuotes .TaskExecutionTimeLimit}}'
$vars = '{{escapeSingleQuotes .Vars}}'
$scriptPath = '{{escapeSingleQuotes .ScriptPath}}'
$exitCode = RunAsScheduledTask -username $username -password $password -taskName $taskName -taskDescription $taskDescription -taskExecutionTimeLimit $taskExecutionTimeLimit -vars $vars -scriptPath $scriptPath
exit $exitCode
`))