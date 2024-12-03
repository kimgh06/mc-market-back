package surge

import (
	"net/http"
	"net/url"
)

type API struct {
	Client        *http.Client
	Configuration APIConfiguration
}

type APIConfiguration struct {
	SurgeURL       string
	ParsedSurgeURL *url.URL
}

func NewAPIConfiguration(surgeURL string) APIConfiguration {
	parsed, err := url.Parse(surgeURL)
	if err != nil {
		panic(err)
	}

	return APIConfiguration{
		SurgeURL:       surgeURL,
		ParsedSurgeURL: parsed,
	}
}
