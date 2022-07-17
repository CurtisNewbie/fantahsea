package data

const (
	// record is deleted
	IS_DEL_Y IS_DEL = 1
	// record is not deleted
	IS_DEL_N IS_DEL = 0
)

type IS_DEL int8

type WEntity struct {
	ID         int64
	CreateTime WTime
	CreateBy   string
	UpdateTime WTime
	UpdateBy   string
}

// Check if the record is deleted
func IsDeleted(isDel IS_DEL) bool {
	return isDel == IS_DEL_Y
}
