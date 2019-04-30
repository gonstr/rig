package build

import (
	"fmt"

	"k8s.io/helm/pkg/strvals"
)

// createValueMap merges a value map with value arrays. Values defined in the
// value arrays supercede values in the value map
func createValueMap(valueMap map[string]interface{}, values []string, stringValues []string) (map[string]interface{}, error) {
	vals := make(map[string]interface{})

	if valueMap != nil {
		for k, v := range valueMap {
			vals[k] = v
		}
	}

	// User specified a value via --value
	for _, value := range values {
		if err := strvals.ParseInto(value, vals); err != nil {
			return nil, fmt.Errorf("failed parsing --value data: %s", err)
		}
	}

	// User specified a value via --string-value
	for _, value := range stringValues {
		if err := strvals.ParseIntoString(value, vals); err != nil {
			return nil, fmt.Errorf("failed parsing --string-value data: %s", err)
		}
	}

	return map[string]interface{}{
		"values": vals,
	}, nil
}
