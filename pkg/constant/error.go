package constant

const (
	SuccessMess           = "Success"
	ConflictErrMess       = "Conflict Error"
	BadRequestErrMess     = "Bad Request Error"
	UnAuthorizedErrMess   = "UnAuthorized Error"
	NotFoundErrMess       = "Not Found Error"
	ForbiddenErrMess      = "Forbidden Error"
	InternalServerErrMess = "Internal Server Error"
)

const (
	// General Errors: 0 -> -49
	// Processing indicate success but the object is being processed
	Processing = 2
	// Success indicates no error
	Success = 1
	// Unknown error indicates unknown state or step
	Unknown = 0
	// BadRequest error
	BadRequestErr = -1
	// NotFound error
	NotFound = -2
	// AuthenFailed error
	AuthenticationFailed = -3
	// Internal server error
	InternalServerError = -4
	// IllegalStateError
	IllegalStateError = -5
	// SendMessageError
	SendMessageError = -6
	// Call Internal API Error
	CallInternalAPIError = -7
	// Invalid Data
	InvalidData = -8
	// SerializeError
	SerializingError = -9
	// DeserializeError
	DeserializingError = -10
	// CastingError
	CastingError = -11
	// ParsingError
	ParsingError = -12
	// ConflictError
	ConflictError = -13
	// ForbiddenError
	ForbiddenError = -14
)
