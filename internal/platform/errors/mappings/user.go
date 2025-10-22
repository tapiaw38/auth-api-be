package mappings

import "net/http"

var (
	// Login errors
	UserLoginInvalidCredentialsError = ErrorDetails{
		"user:login:invalid-credentials",
		http.StatusUnauthorized,
		"invalid email or password",
	}

	UserLoginUserNotFoundError = ErrorDetails{
		"user:login:user-not-found",
		http.StatusNotFound,
		"user not found",
	}

	UserLoginUserNotActiveError = ErrorDetails{
		"user:login:user-not-active",
		http.StatusForbidden,
		"user account is not active",
	}

	UserLoginSSOAuthMethodError = ErrorDetails{
		"user:login:sso-auth-method",
		http.StatusBadRequest,
		"this account uses SSO authentication. Please use Google login",
	}

	UserLoginNoPasswordSetError = ErrorDetails{
		"user:login:no-password-set",
		http.StatusBadRequest,
		"account has no password set",
	}

	UserLoginTokenGenerationError = ErrorDetails{
		"user:login:token-generation-error",
		http.StatusInternalServerError,
		"failed to generate authentication token",
	}

	UserLoginGoogleExchangeError = ErrorDetails{
		"user:login:google-exchange-error",
		http.StatusInternalServerError,
		"failed to exchange Google authorization code",
	}

	UserLoginGoogleUserInfoError = ErrorDetails{
		"user:login:google-user-info-error",
		http.StatusInternalServerError,
		"failed to retrieve user information from Google",
	}

	// Register errors
	UserRegisterInvalidInputError = ErrorDetails{
		"user:register:invalid-input",
		http.StatusBadRequest,
		"invalid registration data",
	}

	UserRegisterEmailRequiredError = ErrorDetails{
		"user:register:email-required",
		http.StatusBadRequest,
		"email is required",
	}

	UserRegisterPasswordRequiredError = ErrorDetails{
		"user:register:password-required",
		http.StatusBadRequest,
		"password is required",
	}

	UserRegisterCreateUserError = ErrorDetails{
		"user:register:create-user-error",
		http.StatusInternalServerError,
		"failed to create user account",
	}

	UserRegisterAssignRoleError = ErrorDetails{
		"user:register:assign-role-error",
		http.StatusInternalServerError,
		"failed to assign default role to user",
	}

	// Get user errors
	UserGetNotFoundError = ErrorDetails{
		"user:get:not-found",
		http.StatusNotFound,
		"user not found",
	}

	UserGetQueryError = ErrorDetails{
		"user:get:query-error",
		http.StatusInternalServerError,
		"failed to retrieve user information",
	}

	// Update user errors
	UserUpdateNotFoundError = ErrorDetails{
		"user:update:not-found",
		http.StatusNotFound,
		"user not found",
	}

	UserUpdateQueryError = ErrorDetails{
		"user:update:query-error",
		http.StatusInternalServerError,
		"failed to update user information",
	}

	// Delete user errors
	UserDeleteNotFoundError = ErrorDetails{
		"user:delete:not-found",
		http.StatusNotFound,
		"user not found",
	}

	UserDeleteQueryError = ErrorDetails{
		"user:delete:query-error",
		http.StatusInternalServerError,
		"failed to delete user",
	}

	// List users errors
	UserListQueryError = ErrorDetails{
		"user:list:query-error",
		http.StatusInternalServerError,
		"failed to retrieve user list",
	}

	// Password reset errors
	UserRequestResetPasswordEmailRequiredError = ErrorDetails{
		"user:request-reset-password:email-required",
		http.StatusBadRequest,
		"email is required",
	}

	UserRequestResetPasswordUserNotFoundError = ErrorDetails{
		"user:request-reset-password:user-not-found",
		http.StatusNotFound,
		"user not found",
	}

	UserRequestResetPasswordTokenGenerationError = ErrorDetails{
		"user:request-reset-password:token-generation-error",
		http.StatusInternalServerError,
		"failed to generate password reset token",
	}

	UserRequestResetPasswordUpdateError = ErrorDetails{
		"user:request-reset-password:update-error",
		http.StatusInternalServerError,
		"failed to save password reset token",
	}

	UserRequestResetPasswordEmailSendError = ErrorDetails{
		"user:request-reset-password:email-send-error",
		http.StatusInternalServerError,
		"failed to send password reset email",
	}

	UserResetPasswordTokenRequiredError = ErrorDetails{
		"user:reset-password:token-required",
		http.StatusBadRequest,
		"reset token is required",
	}

	UserResetPasswordNewPasswordRequiredError = ErrorDetails{
		"user:reset-password:new-password-required",
		http.StatusBadRequest,
		"new password is required",
	}

	UserResetPasswordInvalidTokenError = ErrorDetails{
		"user:reset-password:invalid-token",
		http.StatusBadRequest,
		"invalid or expired reset token",
	}

	UserResetPasswordUpdateError = ErrorDetails{
		"user:reset-password:update-error",
		http.StatusInternalServerError,
		"failed to reset password",
	}

	// Change password errors
	UserChangePasswordOldPasswordRequiredError = ErrorDetails{
		"user:change-password:old-password-required",
		http.StatusBadRequest,
		"old password is required",
	}

	UserChangePasswordNewPasswordRequiredError = ErrorDetails{
		"user:change-password:new-password-required",
		http.StatusBadRequest,
		"new password is required",
	}

	UserChangePasswordInvalidOldPasswordError = ErrorDetails{
		"user:change-password:invalid-old-password",
		http.StatusUnauthorized,
		"old password is incorrect",
	}

	UserChangePasswordUpdateError = ErrorDetails{
		"user:change-password:update-error",
		http.StatusInternalServerError,
		"failed to change password",
	}

	// Verify email errors
	UserVerifyEmailTokenRequiredError = ErrorDetails{
		"user:verify-email:token-required",
		http.StatusBadRequest,
		"verification token is required",
	}

	UserVerifyEmailInvalidTokenError = ErrorDetails{
		"user:verify-email:invalid-token",
		http.StatusBadRequest,
		"invalid or expired verification token",
	}

	UserVerifyEmailUpdateError = ErrorDetails{
		"user:verify-email:update-error",
		http.StatusInternalServerError,
		"failed to verify email",
	}

	// Set password errors
	UserSetPasswordRequiredError = ErrorDetails{
		"user:set-password:password-required",
		http.StatusBadRequest,
		"password is required",
	}

	UserSetPasswordUpdateError = ErrorDetails{
		"user:set-password:update-error",
		http.StatusInternalServerError,
		"failed to set password",
	}
)
