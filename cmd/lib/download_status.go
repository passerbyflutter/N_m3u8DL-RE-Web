package lib

type DownloadStatus string

const (
	Pending     DownloadStatus = "Pending"
	Downloading DownloadStatus = "Downloading"
	Finished    DownloadStatus = "Finished"
	Deleted     DownloadStatus = "Deleted"
	Error       DownloadStatus = "Error"
)
