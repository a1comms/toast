// +build windows

package toast

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	uuid "github.com/satori/go.uuid"

	"syscall"
)

var toastTemplate *template.Template

func init() {
	toastTemplate = template.New("toast")
	toastTemplate.Parse(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '{{if .AppID}}{{.AppID}}{{else}}Windows App{{end}}'

$template = @"
<toast activationType="{{.ActivationType}}" launch="{{.ActivationArguments}}" duration="{{.Duration}}">
    <visual>
        <binding template="ToastGeneric">
			{{if .HeroImage}}
			<image placement="hero" src="{{.HeroImage}}" />
			{{end}}
            {{if .Icon}}
            <image placement="appLogoOverride" src="{{.Icon}}" />
            {{end}}
			{{if .Image}}
            <image src="{{.Image}}" />
            {{end}}
            {{if .Title}}
            <text><![CDATA[{{.Title}}]]></text>
            {{end}}
            {{if .Message}}
            <text><![CDATA[{{.Message}}]]></text>
            {{end}}
			{{if .Attribution}}
			<text placement="attribution">{{.Attribution}}</text>
			{{end}}
        </binding>
    </visual>
    {{if ne .Audio "silent"}}
	<audio src="{{.Audio}}" loop="{{.Loop}}" />
	{{else}}
	<audio silent="true" />
	{{end}}
    {{if .Actions}}
    <actions>
        {{range .Actions}}
        <action activationType="{{.Type}}" content="{{.Label}}" arguments="{{.Arguments}}" />
        {{end}}
    </actions>
    {{end}}
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
    `)
}

func (n *Notification) applyDefaults() {
	if n.ActivationType == "" {
		n.ActivationType = "protocol"
	}
	if n.Duration == "" {
		n.Duration = Short
	}
	if n.Audio == "" {
		n.Audio = Default
	}
}

func (n *Notification) buildXML() (string, error) {
	var out bytes.Buffer
	err := toastTemplate.Execute(&out, n)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Builds the Windows PowerShell script & invokes it, causing the toast to display.
//
// Note: Running the PowerShell script is by far the slowest process here, and can take a few
// seconds in some cases.
//
//     notification := toast.Notification{
//         AppID: "Example App",
//         Title: "My notification",
//         Message: "Some message about how important something is...",
//         Icon: "go.png",
//         Actions: []toast.Action{
//             {"protocol", "I'm a button", ""},
//             {"protocol", "Me too!", ""},
//         },
//     }
//     err := notification.Push()
//     if err != nil {
//         log.Fatalln(err)
//     }
func (n *Notification) Push() error {
	n.applyDefaults()
	xml, err := n.buildXML()
	if err != nil {
		return err
	}
	return invokeTemporaryScript(xml)
}

func invokeTemporaryScript(content string) error {
	id := uuid.NewV4()
	file := filepath.Join(os.TempDir(), id.String()+".ps1")
	defer os.Remove(file)
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	out := append(bomUtf8, []byte(content)...)
	err := ioutil.WriteFile(file, out, 0600)
	if err != nil {
		return err
	}
	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", "-File", file)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}
