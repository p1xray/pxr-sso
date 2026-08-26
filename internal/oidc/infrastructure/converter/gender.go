package converter

import (
	"github.com/guregu/null/v6"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/enum"
)

// GenderEnumToNullInt16 converts pointer of GenderEnum to null.Int16 type.
func GenderEnumToNullInt16(ge *enum.Gender) null.Int16 {
	if ge == nil {
		return null.NewInt16(0, false)
	}
	return null.Int16From(int16(*ge))
}

// GenderEnumFromNullInt16 converts null.Int16 type to pointer of GenderEnum.
func GenderEnumFromNullInt16(value null.Int16) *enum.Gender {
	genderNumeric := value.Ptr()
	var gender *enum.Gender
	if genderNumeric != nil {
		genderValue := enum.Gender(*genderNumeric)
		gender = &genderValue
	}

	return gender
}
