package sdk

// Функция для разбивки среза на несколько срезов по указанному размеру.
func ChunkSlice[T any](slice []T, size int) [][]T {
	var chunks [][]T

	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
