package mistral

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
	Group              *string `json:"group"`
	IsBlocking         bool    `json:"is_blocking"`
}

type ModelCard struct {
	ID         string            `json:"id"`
	Object     string            `json:"object"`
	Created    int64             `json:"created"`
	OwnedBy    string            `json:"owned_by"`
	Root       *string           `json:"root"`
	Parent     *string           `json:"parent"`
	Permission []ModelPermission `json:"permission"`
}

type ModelList struct {
	Object string      `json:"object"`
	Data   []ModelCard `json:"data"`
}

func (mc *MistralClient) ListModels() (resp *ModelList, err error) {
	return resp, err
}
