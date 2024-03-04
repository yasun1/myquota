/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains helper functions for the tests.

package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	// nolint
	. "github.com/onsi/gomega"

	client "github.com/openshift-online/ocm-sdk-go"
)

// Parse parses the given JSON data and returns a map of strings containing the result.
func Parse(data []byte) map[string]interface{} {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		panic(err)
	}
	return object
}

// DigString tries to find an attribute inside the given object with the given path, and returns its
// value, assuming that it is an string. For example, if the object is the result of parsing the
// following JSON document:
//
//	{
//		"kind": "Cluster",
//		"id": "123",
//		"flavour": {
//			"kind": "Flavour",
//			"id": "456",
//			"href": "/api/clusters_mgmt/v1/flavours/456"
//		}
//	}
//
// The the 'id' attribute can be obtained like this:
//
//	clusterID := DigString(object, "id")
//
// And the 'id' attribute inside the 'flavour' can be obtained like this:
//
//	flavourID := DigString(object, "flavour", "id")
//
// If there is no attribute with the given path then the return value will be an empty string.
func DigString(object interface{}, keys ...interface{}) string {
	switch result := dig(object, keys).(type) {
	case nil:
		return ""
	case string:
		return result
	case fmt.Stringer:
		return result.String()
	default:
		return fmt.Sprintf("%s", result)
	}
}

// DigBool tries to find an attribute inside the given object with the given path, and returns its
// value, assuming that it is a boolean. For example, if the object is the result of parsing the
// following JSON document:
//
//	{
//		"kind": "Cluster",
//		"id": "123",
//		"hasId": true,
//		"flavour": {
//			"kind": "Flavour",
//			"hasId": false,
//			"href": "/api/clusters_mgmt/v1/flavours/456"
//		}
//	}
//
// The the 'hasId' attribute can be obtained like this:
//
//	hasID := DigBool(object, "hasId")
//
// And the 'hasId' attribute inside the 'flavour' can be obtained like this:
//
//	flavourHasID := DigBool(object, "flavour", "hasId")
//
// If there is no attribute with the given path then the return value will be false.
func DigBool(object interface{}, keys ...interface{}) bool {
	switch result := dig(object, keys).(type) {
	case nil:
		return false
	case bool:
		return result
	case string:
		b, err := strconv.ParseBool(result)
		if err != nil {
			return false
		}
		return b
	default:
		return false
	}
}

// DigInt tries to find an attribute inside the given object with the given path, and returns its
// value, assuming that it is an integer. If there is no attribute with the given path then the test
// will be aborted with an error.
func DigInt(object interface{}, keys ...interface{}) int {
	value := dig(object, keys)
	// ExpectWithOffset(1, value).ToNot(BeNil())
	var result float64
	// ExpectWithOffset(1, value).To(BeAssignableToTypeOf(result))
	result = value.(float64)
	return int(result)
}

// DigFloat tries to find an attribute inside the given object with the given path, and returns its
// value, assuming that it is an floating point number. If there is no attribute with the given path
// then the test will be aborted with an error.
func DigFloat(object interface{}, keys ...interface{}) float64 {
	value := dig(object, keys)
	ExpectWithOffset(1, value).ToNot(BeNil())
	var result float64
	ExpectWithOffset(1, value).To(BeAssignableToTypeOf(result))
	result = value.(float64)
	return result
}

// DigObject tries to find an attribute inside the given object with the given path, and returns its
// value. If there is no attribute with the given path then the test will be aborted with an error.
func DigObject(object interface{}, keys ...interface{}) interface{} {
	value := dig(object, keys)
	ExpectWithOffset(1, value).ToNot(BeNil())
	return value
}

// DigArray tries to find an array inside the given object with the given path, and returns its
// value. If there is no attribute with the given path then the test will be aborted with an error.
func DigArray(object interface{}, keys ...interface{}) []interface{} {
	value := dig(object, keys)
	//ExpectWithOffset(1, value).ToNot(BeNil())
	var result []interface{}
	//ExpectWithOffset(1, value).To(BeAssignableToTypeOf(result))
	result = value.([]interface{})
	return result
}

