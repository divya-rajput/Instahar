package models

type HarEntry struct {
	Request  HarRequest
	Response HarResponse
}

type HarRequest struct {
	Method  string
	Url     string
	Headers []struct {
		Name  string
		Value string
	}
	QueryString []struct {
		Name  string
		Value string
	}
	Cookies []struct {
		Name  string
		Value string
	}
}

type HarResponse struct {
	Status  int
	Headers []struct {
		Name  string
		Value string
	}
	Cookies []struct {
		Name  string
		Value string
	}
	Content struct {
		MimeType string
		Size     uint64
		Text     string
	}
}

type HarFileData struct {
	Log struct {
		Entries []HarEntry
	}
}
