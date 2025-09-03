package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/mcp"
	"github.com/hayride-dev/bindings/go/hayride/mcp/resources"
	witResources "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/mcp/resources"
	"go.bytecodealliance.org/cm"
)

type Constructor func() (resources.Resources, error)

var resourcesConstructor Constructor

type wasmResources struct {
	resources map[cm.Rep]resources.Resources
	errors    map[cm.Rep]errorResource
}

var resourceTable = &wasmResources{
	resources: make(map[cm.Rep]resources.Resources),
	errors:    make(map[cm.Rep]errorResource),
}

func Resources(c Constructor) {
	resourcesConstructor = c

	witResources.Exports.Resources.Constructor = constructor
	witResources.Exports.Resources.ReadResources = read
	witResources.Exports.Resources.ListResources = list
	witResources.Exports.Resources.ListTemplates = listTemplates
	witResources.Exports.Resources.Destructor = destructor

	witResources.Exports.Error.Code = errorCode
	witResources.Exports.Error.Data = errorData
	witResources.Exports.Error.Destructor = errorDestructor
}

func constructor() witResources.Resources {
	res, err := resourcesConstructor()
	if err != nil {
		return cm.ResourceNone
	}

	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&res))))
	v := witResources.ResourcesResourceNew(key)
	resourceTable.resources[key] = res
	return v
}

func read(self cm.Rep, params witResources.ReadResourceParams) cm.Result[witResources.ReadResourceResultShape, witResources.ReadResourceResult, witResources.Error] {
	resource, ok := resourceTable.resources[self]
	if !ok {
		wasiErr := createError(witResources.ErrorCodeResourceNotFound, "failed to find resource")
		return cm.Err[cm.Result[witResources.ReadResourceResultShape, witResources.ReadResourceResult, witResources.Error]](wasiErr)
	}

	result, err := resource.Read(cm.Reinterpret[mcp.ReadResourceParams](params))
	if err != nil {
		wasiErr := createError(witResources.ErrorCodeResourceNotFound, err.Error())
		return cm.Err[cm.Result[witResources.ReadResourceResultShape, witResources.ReadResourceResult, witResources.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witResources.ReadResourceResultShape, witResources.ReadResourceResult, witResources.Error]](cm.Reinterpret[witResources.ReadResourceResult](*result))
}

func list(self cm.Rep, cursor string) cm.Result[witResources.ListResourcesResultShape, witResources.ListResourcesResult, witResources.Error] {
	resource, ok := resourceTable.resources[self]
	if !ok {
		wasiErr := createError(witResources.ErrorCodeResourceNotFound, "failed to find resource")
		return cm.Err[cm.Result[witResources.ListResourcesResultShape, witResources.ListResourcesResult, witResources.Error]](wasiErr)
	}

	result, err := resource.List(cursor)
	if err != nil {
		wasiErr := createError(witResources.ErrorCodeResourceNotFound, err.Error())
		return cm.Err[cm.Result[witResources.ListResourcesResultShape, witResources.ListResourcesResult, witResources.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witResources.ListResourcesResultShape, witResources.ListResourcesResult, witResources.Error]](cm.Reinterpret[witResources.ListResourcesResult](*result))
}

func listTemplates(self cm.Rep, cursor string) cm.Result[witResources.ListResourceTemplatesResultShape, witResources.ListResourceTemplatesResult, witResources.Error] {
	resource, ok := resourceTable.resources[self]
	if !ok {
		wasiErr := createError(witResources.ErrorCodeResourceNotFound, "failed to find resource")
		return cm.Err[cm.Result[witResources.ListResourceTemplatesResultShape, witResources.ListResourceTemplatesResult, witResources.Error]](wasiErr)
	}

	result, err := resource.ListTemplates(cursor)
	if err != nil {
		wasiErr := createError(witResources.ErrorCodeResourceNotFound, err.Error())
		return cm.Err[cm.Result[witResources.ListResourceTemplatesResultShape, witResources.ListResourceTemplatesResult, witResources.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witResources.ListResourceTemplatesResultShape, witResources.ListResourceTemplatesResult, witResources.Error]](cm.Reinterpret[witResources.ListResourceTemplatesResult](*result))
}

func destructor(self cm.Rep) {
	delete(resourceTable.resources, self)
}
