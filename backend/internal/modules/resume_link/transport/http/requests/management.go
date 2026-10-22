package requests

type CreateResumeLinkRequest struct {
	ResumeID string `json:"resumeId"`
	Name     string `json:"name"`
}

type RenameResumeLinkRequest struct {
	Name    string `json:"name"`
	Version int64  `json:"version"`
}

type ReplaceResumeLinkStatusRequest struct {
	Status  string `json:"status"`
	Version int64  `json:"version"`
}
