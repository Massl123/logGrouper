package loggrouper

// Profile contains common profiles for logGrouper
type Profile struct {
	Name       string
	LogFormat  string
	TimeFormat string
}

// Profiles contains all known profiles
var Profiles = map[string]Profile{profileApacheAccessLogFull.Name: profileApacheAccessLogFull, profileApacheAccessLogFirstGroup.Name: profileApacheAccessLogFirstGroup}

// profileApacheAccessLogFirstGroup is the profile for apache access logs, matching the first group in the URL
var profileApacheAccessLogFirstGroup = Profile{
	Name:       "apacheAccessLog-first-group",
	LogFormat:  `^.*?\[(?P<timestamp>.+?)\].*?"(?P<group>.+ [/\*].*?[/?\ ]).+?" .*$`,
	TimeFormat: "2/Jan/2006:15:04:05 -0700"}

// profileApacheAccessLogFull is the profile for apache access logs, matching the whole URL
var profileApacheAccessLogFull = Profile{
	Name:       "apacheAccessLog-full",
	LogFormat:  `^.*?\[(?P<timestamp>.+?)\].*?"(?P<group>.+? .+?) .+?" .*$`,
	TimeFormat: "2/Jan/2006:15:04:05 -0700"}
