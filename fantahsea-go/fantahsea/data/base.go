package data

const (
	// record is deleted
	IS_DEL_Y int8 = 1
	// record is not deleted
	IS_DEL_N int8 = 0
)

// Check if the record is deleted
func IsDeleted(isDel int8) bool {
	return isDel == IS_DEL_Y
}
