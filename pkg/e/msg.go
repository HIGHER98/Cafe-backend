package e

var MsgFlags = map[int]string{
	SUCCESS:      "ok",
	CREATED:      "Created",
	BAD_REQUEST:  "Bad request",
	UNAUTHORIZED: "Unauthorized",
	FORBIDDEN:    "Forbidden",
	NOT_FOUND:    "Not Found",
	ERROR:        "Internal server error",

	USERNAME_TAKEN: "Username taken",
	MARSHAL_ERROR:  "Internal server error 5001",
	FAILED_TO_BIND: "Bad request 5002",
	FAILED_ATOI:    "Bad request 5003",
	ID_NOT_FOUND:   "Bad request 5004",
	DELETED:        "Deleted successfully",
	INVALID_ROLE:   "Invalid role applied",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
