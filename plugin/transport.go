package plugin

import (
	"github.com/oursky/skygear/skyconfig"
	"github.com/oursky/skygear/skydb"
	"golang.org/x/net/context"
)

// AuthRequest is sent by Skygear to plugin which contains data for authentication
type AuthRequest struct {
	ProviderName string
	Action       string
	AuthData     map[string]interface{}
}

// AuthResponse is sent by plugin to Skygear which contains authenticated data
type AuthResponse struct {
	PrincipalID string                 `json:"principal_id"`
	AuthData    map[string]interface{} `json:"auth_data"`
}

// TransportState refers to the operation state of the transport
//go:generate stringer -type=TransportState
type TransportState int

const (
	// TransportStateUninitialized is the state when the transport has not
	// been initialized
	TransportStateUninitialized TransportState = iota

	// TransportStateReady is the state when the transport is ready for
	// requests
	TransportStateReady

	// TransportStateWorkerUnavailable is the state when all workers
	// for the transport is not available
	TransportStateWorkerUnavailable

	// TransportStateError is the state when an error has occurred
	// in the transport and it is not able to serve requests
	TransportStateError
)

type TransportInitHandler func([]byte, error) error

// A Transport represents the interface of data transfer between skygear
// and remote process.
type Transport interface {
	State() TransportState
	SetInitHandler(TransportInitHandler)
	RequestInit()
	RunInit() ([]byte, error)
	RunLambda(ctx context.Context, name string, in []byte) ([]byte, error)
	RunHandler(ctx context.Context, name string, in []byte) ([]byte, error)

	// RunHook runs the hook with a name recognized by plugin, passing in
	// record as a parameter. Transport may not modify the record passed in.
	//
	// A skydb.Record is returned as a result of invocation. Such record must be
	// a newly allocated instance, and may not share any reference type values
	// in any of its memebers with the record being passed in.
	RunHook(ctx context.Context, hookName string, record *skydb.Record, oldRecord *skydb.Record) (*skydb.Record, error)

	RunTimer(name string, in []byte) ([]byte, error)

	// RunProvider runs the auth provider with the specified AuthRequest.
	RunProvider(request *AuthRequest) (*AuthResponse, error)
}

// A TransportFactory is a generic interface to instantiates different
// kinds of Plugin Transport.
type TransportFactory interface {
	Open(path string, args []string, config skyconfig.Configuration) Transport
}

// ContextMap returns a map of the user request context.
func ContextMap(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return map[string]interface{}{}
	}
	pluginCtx := map[string]interface{}{}
	if userID, ok := ctx.Value("UserID").(string); ok {
		pluginCtx["user_id"] = userID
	}
	return pluginCtx
}
