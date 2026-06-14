package state

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tfdriftctl/tfdriftctl/internal/model"
)

type tfcStateVersion struct {
	Data struct {
		Attributes struct {
			HostedStateDownloadURL string `json:"hosted-state-download-url"`
		} `json:"attributes"`
	} `json:"data"`
}

func readTFC(ctx context.Context, cfg model.StateConfig) ([]byte, error) {
	if cfg.WorkspaceID == "" || cfg.Token == "" {
		return nil, fmt.Errorf("tfc backend requires workspace_id and token")
	}

	url := fmt.Sprintf("https://app.terraform.io/api/v2/workspaces/%s/current-state-version", cfg.WorkspaceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tfc api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tfc api returned status %d", resp.StatusCode)
	}

	var stateVer tfcStateVersion
	if err := json.NewDecoder(resp.Body).Decode(&stateVer); err != nil {
		return nil, fmt.Errorf("decode tfc response: %w", err)
	}

	downloadURL := stateVer.Data.Attributes.HostedStateDownloadURL
	if downloadURL == "" {
		return nil, fmt.Errorf("tfc hosted state download url is empty")
	}

	// Fetch the actual state payload
	stateReq, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}

	stateResp, err := http.DefaultClient.Do(stateReq)
	if err != nil {
		return nil, fmt.Errorf("download tfc state failed: %w", err)
	}
	defer stateResp.Body.Close()

	if stateResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download tfc state returned status %d", stateResp.StatusCode)
	}

	data, err := io.ReadAll(stateResp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
