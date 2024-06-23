package auth

import (
	"AABBCCDD/app/db"
	"fmt"

	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
)

var profileSchema = v.Schema{
	"firstName": v.Rules(v.Min(3), v.Max(50)),
	"lastName":  v.Rules(v.Min(3), v.Max(50)),
}

type ProfileFormValues struct {
	ID        uint   `form:"id"`
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
	Email     string
	Success   string
}

func HandleProfileShow(kit *kit.Kit) error {
	auth := kit.Auth().(Auth)

	var user User
	if err := db.Get().First(&user, auth.UserID).Error; err != nil {
		return err
	}

	formValues := ProfileFormValues{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return kit.Render(ProfileShow(formValues))
}

func HandleProfileUpdate(kit *kit.Kit) error {
	var values ProfileFormValues
	errors, ok := v.Request(kit.Request, &values, profileSchema)
	if !ok {
		return kit.Render(ProfileForm(values, errors))
	}

	auth := kit.Auth().(Auth)
	if auth.UserID != values.ID {
		return fmt.Errorf("unauthorized request for profile %d", values.ID)
	}
	err := db.Get().Model(&User{}).
		Where("id = ?", auth.UserID).
		Updates(&User{
			FirstName: values.FirstName,
			LastName:  values.LastName,
		}).Error
	if err != nil {
		return err
	}

	values.Success = "Profile successfully updated!"
	values.Email = auth.Email

	return kit.Render(ProfileForm(values, v.Errors{}))
}
