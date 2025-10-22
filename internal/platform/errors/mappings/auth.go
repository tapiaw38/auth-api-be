package mappings

import "net/http"

var (
	// Authentication errors
	AuthMissingTokenError = ErrorDetails{
		"auth:missing-token",
		http.StatusUnauthorized,
		"authentication token is required",
	}

	AuthInvalidTokenError = ErrorDetails{
		"auth:invalid-token",
		http.StatusUnauthorized,
		"invalid authentication token",
	}

	AuthExpiredTokenError = ErrorDetails{
		"auth:expired-token",
		http.StatusUnauthorized,
		"authentication token has expired",
	}

	AuthTokenVersionMismatchError = ErrorDetails{
		"auth:token-version-mismatch",
		http.StatusUnauthorized,
		"token version mismatch - please login again",
	}

	AuthUnauthorizedError = ErrorDetails{
		"auth:unauthorized",
		http.StatusUnauthorized,
		"you are not authorized to perform this action",
	}

	AuthForbiddenError = ErrorDetails{
		"auth:forbidden",
		http.StatusForbidden,
		"you do not have permission to access this resource",
	}

	// Password errors
	AuthPasswordHashError = ErrorDetails{
		"auth:password-hash-error",
		http.StatusInternalServerError,
		"failed to hash password",
	}

	AuthPasswordCompareError = ErrorDetails{
		"auth:password-compare-error",
		http.StatusUnauthorized,
		"incorrect password",
	}

	AuthPasswordTooWeakError = ErrorDetails{
		"auth:password-too-weak",
		http.StatusBadRequest,
		"password does not meet security requirements",
	}
)
