package category

import "github.com/gofrs/uuid/v5"

type CategoryID uuid.UUID

func (id CategoryID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func NewCategoryID() CategoryID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return CategoryID(id)
}

func ParseCategoryID(value string) (CategoryID, error) {
	id, err := uuid.FromString(value)
	return CategoryID(id), err
}

func (id CategoryID) String() string {
	return uuid.UUID(id).String()
}
