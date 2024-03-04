package helpers

import (

	// nolint
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/onsi/gomega"

	client "github.com/openshift-online/ocm-sdk-go"
)

const (
	UnderscoreConnector string = "_"
	DotConnector        string = "."
	HyphenConnector     string = "-"
)

// Dig Expose dig to used by others
func Dig(object interface{}, keys []interface{}) interface{} {
	return dig(object, keys)
}

// CheckResponse is used to checking
func CheckResponse(response *client.Response, err error, expectedStatus int) {
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Status()).To(Equal(expectedStatus))
}

// FlatMap parses the data, and stores all the attributes and the sub attributes at the same level with the prefix key.
func FlatMap(inputMap map[string]interface{}, outputMap map[string]interface{}, key string, connector string) {
	if len(inputMap) == 0 {
		outputMap[key] = ""
	} else {
		for k, v := range inputMap {
			outKey := fmt.Sprintf("%s%s%s", key, connector, k)
			if key == "" {
				outKey = k
			}
			// fmt.Printf("[debug1][%v]%v: %+v, %+v\n", outKey, k, reflect.TypeOf(v), v)
			// time.Sleep(1 * time.Second)
			switch v := v.(type) {
			case map[string]interface{}:
				FlatMap(v, outputMap, outKey, connector)
			case []interface{}:
				// If the keys are in an array, it will flat the array with the index, etc, items_0_id or items.0.id
				for nk, nv := range v {
					newOutKey := fmt.Sprintf("%s%s%d", outKey, connector, nk)
					// fmt.Printf("[debug2][%v]%v: %+v, %+v\n", newOutKey, nk, reflect.TypeOf(nv), nv)
					switch nv := nv.(type) {
					case map[string]interface{}:
						FlatMap(nv, outputMap, newOutKey, connector)
					default:
						outputMap[newOutKey] = nv
					}
				}
			default:
				outputMap[outKey] = v
			}
		}
	}

}

// MapStructure will map the map to the address of the structre *i
func MapStructure(m map[string]interface{}, i interface{}) error {
	jsonbody, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonbody, i)
	if err != nil {
		return err
	}
	return nil
}

// FlatInitialMap parses the data to store all the attributes at the first level.
func FlatInitialMap(inputMap map[string]interface{}, connector string) map[string]interface{} {
	outputMap := make(map[string]interface{})
	FlatMap(inputMap, outputMap, "", connector)
	return outputMap
}

// ConvertRequestBodyByAttr will return the request bodies by attributes
func ConvertRequestBodyByAttr(inputString string, connector string) (outArray []string, err error) {
	inputMap := Parse([]byte(inputString))
	outMap := FlatInitialMap(inputMap, connector)
	outArray = ConvertFlatMapToArray(outMap, connector)
	return
}

// ConvertFlatMapToArray will return the request body converted from flatmap
// {"id":"xxx", "product.id":"xxxx", "managed":true}
// will be converted to
// ["{"id":"xxxx"}","{"product":{"id":"xxxx"}}",...]
func ConvertFlatMapToArray(flatMap map[string]interface{}, connector string) (outArray []string) {
	for key, value := range flatMap {
		keys := strings.Split(key, connector)
		resultmap := make(map[string]interface{})
		for i := len(keys) - 1; i >= 0; i-- {
			middleMap := make(map[string]interface{})
			if i == len(keys)-1 {
				middleMap[keys[i]] = value

			} else {
				middleMap[keys[i]] = resultmap
			}
			resultmap = middleMap

		}
		requestBody, err := json.Marshal(resultmap)
		if err != nil {
			return
		}
		outArray = append(outArray, string(requestBody))
	}
	return
}
