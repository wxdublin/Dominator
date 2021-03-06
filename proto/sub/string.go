package sub

import (
	"fmt"
)

func (configuration Configuration) String() string {
	retval := fmt.Sprintf("ScanSpeedPercent: %d\nNetworkSpeedPercent: %d",
		configuration.ScanSpeedPercent, configuration.NetworkSpeedPercent)
	if len(configuration.ScanExclusionList) > 0 {
		retval += "\n" + "ScanExclusionList:"
		for _, exclusion := range configuration.ScanExclusionList {
			retval += "\n  " + exclusion
		}
	}
	return retval
}

func (configuration GetConfigurationResponse) String() string {
	return Configuration(configuration).String()
}
