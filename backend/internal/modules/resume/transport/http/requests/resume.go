package requests

type CreateResumeRequest struct {
	Title    string `json:"title"`
	Language string `json:"language"`
}

type ReplaceResumeTitleRequest struct {
	Version int64  `json:"version"`
	Title   string `json:"title"`
}

type ChangeResumeStatusRequest struct {
	Status  string `json:"status"`
	Version int64  `json:"version"`
}

type DeleteResumeRequest struct {
	Version int64 `json:"version"`
}
