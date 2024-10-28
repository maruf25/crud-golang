package utils

func IsSupportedFileType(fileType string) bool {
	supportedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	return supportedTypes[fileType]
}
