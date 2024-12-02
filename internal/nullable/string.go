package nullable

import "database/sql"

func PointerToString(pointer *string) sql.NullString {
	if pointer == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *pointer, Valid: true}
}

func StringToPointer(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}
