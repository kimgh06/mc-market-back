package nullable

import "database/sql"

func PointerToInt16(pointer *int16) sql.NullInt16 {
	if pointer == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{Int16: *pointer, Valid: true}
}

func UPointerToInt16(pointer *uint16) sql.NullInt16 {
	if pointer == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{Int16: int16(*pointer), Valid: true}
}

func PointerToInt32(pointer *int32) sql.NullInt32 {
	if pointer == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *pointer, Valid: true}
}

func UPointerToInt32(pointer *uint32) sql.NullInt32 {
	if pointer == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: int32(*pointer), Valid: true}
}

func PointerToInt64(pointer *int64) sql.NullInt64 {
	if pointer == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *pointer, Valid: true}
}

func UPointerToInt64(pointer *uint64) sql.NullInt64 {
	if pointer == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*pointer), Valid: true}
}

func Int16ToPointer(value sql.NullInt16) *int16 {
	if value.Valid {
		return &value.Int16
	}
	return nil
}

func Int16ToUPointer(value sql.NullInt16) *uint16 {
	if value.Valid {
		casted := uint16(value.Int16)
		return &casted
	}
	return nil
}

func Int32ToPointer(value sql.NullInt32) *int32 {
	if value.Valid {
		return &value.Int32
	}
	return nil
}

func Int32ToUPointer(value sql.NullInt32) *uint32 {
	if value.Valid {
		casted := uint32(value.Int32)
		return &casted
	}
	return nil
}

func Int64ToPointer(value sql.NullInt64) *int64 {
	if value.Valid {
		return &value.Int64
	}
	return nil
}

func Int64ToUPointer(value sql.NullInt64) *uint64 {
	if value.Valid {
		casted := uint64(value.Int64)
		return &casted
	}
	return nil
}
