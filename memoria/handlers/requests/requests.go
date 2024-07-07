package requests

type NewProcessRequest struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}
