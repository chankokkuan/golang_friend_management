package domain

type CustomErr struct {
	code    string
	message string
	details map[string]interface{}
}

func (e CustomErr) Code() string {
	return e.code
}

func (e CustomErr) Error() string {
	return e.message
}

func (e CustomErr) Details() map[string]interface{} {
	return e.details
}

func CustomizeError(code string, message string, details map[string]interface{}) *CustomErr {
	return &CustomErr{code, message, details}
}

const (
	NotFound        = "not-found"
	Unknown         = "unknown"
	VersionConflict = "version-conflict"
	AccessDenied    = "access-denied"
	IsFriend        = "is-friend"
)

var (
	ErrNotFound        = &CustomErr{NotFound, "No user found", nil}
	ErrUnknown         = &CustomErr{Unknown, "Unknown error occured. Please try again in a few minutes", nil}
	ErrVersionConflict = &CustomErr{VersionConflict, "The data version is not the latest version. Please re-enter / refresh the page and try again later", nil}
	ErrAccessDenied    = &CustomErr{AccessDenied, "Access denied", nil}
	ErrFriendConflict  = &CustomErr{IsFriend, "The users are already friend to each other", nil}
)
