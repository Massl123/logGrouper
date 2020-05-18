package loggrouper

// Profile contains common profiles for logGrouper
type Profile struct {
	Name       string
	LogFormat  string
	TimeFormat string
}

// Profiles contains all known profiles
var Profiles = map[string]Profile{"apacheAccessLog": ApacheAccessLog}

// ApacheAccessLog is the profile for apache access logs (common and combined are tested)
var ApacheAccessLog = Profile{Name: "apacheAccessLog", LogFormat: `^.*\[(?P<timestamp>.+)\].*"(?P<group>.+ [/\*].*?[/?\ ]).*".*$`, TimeFormat: "2/Jan/2006:15:04:05 -0700"}
