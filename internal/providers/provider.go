package providers

import (
	"context"

	"github.com/tfdriftctl/tfdriftctl/internal/model"
	"github.com/tfdriftctl/tfdriftctl/internal/providers/aws"
	"github.com/tfdriftctl/tfdriftctl/internal/providers/azure"
)

// CloudProvider fetches live cloud resources.
type CloudProvider interface {
	Name() string
	FetchResources(ctx context.Context, expected []model.Resource, workspace model.Workspace) ([]model.Resource, error)
	SupportedTypes() []string
}

// Registry holds registered cloud providers.
type Registry struct {
	providers map[string]CloudProvider
}

func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]CloudProvider)}
}

func (r *Registry) Register(p CloudProvider) {
	r.providers[p.Name()] = p
}

func (r *Registry) Get(name string) (CloudProvider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

func DefaultRegistry() *Registry {
	r := NewRegistry()
	r.Register(aws.NewProvider())
	r.Register(azure.NewProvider())
	return r
}
