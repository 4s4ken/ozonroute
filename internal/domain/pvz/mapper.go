package pvz

import (
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"time"
)

const DatetimeFormat = "2006-01-02"

// MapPackageToDomain преобразует строку в PackageType с валидацией.
func MapPackageToDomain(pkgStr string) (PackageType, error) {
	switch pkgStr {
	case string(Bag), string(Box), string(Film), string(BagAndFilm), string(BoxAndFilm):
		return PackageType(pkgStr), nil
	default:
		return "", businesserror.EntityValidationError("", "")
	}
}

// MapPackageFromDomain преобразует PackageType в строку.
func MapPackageFromDomain(pkg PackageType) string {
	return string(pkg)
}

// MapStringToTime преобразует временную строку по формату DatetimeFormat.
func MapStringToTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(DatetimeFormat, timeStr)
	if err != nil {
		return time.Time{}, businesserror.EntityValidationError("", "")
	}
	return t, nil
}

// MapTimeToString форматирует временную строку в необходимый по контракту тип
func MapTimeToString(t time.Time) string {
	return t.Format(DatetimeFormat)
}
