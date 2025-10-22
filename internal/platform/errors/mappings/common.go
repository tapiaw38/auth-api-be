package mappings

import "net/http"

var (
	// Common validation errors
	RequestBodyParsingError = ErrorDetails{
		"common:request-body-parsing-error",
		http.StatusBadRequest,
		"invalid request body format",
	}

	InvalidParamsError = ErrorDetails{
		"common:invalid-params",
		http.StatusBadRequest,
		"invalid request parameters",
	}

	// Database errors
	DatabaseConnectionError = ErrorDetails{
		"common:database-connection-error",
		http.StatusInternalServerError,
		"database connection error",
	}

	DatabaseQueryError = ErrorDetails{
		"common:database-query-error",
		http.StatusInternalServerError,
		"database query error",
	}

	// External service errors
	ExternalServiceError = ErrorDetails{
		"common:external-service-error",
		http.StatusInternalServerError,
		"external service unavailable",
	}

	// Generic errors
	InternalServerError = ErrorDetails{
		"common:internal-server-error",
		http.StatusInternalServerError,
		"an internal server error occurred",
	}

	NotFoundError = ErrorDetails{
		"common:not-found",
		http.StatusNotFound,
		"resource not found",
	}

	ConflictError = ErrorDetails{
		"common:conflict",
		http.StatusConflict,
		"resource conflict",
	}
)
