package resources

import (
	"fmt"

	"github.com/hayride-dev/bindings/go/hayride/types"
	witResources "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/mcp/resources"

	"go.bytecodealliance.org/cm"
)

var _ Resources = (*MCPResource)(nil)

type Resources interface {
	Read(params types.ReadResourceParams) (*types.ReadResourceResult, error)
	List(cursor string) (*types.ListResourcesResult, error)
	ListTemplates(cursor string) (*types.ListResourceTemplatesResult, error)
}

type MCPResource cm.Resource

func New() (MCPResource, error) {
	return MCPResource(witResources.NewResources()), nil
}

func (t MCPResource) Read(params types.ReadResourceParams) (*types.ReadResourceResult, error) {
	witMCPResource := cm.Reinterpret[witResources.Resources](t)

	result := witMCPResource.ReadResources(cm.Reinterpret[witResources.ReadResourceParams](params))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to read resources: %s", result.Err().Data())
	}

	return cm.Reinterpret[*types.ReadResourceResult](result.OK()), nil
}

func (t MCPResource) List(cursor string) (*types.ListResourcesResult, error) {
	witMCPResource := cm.Reinterpret[witResources.Resources](t)

	result := witMCPResource.ListResources(cursor)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get resources: %s", result.Err().Data())
	}

	return cm.Reinterpret[*types.ListResourcesResult](result.OK()), nil
}

func (t MCPResource) ListTemplates(cursor string) (*types.ListResourceTemplatesResult, error) {
	witMCPResource := cm.Reinterpret[witResources.Resources](t)

	result := witMCPResource.ListTemplates(cursor)
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get resource templates: %s", result.Err().Data())
	}

	return cm.Reinterpret[*types.ListResourceTemplatesResult](result.OK()), nil
}
