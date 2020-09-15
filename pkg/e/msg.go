package e

var MsgFlags = map[int]string{
	SUCCESS:      "ok",
	CREATED:      "Created",
	BAD_REQUEST:  "Bad request",
	UNAUTHORIZED: "Unauthorized",
	FORBIDDEN:    "Forbidden",
	NOT_FOUND:    "Not Found",
	ERROR:        "Internal server error",

	MARSHAL_ERROR:  "Internal server error 5001",
	FAILED_TO_BIND: "Bad request 5002",
	USERNAME_TAKEN: "Username taken",
	FAILED_ATOI:    "Bad request 5003",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
