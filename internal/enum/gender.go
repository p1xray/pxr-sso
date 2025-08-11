package enum

import "github.com/guregu/null/v6"

// GenderEnum is type for gender enum.
type GenderEnum int16

// Gender enum.
const (
	MALE   GenderEnum = 1
	FEMALE GenderEnum = 2
)

// ToNullInt16 converts GenderEnum to nullable int16 type.
func (ge *GenderEnum) ToNullInt16() null.Int16 {
	if ge == nil {
		return null.NewInt16(0, false)
	}
	return null.Int16From(int16(*ge))
}

// GenderEnumFromNullInt16 converts null.Int16 type to pointer of GenderEnum.
func GenderEnumFromNullInt16(value null.Int16) *GenderEnum {
	genderNumeric := value.Ptr()
	var gender *GenderEnum
	if genderNumeric != nil {
		genderValue := GenderEnum(*genderNumeric)
		gender = &genderValue
	}

	return gender
}
