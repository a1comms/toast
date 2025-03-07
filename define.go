package toast

import (
	"errors"
	"strings"
)

var (
	ErrorInvalidAudio    error = errors.New("toast: invalid audio")
	ErrorInvalidDuration       = errors.New("toast: invalid duration")
)

type toastAudio string

const (
	Default        toastAudio = "ms-winsoundevent:Notification.Default"
	IM                        = "ms-winsoundevent:Notification.IM"
	Mail                      = "ms-winsoundevent:Notification.Mail"
	Reminder                  = "ms-winsoundevent:Notification.Reminder"
	SMS                       = "ms-winsoundevent:Notification.SMS"
	LoopingAlarm              = "ms-winsoundevent:Notification.Looping.Alarm"
	LoopingAlarm2             = "ms-winsoundevent:Notification.Looping.Alarm2"
	LoopingAlarm3             = "ms-winsoundevent:Notification.Looping.Alarm3"
	LoopingAlarm4             = "ms-winsoundevent:Notification.Looping.Alarm4"
	LoopingAlarm5             = "ms-winsoundevent:Notification.Looping.Alarm5"
	LoopingAlarm6             = "ms-winsoundevent:Notification.Looping.Alarm6"
	LoopingAlarm7             = "ms-winsoundevent:Notification.Looping.Alarm7"
	LoopingAlarm8             = "ms-winsoundevent:Notification.Looping.Alarm8"
	LoopingAlarm9             = "ms-winsoundevent:Notification.Looping.Alarm9"
	LoopingAlarm10            = "ms-winsoundevent:Notification.Looping.Alarm10"
	LoopingCall               = "ms-winsoundevent:Notification.Looping.Call"
	LoopingCall2              = "ms-winsoundevent:Notification.Looping.Call2"
	LoopingCall3              = "ms-winsoundevent:Notification.Looping.Call3"
	LoopingCall4              = "ms-winsoundevent:Notification.Looping.Call4"
	LoopingCall5              = "ms-winsoundevent:Notification.Looping.Call5"
	LoopingCall6              = "ms-winsoundevent:Notification.Looping.Call6"
	LoopingCall7              = "ms-winsoundevent:Notification.Looping.Call7"
	LoopingCall8              = "ms-winsoundevent:Notification.Looping.Call8"
	LoopingCall9              = "ms-winsoundevent:Notification.Looping.Call9"
	LoopingCall10             = "ms-winsoundevent:Notification.Looping.Call10"
	Silent                    = "silent"
)

type toastDuration string

const (
	Short toastDuration = "short"
	Long                = "long"
)

// Notification
//
// The toast notification data. The following fields are strongly recommended;
//   - AppID
//   - Title
//
// If no toastAudio is provided, then the toast notification will be silent.
// You can set the toast to have a default audio by setting "Audio" to "toast.Default", or if your go app takes
// user-provided input for audio, call the "toast.Audio(name)" func.
//
// The AppID is shown beneath the toast message (in certain cases), and above the notification within the Action
// Center - and is used to group your notifications together. It is recommended that you provide a "pretty"
// name for your app, and not something like "com.example.MyApp".
//
// If no Title is provided, but a Message is, the message will display as the toast notification's title -
// which is a slightly different font style (heavier).
//
// The Icon should be an absolute path to the icon (as the toast is invoked from a temporary path on the user's
// system, not the working directory).
//
// If you would like the toast to call an external process/open a webpage, then you can set ActivationArguments
// to the uri you would like to trigger when the toast is clicked. For example: "https://google.com" would open
// the Google homepage when the user clicks the toast notification.
// By default, clicking the toast just hides/dismisses it.
//
// The following would show a notification to the user letting them know they received an email, and opens
// gmail.com when they click the notification. It also makes the Windows 10 "mail" sound effect.
//
//     toast := toast.Notification{
//         AppID:               "Google Mail",
//         Title:               email.Subject,
//         Message:             email.Preview,
//         Icon:                "C:/Program Files/Google Mail/icons/logo.png",
//         ActivationArguments: "https://gmail.com",
//         Audio:               toast.Mail,
//     }
//
//     err := toast.Push()
type Notification struct {
	// The name of your app. This value shows up in Windows 10's Action Centre, so make it
	// something readable for your users. It can contain spaces, however special characters
	// (eg. é) are not supported.
	AppID string `json:"appID"`

	// The main title/heading for the toast notification.
	Title string `json:"title"`

	// The single/multi line message to display for the toast notification.
	Message string `json:"message,omitempty"`

	// A message attribution string
	Attribution string `json:"attribution,omitempty"`

	// An optional path to an image on the OS to display to the left of the title & message.
	Icon string `json:"icon,omitempty"`

	// An optional hero image, local image path...
	HeroImage string `json:"hero,omitempty"`

	// An optional image URL...
	Image string `json:"image,omitempty"`

	// The type of notification level action (like toast.Action)
	ActivationType string `json:"activationType,omitempty"`

	// The activation/action arguments (invoked when the user clicks the notification)
	ActivationArguments string `json:"activationArguments,omitempty"`

	// Optional action buttons to display below the notification title & message.
	Actions []Action `json:"actions,omitempty"`

	// The audio to play when displaying the toast
	Audio toastAudio `json:"audio,omitempty"`

	// Whether to loop the audio (default false)
	Loop bool `json:"loop,omitempty"`

	// How long the toast should show up for (short/long)
	Duration toastDuration `json:"duration,omitempty"`
}

