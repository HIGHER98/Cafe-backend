package e

const (
	SUCCESS      = 200
	CREATED      = 201
	BAD_REQUEST  = 400
	UNAUTHORIZED = 401
	FORBIDDEN    = 403
	NOT_FOUND    = 404
	ERROR        = 500

	USERNAME_TAKEN              = 5000
	MARSHAL_ERROR               = 5001
	FAILED_TO_BIND              = 5002
	FAILED_ATOI                 = 5003
	ID_NOT_FOUND                = 5004
	DELETED                     = 5005
	INVALID_ROLE                = 5006
	STRIPE_CREATE_SESSION_ERROR = 5007
)
