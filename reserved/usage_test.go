package reserved

import (
	poolV1 "github.com/Netflix/titus-controllers-api/api/resourcepool/v1"
	. "github.com/Netflix/titus-resource-pool/resourcepool"
	"github.com/Netflix/titus-resource-pool/util"
	"github.com/stretchr/testify/require"
	k8sCore "k8s.io/api/core/v1"
	"testing"
)

func TestNewCapacityReservationUsage(t *testing.T) {
	pool := ButResourcePoolName(EmptyResourcePool(), PoolNameIntegration)
	pool.Spec.ResourceCount = 20

	node := util.NewNode("node1", PoolNameIntegration, util.R5Metal())

	pod1 := util.ButPodResourcePools(util.NewRandomNotScheduledPod(), PoolNameIntegration)
	pod1 = util.ButPodCapacityGroup(pod1, "group1")
	pod1 = util.ButPodAssignedToNode(pod1, node)

	poolSnapshot := NewStaticResourceSnapshot(pool, []*poolV1.MachineTypeConfig{}, []*k8sCore.Node{node}, []*k8sCore.Pod{pod1},
		0, true)

	group1 := NewCapacityGroup("group1", PoolNameIntegration)
	group1.Spec.InstanceCount = 10
	capacityGroups := []*poolV1.CapacityGroup{
		group1,
		NewCapacityGroup("group2", PoolNameIntegration),
	}

	usage := NewCapacityReservationUsage(poolSnapshot, capacityGroups)
	require.Len(t, usage.InCapacityGroup, 2)

	expectedGroup1Allocated := util.FromPodToComputeResource(pod1)
	expectedGroup1Unallocated := CapacityGroupResources(capacityGroups[0]).Sub(expectedGroup1Allocated)
	require.Equal(t, expectedGroup1Allocated, usage.InCapacityGroup["group1"].Allocated)
	require.Equal(t, expectedGroup1Unallocated, usage.InCapacityGroup["group1"].Unallocated)
	require.Equal(t, expectedGroup1Allocated, usage.AllReserved.Allocated)
	require.Equal(t, expectedGroup1Unallocated.Add(CapacityGroupResources(capacityGroups[1])), usage.AllReserved.Unallocated)
}