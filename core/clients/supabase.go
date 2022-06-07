package clients

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/websublime/sublime-cli/utils"
)

const (
	AuthEndpoint = "auth/v1"
	RestEndpoint = "rest/v1"
)

type Supabase struct {
	BaseURL string
	// apiKey can be a client API key or a service key
	apiKey     string
	HTTPClient *http.Client
}

func NewSupabase(baseURL string, supabaseKey string) *Supabase {
	return &Supabase{
		BaseURL: baseURL,
		apiKey:  supabaseKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (ctx *Supabase) Upload(directory string) {
	fileList, err := utils.PathWalk(directory)
	if err != nil {
		panic(err)
	}

	for idx := range fileList {
		file, err := os.Open(fileList[idx])
		if err != nil {
			panic(fmt.Errorf("os.Open: %v", err))
		}
		file.Close()
	}
}