func dig(object interface{}, keys []interface{}) interface{} {
	if object == nil || len(keys) == 0 {
		return nil
	}
	switch key := keys[0].(type) {
	case string:
		switch data := object.(type) {
		case map[string]interface{}:
			value := data[key]
			if len(keys) == 1 {
				return value
			}
			return dig(value, keys[1:])
		}
	case int:
		switch data := object.(type) {
		case []interface{}:
			value := data[key]
			if len(keys) == 1 {
				return value
			}
			return dig(value, keys[1:])
		}
	}
	return nil
}

// Template processes the given template using as data the set of name value pairs that are given as
// arguments. For example, to the following code:
//
//	result, err := Template(`
//		{
//			"name": "{{ .Name }}",
//			"flavour": {
//				"id": "{{ .Flavour }}"
//			}
//		}
//		`,
//		"Name", "mycluster",
//		"Flavour", "4",
//	)
//
// Produces the following result:
//
//	{
//		"name": "mycluster",
//		"flavour": {
//			"id": "4"
//		}
//	}
func Template(source string, args ...interface{}) string {
	// Check that there is an even number of args, and that the first of each pair is an string:
	count := len(args)
	ExpectWithOffset(1, count%2).To(
		Equal(0),
		"Template '%s' should have an even number of arguments, but it has %d",
		source, count,
	)
	for i := 0; i < count; i = i + 2 {
		name := args[i]
		_, ok := name.(string)
		ExpectWithOffset(1, ok).To(
			BeTrue(),
			"Argument %d of template '%s' is a key, so it should be a string, "+
				"but its type is %T",
			i, source, name,
		)
	}

	// Put the variables in the map that will be passed as the data object for the execution of
	// the template:
	data := make(map[string]interface{})
	for i := 0; i < count; i = i + 2 {
		name := args[i].(string)
		value := args[i+1]
		data[name] = value
	}

	// Parse the template:
	tmpl, err := template.New("").Parse(source)
	ExpectWithOffset(1, err).ToNot(
		HaveOccurred(),
		"Can't parse template '%s': %v",
		source, err,
	)

	// Execute the template:
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, data)
	ExpectWithOffset(1, err).ToNot(
		HaveOccurred(),
		"Can't execute template '%s': %v",
		source, err,
	)
	return buffer.String()
}

func readGetResult(response *client.Response, err error, expectedStatus int) map[string]interface{} {
	checkGetResponse(response, err, expectedStatus)
	return Parse(response.Bytes())
}

func checkGetResponse(response *client.Response, err error, expectedStatus int) {
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Status()).To(Equal(expectedStatus))
}

func readPostResult(response *client.Response, err error, expectedStatus int) map[string]interface{} {
	checkPostResponse(response, err, expectedStatus)
	return Parse(response.Bytes())
}

func checkPostResponse(response *client.Response, err error, expectedStatus int) {
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Status()).To(Equal(expectedStatus))
}

func checkDeleteResponse(response *client.Response, err error, expectedStatus int) {
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Status()).To(Equal(expectedStatus))
}

func readPatchResponse(response *client.Response, err error, expectedStatus int) map[string]interface{} {
	checkPatchResponse(response, err, expectedStatus)
	return Parse(response.Bytes())
}

func checkPatchResponse(response *client.Response, err error, expectedStatus int) {
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Status()).To(Equal(expectedStatus))
}

func isLocal(url string) bool {
	return strings.Contains(url, "localhost") ||
		strings.Contains(url, "127.0.0.1") ||
		strings.Contains(url, "::1")
}

// runAttempt will run a function (attempt), until the function returns false - meaning no further attempts should
// follow, or until the number of attempts reached maxAttempts. Between each 2 attempts, it will wait for a given
// delay time.
// In case maxAttempts have been reached, an error will be returned, with the latest attempt result.
// The attempt function should return true as long as another attempt should be made, false when no further attempts
// are required - i.e. the attempt succeeded, and the result is available as returned value.
func runAttempt(attempt func() (interface{}, bool), maxAttempts int, delay time.Duration) (interface{}, error) {
	var result interface{}
	for i := 0; i < maxAttempts; i++ {
		fmt.Println("Running attempt", i)
		result, toContinue := attempt()
		if toContinue {
			fmt.Println("Need to continue for another attempt")
		} else {
			fmt.Println("Attempt successful, returning result")
			return result, nil
		}
		fmt.Println("Sleeping for", delay)
		time.Sleep(delay)
	}
	return result, fmt.Errorf("Got to max attempts %d", maxAttempts)
}
