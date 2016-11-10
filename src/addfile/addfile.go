package addfile

import (
	"encoding/base64"
)

const (
	str_file = "UmFyIRoHAM+QcwAADQAAAAAAAAB3FHQgkC0AHgAAACAAAAAC+ZyN0JOGD0kdMwgAIAAAAHRlc3QudHh0APD1CWQI1UvtCnL+Be/f4siKD0eDDsGneDQ1ofBaqv18QSjEPXsAQAcA"
	filename = "test.rar"
)

func GetAdditionFile() (file []byte) {
	decoded, _ := base64.StdEncoding.DecodeString(str_file)
	return decoded
}

func GetAdditionFileName() (name string) {
	return filename
}
