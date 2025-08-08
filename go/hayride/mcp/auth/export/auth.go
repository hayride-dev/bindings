package export

import (
	"unsafe"

	"github.com/hayride-dev/bindings/go/hayride/mcp/auth"
	witAuth "github.com/hayride-dev/bindings/go/internal/gen/exports/hayride/mcp/auth"
	"go.bytecodealliance.org/cm"
)

type Constructor func() (auth.Provider, error)

var providerConstructor Constructor

type resources struct {
	providers map[cm.Rep]auth.Provider
	errors    map[cm.Rep]errorResource
}

var resourceTable = &resources{
	providers: make(map[cm.Rep]auth.Provider),
	errors:    make(map[cm.Rep]errorResource),
}

func Provider(c Constructor) {
	providerConstructor = c

	witAuth.Exports.Provider.Constructor = constructor
	witAuth.Exports.Provider.AuthURL = authURL
	witAuth.Exports.Provider.Registration = registration
	witAuth.Exports.Provider.ExchangeCode = exchangeCode
	witAuth.Exports.Provider.Validate = validate
	witAuth.Exports.Provider.Destructor = destructor

	witAuth.Exports.Error.Code = errorCode
	witAuth.Exports.Error.Data = errorData
	witAuth.Exports.Error.Destructor = errorDestructor
}

func constructor() witAuth.Provider {
	auth, err := providerConstructor()
	if err != nil {
		return cm.ResourceNone
	}

	key := cm.Rep(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&auth))))
	v := witAuth.ProviderResourceNew(key)
	resourceTable.providers[key] = auth
	return v
}

func authURL(self cm.Rep) cm.Result[string, string, witAuth.Error] {
	provider, ok := resourceTable.providers[self]
	if !ok {
		wasiErr := createError(witAuth.ErrorCodeAuthURLFailed, "failed to find provider resource")
		return cm.Err[cm.Result[string, string, witAuth.Error]](wasiErr)
	}

	result, err := provider.AuthURL()
	if err != nil {
		wasiErr := createError(witAuth.ErrorCodeAuthURLFailed, err.Error())
		return cm.Err[cm.Result[string, string, witAuth.Error]](wasiErr)
	}

	return cm.OK[cm.Result[string, string, witAuth.Error]](result)
}

func registration(self cm.Rep, data cm.List[uint8]) cm.Result[cm.List[uint8], cm.List[uint8], witAuth.Error] {
	provider, ok := resourceTable.providers[self]
	if !ok {
		wasiErr := createError(witAuth.ErrorCodeRegistrationFailed, "failed to find provider resource")
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witAuth.Error]](wasiErr)
	}

	result, err := provider.Registration(data.Slice())
	if err != nil {
		wasiErr := createError(witAuth.ErrorCodeRegistrationFailed, err.Error())
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witAuth.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[uint8], cm.List[uint8], witAuth.Error]](cm.ToList(result))
}

func exchangeCode(self cm.Rep, data cm.List[uint8]) cm.Result[cm.List[uint8], cm.List[uint8], witAuth.Error] {
	provider, ok := resourceTable.providers[self]
	if !ok {
		wasiErr := createError(witAuth.ErrorCodeExchangeCodeFailed, "failed to find provider resource")
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witAuth.Error]](wasiErr)
	}

	result, err := provider.ExchangeCode(data.Slice())
	if err != nil {
		wasiErr := createError(witAuth.ErrorCodeExchangeCodeFailed, err.Error())
		return cm.Err[cm.Result[cm.List[uint8], cm.List[uint8], witAuth.Error]](wasiErr)
	}

	return cm.OK[cm.Result[cm.List[uint8], cm.List[uint8], witAuth.Error]](cm.ToList(result))
}

func validate(self cm.Rep, token string) cm.Result[witAuth.Error, bool, witAuth.Error] {
	provider, ok := resourceTable.providers[self]
	if !ok {
		wasiErr := createError(witAuth.ErrorCodeValidateFailed, "failed to find provider resource")
		return cm.Err[cm.Result[witAuth.Error, bool, witAuth.Error]](wasiErr)
	}

	result, err := provider.Validate(token)
	if err != nil {
		wasiErr := createError(witAuth.ErrorCodeValidateFailed, err.Error())
		return cm.Err[cm.Result[witAuth.Error, bool, witAuth.Error]](wasiErr)
	}

	return cm.OK[cm.Result[witAuth.Error, bool, witAuth.Error]](result)
}

func destructor(self cm.Rep) {
	delete(resourceTable.providers, self)
}
