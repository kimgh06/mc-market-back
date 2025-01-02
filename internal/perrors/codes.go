package perrors

const (
	UnauthorizedMissingBearer Code = "unauthorized_missing_bearer"
	UnauthorizedFailedParse   Code = "unauthorized_failed_parse"
	UnauthorizedFailedSubject Code = "unauthorized_failed_subject"

	InsufficientUserPermission Code = "insufficient_user_permission"

	FailedFormFile Code = "failed_form_file"
	FailedStorage  Code = "failed_storage"
	FailedDatabase Code = "failed_database"
	FailedAPI      Code = "failed_api"
	FailedValidate Code = "failed_validate"
	FailedMarshal  Code = "failed_marshal"

	Mismatching Code = "mismatching"

	InvalidQuery     Code = "invalid_query"
	InvalidJSON      Code = "invalid_json"
	InvalidParameter Code = "invalid_parameter"
	InvalidUUID      Code = "invalid_uuid"
	InvalidSnowflake Code = "invalid_snowflake"
	InvalidHTML      Code = "invalid_html"

	UserNotFound Code = "user_not_found"

	ProductNotFound Code = "maple_not_found"

	InsufficientFunds Code = "insufficient_funds"
	DuplicatePurchase Code = "duplicate_purchase"
	PurchaseNotFound
	PaymentNotFound Code = "payment_not_found"

	ArticleNotFound Code = "article_not_found"

	UnknownInternalError Code = "unknown_internal_error"
)
