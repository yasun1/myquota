package helpers

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// NewRand returns a rand with the time seed
func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Join will link the strings with "."
func Join(s ...string) string {
	return strings.Join(s, ".")
}

// IsSorted will return whether the array is sorted by mode
func IsSorted(arry []string, mode string) (flag bool) {
	switch mode {
	case "desc":
		for i := 0; i < len(arry)-1; i++ {
			if arry[i] < arry[i+1] {
				flag = false
			}
		}
		flag = true
	case "asc":
		for i := 0; i < len(arry)-1; i++ {
			if arry[i] > arry[i+1] {
				flag = false
			}
		}
		flag = true
	default:
		for i := 0; i < len(arry)-1; i++ {
			if arry[i] > arry[i+1] {
				flag = false
			}
		}
		flag = true
	}
	return
}

// Min will return the minimize value
func Min(a int, b int) int {
	if a >= b {
		return b
	}
	return a
}

// Min will return the minimize value
func Max(a int, b int) int {
	if a <= b {
		return b
	}
	return a
}

// Contains will return bool balue that whether the arry contains string val
func Contains(arry []string, val string) (index int, flag bool) {
	var i int
	flag = false
	if len(arry) == 0 {
		return
	}
	for i = 0; i < len(arry); i++ {
		if arry[i] == val {
			index = i
			flag = true
			break
		}
	}
	return
}

// EndsWith will return the bool value that whether the st is end with substring
func EndsWith(st string, substring string) (flag bool) {

	if len(st) < len(substring) {
		return
	}
	flag = (st[len(st)-len(substring):] == substring)
	return
}

// Strip will return the value striped with substring
func Strip(st string, substring string) string {
	st = Lstrip(st, substring)
	st = Rstrip(st, substring)
	return st
}

// Lstrip will return the string left striped with substring
func Lstrip(st string, substring string) string {
	if StartsWith(st, substring) {
		st = st[len(substring):]
	}
	return st
}

// Rstrip will return the string right striped with substring
func Rstrip(st string, substring string) string {
	if EndsWith(st, substring) {
		st = st[:len(st)-len(substring)]
	}
	return st
}

// StartsWith return bool whether st start with substring
func StartsWith(st string, substring string) (flag bool) {
	if len(st) < len(substring) {
		return
	}
	flag = (st[:len(substring)] == substring)
	return
}

// NegateBoolToString reverts the boolean to its oppositely value as a string.
func NegateBoolToString(value bool) string {
	boolString := "true"
	if value {
		boolString = "false"
	}
	return boolString
}

// ConvertMapToJSONString converts the map to a json string.
func ConvertMapToJSONString(inputMap map[string]interface{}) string {
	jsonBytes, _ := json.Marshal(inputMap)
	return string(jsonBytes)
}

func ConvertStringToInt(mystring string) int {
	myint, _ := strconv.Atoi(mystring)
	return myint
}

func ConvertStructToMap(s interface{}) map[string]interface{} {
	structMap := make(map[string]interface{})

	j, _ := json.Marshal(s)
	err := json.Unmarshal(j, &structMap)
	if err != nil {
		panic(err)
	}

	return structMap
}

func ConvertStructToString(s interface{}) string {
	structMap := ConvertStructToMap(s)
	return ConvertMapToJSONString(structMap)
}

func IsInMap(inputMap map[string]interface{}, key string) bool {
	_, contain := inputMap[key]
	return contain
}

// NeedFiltered will return the attribute that should be filtered
// filterList should be array with regex like ["excluded\..+","excluded_[\s\S]+"]
func NeedFiltered(filterList []string, key string) bool {
	for _, regex := range filterList {
		pattern := regexp.MustCompile(regex)
		if pattern.MatchString(key) {
			return true
		}
	}
	return false
}
