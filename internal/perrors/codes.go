package perrors

const (
	UnauthorizedMissingBearer Code = "unauthorized_missing_bearer"
	UnauthorizedFailedParse   Code = "unauthorized_failed_parse"
	UnauthorizedFailedSubject Code = "unauthorized_failed_subject"

	InsufficientUserPermission Code = "insufficient_user_permission"

	FailedDatabase Code = "failed_database"
	FailedAPI      Code = "failed_api"
	FailedValidate Code = "failed_validate"
	FailedMarshal  Code = "failed_marshal"

	InvalidJSON      Code = "invalid_json"
	InvalidParameter Code = "invalid_parameter"
	InvalidUUID      Code = "invalid_uuid"
	InvalidSnowflake Code = "invalid_snowflake"

	UserNotFound Code = "user_not_found"

	ProductNotFound Code = "maple_not_found"
)
