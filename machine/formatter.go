package machine

import (
	poolV1 "github.com/Netflix/titus-controllers-api/api/resourcepool/v1"
	. "github.com/Netflix/titus-resource-pool/util"
)

func FormatMachineType(machineType *poolV1.MachineTypeConfig, options FormatterOptions) string {
	if options.Level != FormatDetails {
		return formatMachineTypeCompact(machineType)
	}
	return ToJSONString(machineType)
}

func formatMachineTypeCompact(machineType *poolV1.MachineTypeConfig) string {
	type Compact struct {
		Name            string
		ComputeResource poolV1.ComputeResource
	}
	value := Compact{
		Name:            machineType.Name,
		ComputeResource: machineType.Spec.ComputeResource,
	}
	return ToJSONString(value)
}