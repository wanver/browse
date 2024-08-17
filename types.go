package browse

import (
	"github.com/go-rod/rod/lib/proto"
)

type BrowseRequest struct {
	PageURL        string                      `json:"page_url"`
	Cookies        []*proto.NetworkCookieParam `json:"cookies"`
	Instructions   []BrowseRequestInstruction  `json:"instructions"`
	Headless       bool                        `json:"headless"`
	Proxy          string                      `json:"proxy"`
	HijackRequests []string                    `json:"hijack_requests"`
}

type BrowseResponse struct {
	Hijacks map[string]any `json:"hijacks"`
	Error   string         `json:"error,omitempty"`
}

type BrowseRequestInstruction struct {
	Selector    string              `json:"selector"`
	Action      BrowseRequestAction `json:"action"`
	Input       string              `json:"input"`
	Frames      []string            `json:"wait_for_frame"`
	Fatal       bool                `json:"fatal"`
	WaitVisible bool                `json:"wait_visible"`
}

type BrowseRequestAction string

const (
	Click BrowseRequestAction = "click"
	Type  BrowseRequestAction = "type"
)

func (bra BrowseRequestAction) String() string {
	return string(bra)
}
