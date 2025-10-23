package mappings

import "net/http"

var (
	RoleCreateNameRequiredError = ErrorDetails{
		"role:create:name-required",
		http.StatusBadRequest,
		"role name is required",
	}

	RoleCreateInvalidNameError = ErrorDetails{
		"role:create:invalid-name",
		http.StatusBadRequest,
		"invalid role name",
	}

	RoleCreateQueryError = ErrorDetails{
		"role:create:query-error",
		http.StatusInternalServerError,
		"failed to create role",
	}

	RoleCreateAlreadyExistsError = ErrorDetails{
		"role:create:already-exists",
		http.StatusConflict,
		"role with this name already exists",
	}

	RoleGetNotFoundError = ErrorDetails{
		"role:get:not-found",
		http.StatusNotFound,
		"role not found",
	}

	RoleGetQueryError = ErrorDetails{
		"role:get:query-error",
		http.StatusInternalServerError,
		"failed to retrieve role information",
	}

	RoleUpdateIDRequiredError = ErrorDetails{
		"role:update:id-required",
		http.StatusBadRequest,
		"role ID is required",
	}

	RoleUpdateNameRequiredError = ErrorDetails{
		"role:update:name-required",
		http.StatusBadRequest,
		"role name is required",
	}

	RoleUpdateNotFoundError = ErrorDetails{
		"role:update:not-found",
		http.StatusNotFound,
		"role not found",
	}

	RoleUpdateQueryError = ErrorDetails{
		"role:update:query-error",
		http.StatusInternalServerError,
		"failed to update role",
	}

	RoleUpdateInvalidNameError = ErrorDetails{
		"role:update:invalid-name",
		http.StatusBadRequest,
		"invalid role name",
	}

	RoleDeleteIDRequiredError = ErrorDetails{
		"role:delete:id-required",
		http.StatusBadRequest,
		"role ID is required",
	}

	RoleDeleteNotFoundError = ErrorDetails{
		"role:delete:not-found",
		http.StatusNotFound,
		"role not found",
	}

	RoleDeleteQueryError = ErrorDetails{
		"role:delete:query-error",
		http.StatusInternalServerError,
		"failed to delete role",
	}

	RoleDeleteInUseError = ErrorDetails{
		"role:delete:in-use",
		http.StatusConflict,
		"cannot delete role that is assigned to users",
	}

	RoleListQueryError = ErrorDetails{
		"role:list:query-error",
		http.StatusInternalServerError,
		"failed to retrieve role list",
	}

	RoleEnsureDefaultRoleNotFoundError = ErrorDetails{
		"role:ensure:default-role-not-found",
		http.StatusInternalServerError,
		"default user role not found in the system",
	}
)
