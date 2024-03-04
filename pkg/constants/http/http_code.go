package http

// HTTP response code
const (
	HTTPOK        = 200
	HTTPCreated   = 201
	HTTPNoContent = 204

	HTTPBadRequest          = 400
	HTTPUnauthorized        = 401
	HTTPForbidden           = 403
	HTTPNotFound            = 404
	HTTPConflict            = 409
	HTTPMethodNotAllowed    = 405
	HTTPUnprocessableEntity = 422
	HTTPTooManyRequests     = 429

	HTTPInternalServerError = 500
	HTTPUnimplemented       = 501
	HTTPAccepted            = 202
)
