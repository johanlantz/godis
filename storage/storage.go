package storage

var storage = make(map[string]StorageEntry)

// Return a copy so caller does not get a data access pointer.
func Get(key string) StorageEntry {
	v, exists := storage[key]
	if !exists {
		return StorageEntry{}
	} else {
		return v
	}
}

// Set can't fail, if a previous existing dataType is different, we overwrite.
func Set(key string, entry StorageEntry) {
	storage[key] = entry
}
