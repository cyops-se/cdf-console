package messages

type Command struct {
	Command string `json:"command"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

type CommandResponse struct {
	RequestResponse
	Commands []Command `json:"commands"`
}

type SystemInformation struct {
	GitVersion string `json:"gitversion"`
	GitCommit  string `json:"gitcommit"`
}

type SysInfoResponse struct {
	RequestResponse
	SysInfo SystemInformation `json:"sysinfo"`
}
