package azure

import (
	"context"

	"github.com/tfdriftctl/tfdriftctl/internal/model"
)

// Provider implements Azure cloud resource fetching.
type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Name() string { return "azure" }

func (p *Provider) SupportedTypes() []string {
	return []string{
		"azurerm_resource_group",
		"azurerm_virtual_network",
		// Add more Azure resources here
	}
}

func (p *Provider) FetchResources(ctx context.Context, expected []model.Resource, workspace model.Workspace) ([]model.Resource, error) {
	// TODO: Implement Azure Resource Graph API calls or standard Azure SDK
	// This acts as a foundation demonstrating multi-cloud plugin architecture.
	return []model.Resource{}, nil
}
