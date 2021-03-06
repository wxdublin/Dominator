package rpcd

import (
	"github.com/Symantec/Dominator/lib/filter"
	"github.com/Symantec/Dominator/proto/sub"
)

func (t *rpcType) SetConfiguration(request sub.SetConfigurationRequest,
	reply *sub.SetConfigurationResponse) error {
	var response sub.SetConfigurationResponse
	scannerConfiguration.FsScanContext.GetContext().SetSpeedPercent(
		request.ScanSpeedPercent)
	scannerConfiguration.NetworkReaderContext.SetSpeedPercent(
		request.NetworkSpeedPercent)
	newFilter, err := filter.NewFilter(request.ScanExclusionList)
	if err != nil {
		return err
	}
	scannerConfiguration.ScanFilter = newFilter
	response.Success = true
	*reply = response
	return nil
}
