package mistral

import (
	"context"
	"net/http"
)

// ModelPermission represents the permissions for a Mistral model.
type ModelPermission struct {
	ID                 string  `json:"id"`
	Object             string  `json:"object"`
	Created            int64   `json:"created"`
	AllowCreateEngine  bool    `json:"allow_create_engine"`
	AllowSampling      bool    `json:"allow_sampling"`
	AllowLogprobs      bool    `json:"allow_logprobs"`
	AllowSearchIndices bool    `json:"allow_search_indices"`
	AllowView          bool    `json:"allow_view"`
	AllowFineTuning    bool    `json:"allow_fine_tuning"`
	Organization       string  `json:"organization"`
	Group              *string `json:"group,omitempty"`
	IsBlocking         bool    `json:"is_blocking"`
}

// ModelCard represents a Mistral model.
type ModelCard struct {
	ID         string            `json:"id"`
	Object     string            `json:"object"`
	Created    int64             `json:"created"`
	OwnedBy    string            `json:"owned_by"`
	Root       *string           `json:"root"`
	Parent     *string           `json:"parent"`
	Permission []ModelPermission `json:"permission"`
}

// ModelList represents the list of available Mistral models.
type ModelList struct {
	Object string      `json:"object"`
	Data   []ModelCard `json:"data"`
}

func (mc *MistralClient) ListModels(ctx context.Context) (resp *ModelList, err error) {
	req, err := mc.newRequest(ctx, http.MethodGet, mc.endpoint("/models"), nil)
	if err != nil {
		return
	}

	err = mc.sendRequest(req, &resp)
	return
}
