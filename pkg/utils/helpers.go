package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This function is used to detect two pieces of information:
// A flag telling if the workload is being overwritten by an external manager **after** original manager
// it will be true if the external manager wrote the field after original manager by comparing timestamps
// and if any field was altered by the external manager
// the second piece of information is the name of the external manager, regardless the flag value
func DetectExternalManager(originalManager string, managedFields []metav1.ManagedFieldsEntry) (bool, string) {

	overwrittenByExternalManager := false
	otherManager := ""

	// First, let's get the latest managed field entry
	// of the original manager

	managedByOriginalManager, managedFieldsV1, mfTime := DetectManagedFields(originalManager, managedFields)

	if !managedByOriginalManager {
		return overwrittenByExternalManager, otherManager
	}

	if managedFieldsV1 == nil {
		return overwrittenByExternalManager, otherManager
	}

	// Now, let's get the JSON paths of the managed fields
	// of the original manager in regex format

	lookFor, err := FieldsV1ToJSONPaths(managedFieldsV1, true)
	if err != nil {
		return overwrittenByExternalManager, otherManager
	}

	if len(lookFor) == 0 {
		return overwrittenByExternalManager, otherManager
	}

	mfPattern := strings.Join(lookFor, "|")
	matchFields := regexp.MustCompile(mfPattern)

	// managedFields: sorting by time
	sort.Slice(managedFields, func(i, j int) bool {
		// Ensure that Time is not nil, just in case
		if managedFields[i].Time == nil {
			return false
		}
		if managedFields[j].Time == nil {
			return true
		}
		// Compare the time values
		return managedFields[i].Time.Before(managedFields[j].Time)
	})

	for _, managedField := range managedFields {
		if managedField.FieldsV1 == nil {
			continue
		}
		// we want only updates, not creation
		if managedField.Operation == "Create" {
			continue
		}
		// we ignore the original manager
		if managedField.Manager == originalManager {
			continue
		}

		// normalize the managed fields V1 into JSON Path
		managedFieldsAsJSONPaths, err := FieldsV1ToJSONPaths(managedField.FieldsV1)
		if err != nil {
			continue
		}

		// regex match the external manager managed fields
		for _, mfPath := range managedFieldsAsJSONPaths {
			if matchFields.MatchString(mfPath) {
				otherManager = managedField.Manager
				if managedField.Time.After(mfTime.Time) {
					overwrittenByExternalManager = true
				}
			}

		}

	}

	return overwrittenByExternalManager, otherManager
}

// DetectManagedFieldsByStormForge
// returns the last entry of managedFieldsEntry of the
// managedFieldsEntry array that was managed by StormForge

func DetectManagedFields(originalManager string, managedFields []metav1.ManagedFieldsEntry) (bool, *metav1.FieldsV1, metav1.Time) {

	managedByOriginalManager := false

	var idxLatestField int
	var timeLatestField = metav1.Time{Time: time.Time{}}

	// not super required, but sorting the managedFields
	// by time
	sort.Slice(managedFields, func(i, j int) bool {
		// Ensure that Time is not nil, just in case
		if managedFields[i].Time == nil {
			return false
		}
		if managedFields[j].Time == nil {
			return true
		}
		// Compare the time values
		return managedFields[i].Time.Before(managedFields[j].Time)
	})

	// detect there is a field managed by originalManager and
	// the secure the index of last one
	for idx, managedField := range managedFields {
		if managedField.FieldsV1 == nil {
			continue
		}
		if managedField.Operation == "Create" {
			continue
		}
		if managedField.Manager == originalManager {
			managedByOriginalManager = true
			if managedField.Time == nil || managedField.Time.After(timeLatestField.Time) {
				idxLatestField = idx
				timeLatestField = *managedField.Time
			}
		}
	}

	if managedByOriginalManager {
		return managedByOriginalManager, managedFields[idxLatestField].FieldsV1, timeLatestField
	}

	return managedByOriginalManager, nil, timeLatestField
}

func MustParseTime(value string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(fmt.Sprintf("Error parsing date string: %v", err))
	}
	return parsedTime
}

func FieldsV1ToJSONPaths(fieldsV1 *metav1.FieldsV1, flags ...bool) ([]string, error) {

	paths := []string{}

	if fieldsV1 == nil {
		return paths, fmt.Errorf("fieldsV1 nil")
	}

	optionalRegExFlag := false

	if len(flags) > 0 {
		optionalRegExFlag = flags[0]
	}

	var fieldsMap map[string]interface{}
	if err := json.Unmarshal(fieldsV1.Raw, &fieldsMap); err != nil {
		return paths, err
	}

	extractPaths("", fieldsMap, &paths)

	sort.Strings(paths)

	// convert paths into regexes
	if optionalRegExFlag {

		regExPaths := []string{}

		for _, path := range paths {

			// adding escaled on each /
			escaped := strings.ReplaceAll(path, "/", `\/`)

			//replacing keys with wildcard
			re := regexp.MustCompile(`\[\{\".*?\"\}\]`)
			withoutKeys := re.ReplaceAllString(escaped, `*.*`)

			//replacing the last dot with wildcard
			re = regexp.MustCompile(`/\.$`)
			withoutLastDot := re.ReplaceAllString(withoutKeys, `/*.*`)

			regExPaths = append(regExPaths, withoutLastDot)

		}

		paths = MakeUnique(regExPaths)
		sort.Strings(paths)

	}

	return paths, nil

}

// Helper function to recursively extract paths from the fields map
func extractPaths(prefix string, m map[string]interface{}, paths *[]string) {
	for key, val := range m {
		// Handle both "f:" for fields and "k:" for keys
		trimmedKey := strings.TrimPrefix(key, "f:")
		if strings.HasPrefix(trimmedKey, "k:") {
			trimmedKey = fmt.Sprintf("[%s]", strings.TrimPrefix(trimmedKey, "k:"))
		}

		// Build the full path for the current key
		fullPath := fmt.Sprintf("%s/%s", prefix, trimmedKey)

		// Check if value is a nested map, if so recurse
		if nested, ok := val.(map[string]interface{}); ok {
			if len(nested) > 0 {
				extractPaths(fullPath, nested, paths)
			} else {
				*paths = append(*paths, fullPath)
			}
		}
	}

}

func MakeUnique(input []string) []string {
	uniqueMap := make(map[string]struct{}) // Use struct{} to save memory
	var result []string

	for _, entry := range input {
		if _, exists := uniqueMap[entry]; !exists {
			uniqueMap[entry] = struct{}{}
			result = append(result, entry)
		}
	}

	return result
}
