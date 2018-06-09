package models

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := galleryValidationFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired)

	if err != nil {
		return err
	}

	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := galleryValidationFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired)

	if err != nil {
		return err
	}

	return gv.GalleryDB.Update(gallery)
}

///////////////////////////////////////////////////////////
// Private functions
///////////////////////////////////////////////////////////

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserId <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

///////////////////////////////////////////////////////////
// Reusable validation functions helper
///////////////////////////////////////////////////////////

type galleryValidationFunc func(*Gallery) error

func galleryValidationFuncs(gallery *Gallery, funcs ...galleryValidationFunc) error {
	for _, fn := range funcs {
		err := fn(gallery)
		if err != nil {
			return err
		}
	}
	return nil
}
