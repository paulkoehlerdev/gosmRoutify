package datadecoder

// Make tags map from stringtable and two parallel arrays of IDs.
func extractTags(stringTable [][]byte, keyIDs []uint32, valueIDs []uint32) map[string]string {
	tags := make(map[string]string, len(keyIDs))
	for index, keyID := range keyIDs {
		key := stringTable[keyID]
		val := stringTable[valueIDs[index]]
		tags[string(key)] = string(val)
	}
	return tags
}

type tagUnpacker struct {
	stringTable [][]byte
	keysVals    []int32
	index       int
}

// Make tags map from stringtable and array of IDs (used in DenseNodes encoding).
func (tu *tagUnpacker) next() map[string]string {
	tags := make(map[string]string)
	for tu.index < len(tu.keysVals) {
		keyID := tu.keysVals[tu.index]
		tu.index++
		if keyID == 0 {
			break
		}

		valID := tu.keysVals[tu.index]
		tu.index++

		key := tu.stringTable[keyID]
		val := tu.stringTable[valID]
		tags[string(key)] = string(val)
	}
	return tags
}
