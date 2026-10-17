package requests

type CreateResumeLinkRequest struct {
	ResumeID string `json:"resumeId"`
	Name     string `json:"name"`
}

type ReplaceResumeLinkStatusRequest struct {
	Status  string `json:"status"`
	Version int64  `json:"version"`
}
