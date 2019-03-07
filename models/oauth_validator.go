package models

type oauthValidator struct {
	OAuthDB
}

func (ov *oauthValidator) Create(oauth *OAuth) error {
	err := oauthValidationFuncs(oauth,
		ov.userIDRequired,
		ov.serviceRequired)

	if err != nil {
		return err
	}

	return ov.OAuthDB.Create(oauth)
}

func (ov *oauthValidator) Delete(id uint) error {
	var oauth OAuth
	oauth.ID = id

	err := oauthValidationFuncs(&oauth, ov.validId)

	if err != nil {
		return err
	}

	return ov.OAuthDB.Delete(oauth.ID)
}

// Private

func (ov *oauthValidator) userIDRequired(oa *OAuth) error {
	if oa.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (ov *oauthValidator) serviceRequired(oa *OAuth) error {
	if oa.Service == "" {
		return ErrServiceRequired
	}
	return nil
}

func (ov *oauthValidator) validId(oa *OAuth) error {
	if oa.ID <= 0 {
		return ErrInvalidID
	}
	return nil
}

///////////////////////////////////////////////////////////
// Reusable validation functions helper
///////////////////////////////////////////////////////////

type oauthValidationFunc func(*OAuth) error

func oauthValidationFuncs(oauth *OAuth, funcs ...oauthValidationFunc) error {
	for _, fn := range funcs {
		err := fn(oauth)
		if err != nil {
			return err
		}
	}
	return nil
}
