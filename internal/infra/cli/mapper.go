package cli

import (
	storageErorrs "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"strings"
	"time"
)

const DatetimeFormat = "2006-01-02"

// MapPackageToDomain преобразует строку в PackageType с валидацией.
func MapPackageToDomain(pkgStr string) (pvz.PackageType, error) {
	switch pkgStr {
	case "bag":
		return pvz.PackageTypeBag, nil
	case "box":
		return pvz.PackageTypeBox, nil
	case "tape":
		return pvz.PackageTypeTape, nil
	case "bag+tape":
		return pvz.PackageTypeBagTape, nil
	case "box+tape":
		return pvz.PackageTypeBoxTape, nil
	default:
		return 0, storageErorrs.EntityValidationError("package", "invalid package type")
	}
}

// MapStringToTime преобразует временную строку по формату DatetimeFormat.
func MapStringToTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(DatetimeFormat, timeStr)
	if err != nil {
		return time.Time{}, storageErorrs.EntityValidationError("", "")
	}
	return t, nil
}

// MapTimeToString форматирует временную строку в необходимый по контракту тип
func MapTimeToString(t time.Time) string {
	return t.Format(DatetimeFormat)
}

func ActionTypeFromString(s string) pvz.ActionType {
	switch strings.ToLower(s) {
	case "issue":
		return pvz.ActionTypeIssue
	case "return":
		return pvz.ActionTypeReturn
	default:
		return pvz.ActionTypeUnspecified
	}
}

func MapStatusToString(s pvz.OrderStatus) string {
	switch s {
	case pvz.OrderStatusUnspecified:
		return "Unknown"
	case pvz.OrderStatusAccepted:
		return "Accepted"
	case pvz.OrderStatusIssued:
		return "Issued"
	case pvz.OrderStatusReturned:
		return "Returned"
	case pvz.OrderStatusDeleted:
		return "Deleted"
	default:
		return ""
	}
}