// Action
//
// Defines an actionable button.
// See https://msdn.microsoft.com/en-us/windows/uwp/controls-and-patterns/tiles-and-notifications-adaptive-interactive-toasts for more info.
//
// Only protocol type action buttons are actually useful, as there's no way of receiving feedback from the
// user's choice. Examples of protocol type action buttons include: "bingmaps:?q=sushi" to open up Windows 10's
// maps app with a pre-populated search field set to "sushi".
//
//     toast.Action{"protocol", "Open Maps", "bingmaps:?q=sushi"}
type Action struct {
	Type      string `json:"type,omitempty"`
	Label     string `json:"label,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// Returns a toastAudio given a user-provided input (useful for cli apps).
//
// If the "name" doesn't match, then the default toastAudio is returned, along with ErrorInvalidAudio.
//
// The following names are valid;
//   - default
//   - im
//   - mail
//   - reminder
//   - sms
//   - loopingalarm
//   - loopimgalarm[2-10]
//   - loopingcall
//   - loopingcall[2-10]
//   - silent
//
// Handle the error appropriately according to how your app should work.
func Audio(name string) (toastAudio, error) {
	switch strings.ToLower(name) {
	case "default":
		return Default, nil
	case "im":
		return IM, nil
	case "mail":
		return Mail, nil
	case "reminder":
		return Reminder, nil
	case "sms":
		return SMS, nil
	case "loopingalarm":
		return LoopingAlarm, nil
	case "loopingalarm2":
		return LoopingAlarm2, nil
	case "loopingalarm3":
		return LoopingAlarm3, nil
	case "loopingalarm4":
		return LoopingAlarm4, nil
	case "loopingalarm5":
		return LoopingAlarm5, nil
	case "loopingalarm6":
		return LoopingAlarm6, nil
	case "loopingalarm7":
		return LoopingAlarm7, nil
	case "loopingalarm8":
		return LoopingAlarm8, nil
	case "loopingalarm9":
		return LoopingAlarm9, nil
	case "loopingalarm10":
		return LoopingAlarm10, nil
	case "loopingcall":
		return LoopingCall, nil
	case "loopingcall2":
		return LoopingCall2, nil
	case "loopingcall3":
		return LoopingCall3, nil
	case "loopingcall4":
		return LoopingCall4, nil
	case "loopingcall5":
		return LoopingCall5, nil
	case "loopingcall6":
		return LoopingCall6, nil
	case "loopingcall7":
		return LoopingCall7, nil
	case "loopingcall8":
		return LoopingCall8, nil
	case "loopingcall9":
		return LoopingCall9, nil
	case "loopingcall10":
		return LoopingCall10, nil
	case "silent":
		return Silent, nil
	default:
		return Default, ErrorInvalidAudio
	}
}

// Returns a toastDuration given a user-provided input (useful for cli apps).
//
// The default duration is short. If the "name" doesn't match, then the default toastDuration is returned,
// along with ErrorInvalidDuration. Most of the time "short" is the most appropriate for a toast notification,
// and Microsoft recommend not using "long", but it can be useful for important dialogs or looping sound toasts.
//
// The following names are valid;
//   - short
//   - long
//
// Handle the error appropriately according to how your app should work.
func Duration(name string) (toastDuration, error) {
	switch strings.ToLower(name) {
	case "short":
		return Short, nil
	case "long":
		return Long, nil
	default:
		return Short, ErrorInvalidDuration
	}
}
