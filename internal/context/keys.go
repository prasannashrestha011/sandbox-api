package request_context

type (
	userIDKeyType   struct{}
	userRoleKeyType struct{}
	userTypeKeyType struct{}
)

var (
	userIDKey   = userIDKeyType{}
	userRoleKey = userRoleKeyType{}
	userTypeKey = userTypeKeyType{}
)
