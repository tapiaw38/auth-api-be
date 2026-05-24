package domain

import (
	"errors"
	"regexp"
	"strings"
)

type (
	Optional[T any] struct {
		Set   bool
		Value *T
	}

	UserPatch struct {
		FirstName     Optional[string]
		LastName      Optional[string]
		Email         Optional[string]
		Picture       Optional[string]
		PhoneNumber   Optional[string]
		Address       Optional[string]
		IsActive      Optional[bool]
		VerifiedEmail Optional[bool]
	}
)

func (p UserPatch) HasChanges() bool {
	return p.FirstName.Set ||
		p.LastName.Set ||
		p.Email.Set ||
		p.Picture.Set ||
		p.PhoneNumber.Set ||
		p.Address.Set ||
		p.IsActive.Set ||
		p.VerifiedEmail.Set
}

func (u *User) ApplyPatch(p UserPatch) error {
	if !p.HasChanges() {
		return errors.New("at least one field is required")
	}

	if err := validateRequiredStringField("first_name", p.FirstName); err != nil {
		return err
	}
	if err := validateRequiredStringField("last_name", p.LastName); err != nil {
		return err
	}
	if p.Email.Set {
		if p.Email.Value == nil {
			return errors.New("email cannot be null")
		}
		if err := validateEmail(strings.TrimSpace(*p.Email.Value)); err != nil {
			return err
		}
	}
	if err := validateNullableStringField("picture", p.Picture); err != nil {
		return err
	}
	if err := validateNullableStringField("phone_number", p.PhoneNumber); err != nil {
		return err
	}
	if err := validateNullableStringField("address", p.Address); err != nil {
		return err
	}
	if err := validateBoolField("is_active", p.IsActive); err != nil {
		return err
	}
	if err := validateBoolField("verified_email", p.VerifiedEmail); err != nil {
		return err
	}

	if p.FirstName.Set {
		u.FirstName = strings.TrimSpace(*p.FirstName.Value)
	}
	if p.LastName.Set {
		u.LastName = strings.TrimSpace(*p.LastName.Value)
	}
	if p.Email.Set {
		u.Email = strings.TrimSpace(*p.Email.Value)
	}
	if p.Picture.Set {
		u.Picture = normalizeOptionalString(p.Picture.Value)
	}
	if p.PhoneNumber.Set {
		u.PhoneNumber = normalizeOptionalString(p.PhoneNumber.Value)
	}
	if p.Address.Set {
		u.Address = normalizeOptionalString(p.Address.Value)
	}
	if p.IsActive.Set {
		u.IsActive = *p.IsActive.Value
	}
	if p.VerifiedEmail.Set {
		u.VerifiedEmail = *p.VerifiedEmail.Value
	}

	return nil
}

func validateRequiredStringField(fieldName string, field Optional[string]) error {
	if !field.Set {
		return nil
	}
	if field.Value == nil {
		return errors.New(fieldName + " cannot be null")
	}

	value := strings.TrimSpace(*field.Value)
	if value == "" {
		return errors.New(fieldName + " cannot be empty")
	}
	if len(value) > 255 {
		return errors.New(fieldName + " must not exceed 255 characters")
	}

	return nil
}

func validateNullableStringField(fieldName string, field Optional[string]) error {
	if !field.Set || field.Value == nil {
		return nil
	}

	value := strings.TrimSpace(*field.Value)
	if value == "" {
		return errors.New(fieldName + " cannot be empty; use null to clear it")
	}
	if len(value) > 255 {
		return errors.New(fieldName + " must not exceed 255 characters")
	}

	return nil
}

func validateBoolField(fieldName string, field Optional[bool]) error {
	if !field.Set {
		return nil
	}
	if field.Value == nil {
		return errors.New(fieldName + " cannot be null")
	}

	return nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	if len(email) > 254 {
		return errors.New("email must not exceed 254 characters")
	}

	return nil
}
