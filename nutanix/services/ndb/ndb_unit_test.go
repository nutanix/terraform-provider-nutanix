package ndb

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitNDBSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, DataSourceNutanixNDBClone().InternalValidate(DataSourceNutanixNDBClone().Schema, false))
	require.NoError(t, DataSourceNutanixNDBClones().InternalValidate(DataSourceNutanixNDBClones().Schema, false))
	require.NoError(t, DataSourceNutanixEraCluster().InternalValidate(DataSourceNutanixEraCluster().Schema, false))
	require.NoError(t, DataSourceNutanixEraClusters().InternalValidate(DataSourceNutanixEraClusters().Schema, false))
	require.NoError(t, DataSourceNutanixEraDatabase().InternalValidate(DataSourceNutanixEraDatabase().Schema, false))
	require.NoError(t, DataSourceNutanixEraDatabases().InternalValidate(DataSourceNutanixEraDatabases().Schema, false))
	require.NoError(t, DataSourceNutanixNDBDBServer().InternalValidate(DataSourceNutanixNDBDBServer().Schema, false))
	require.NoError(t, DataSourceNutanixNDBDBServers().InternalValidate(DataSourceNutanixNDBDBServers().Schema, false))
	require.NoError(t, DataSourceNutanixNDBMaintenanceWindow().InternalValidate(DataSourceNutanixNDBMaintenanceWindow().Schema, false))
	require.NoError(t, DataSourceNutanixNDBMaintenanceWindows().InternalValidate(DataSourceNutanixNDBMaintenanceWindows().Schema, false))
	require.NoError(t, DataSourceNutanixEraNetwork().InternalValidate(DataSourceNutanixEraNetwork().Schema, false))
	require.NoError(t, DataSourceNutanixEraNetworks().InternalValidate(DataSourceNutanixEraNetworks().Schema, false))
	require.NoError(t, DataSourceNutanixEraProfile().InternalValidate(DataSourceNutanixEraProfile().Schema, false))
	require.NoError(t, DataSourceNutanixEraProfiles().InternalValidate(DataSourceNutanixEraProfiles().Schema, false))
	require.NoError(t, DataSourceNutanixNDBProfileAvailableIPs().InternalValidate(DataSourceNutanixNDBProfileAvailableIPs().Schema, false))
	require.NoError(t, DataSourceNutanixEraSLA().InternalValidate(DataSourceNutanixEraSLA().Schema, false))
	require.NoError(t, DataSourceNutanixEraSLAs().InternalValidate(DataSourceNutanixEraSLAs().Schema, false))
	require.NoError(t, DataSourceNutanixNDBSnapshot().InternalValidate(DataSourceNutanixNDBSnapshot().Schema, false))
	require.NoError(t, DataSourceNutanixNDBSnapshots().InternalValidate(DataSourceNutanixNDBSnapshots().Schema, false))
	require.NoError(t, DataSourceNutanixNDBTag().InternalValidate(DataSourceNutanixNDBTag().Schema, false))
	require.NoError(t, DataSourceNutanixNDBTags().InternalValidate(DataSourceNutanixNDBTags().Schema, false))
	require.NoError(t, DataSourceNutanixNDBTimeMachine().InternalValidate(DataSourceNutanixNDBTimeMachine().Schema, false))
	require.NoError(t, DataSourceNutanixNDBTmsCapability().InternalValidate(DataSourceNutanixNDBTmsCapability().Schema, false))
	require.NoError(t, DataSourceNutanixNDBTimeMachines().InternalValidate(DataSourceNutanixNDBTimeMachines().Schema, false))
	require.NoError(t, ResourceNutanixNDBAuthorizeDBServer().InternalValidate(ResourceNutanixNDBAuthorizeDBServer().Schema, true))
	require.NoError(t, ResourceNutanixNDBClone().InternalValidate(ResourceNutanixNDBClone().Schema, true))
	require.NoError(t, ResourceNutanixNDBCloneRefresh().InternalValidate(ResourceNutanixNDBCloneRefresh().Schema, true))
	require.NoError(t, ResourceNutanixNDBCluster().InternalValidate(ResourceNutanixNDBCluster().Schema, true))
	require.NoError(t, ResourceDatabaseInstance().InternalValidate(ResourceDatabaseInstance().Schema, true))
	require.NoError(t, ResourceNutanixNDBDatabaseRestore().InternalValidate(ResourceNutanixNDBDatabaseRestore().Schema, true))
	require.NoError(t, ResourceNutanixNDBScaleDatabase().InternalValidate(ResourceNutanixNDBScaleDatabase().Schema, true))
	require.NoError(t, ResourceNutanixNDBDatabaseSnapshot().InternalValidate(ResourceNutanixNDBDatabaseSnapshot().Schema, true))
	require.NoError(t, ResourceNutanixNDBServerVM().InternalValidate(ResourceNutanixNDBServerVM().Schema, true))
	require.NoError(t, ResourceNutanixNDBRegisterDBServer().InternalValidate(ResourceNutanixNDBRegisterDBServer().Schema, true))
	require.NoError(t, ResourceNutanixNDBLinkedDB().InternalValidate(ResourceNutanixNDBLinkedDB().Schema, true))
	require.NoError(t, ResourceNutanixNDBLogCatchUps().InternalValidate(ResourceNutanixNDBLogCatchUps().Schema, true))
	require.NoError(t, ResourceNutanixNDBMaintenanceTask().InternalValidate(ResourceNutanixNDBMaintenanceTask().Schema, true))
	require.NoError(t, ResourceNutanixNDBMaintenanceWindow().InternalValidate(ResourceNutanixNDBMaintenanceWindow().Schema, true))
	require.NoError(t, ResourceNutanixNDBNetwork().InternalValidate(ResourceNutanixNDBNetwork().Schema, true))
	require.NoError(t, ResourceNutanixNDBProfile().InternalValidate(ResourceNutanixNDBProfile().Schema, true))
	require.NoError(t, ResourceNutanixNDBRegisterDatabase().InternalValidate(ResourceNutanixNDBRegisterDatabase().Schema, true))
	require.NoError(t, ResourceNutanixNDBSla().InternalValidate(ResourceNutanixNDBSla().Schema, true))
	require.NoError(t, ResourceNutanixNDBSoftwareVersionProfile().InternalValidate(ResourceNutanixNDBSoftwareVersionProfile().Schema, true))
	require.NoError(t, ResourceNutanixNDBStretchedVlan().InternalValidate(ResourceNutanixNDBStretchedVlan().Schema, true))
	require.NoError(t, ResourceNutanixNDBTags().InternalValidate(ResourceNutanixNDBTags().Schema, true))
	require.NoError(t, ResourceNutanixNDBTmsCluster().InternalValidate(ResourceNutanixNDBTmsCluster().Schema, true))
}

func TestUnitNDBContextAndBoolHelpers(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		t.Parallel()
		_, ok := FromContext(context.Background())
		require.False(t, ok)
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		ctx := NewContext(context.Background(), "")
		got, ok := FromContext(ctx)
		require.True(t, ok)
		require.Equal(t, "", got)
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		ctx := NewContext(context.Background(), "db-1")
		got, ok := FromContext(ctx)
		require.True(t, ok)
		require.Equal(t, "db-1", got)
		b, ok := tryToConvertBool("true")
		require.True(t, ok)
		require.True(t, b)
		_, ok = tryToConvertBool("nope")
		require.False(t, ok)
	})
}

func TestUnitExpandPrimarySLAAndIPInfos(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandPrimarySLA(nil))
		require.Nil(t, expandIPInfos(nil))
		require.Nil(t, expandTags(nil))
		require.Nil(t, buildSLADetails(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandPrimarySLA([]interface{}{}))
		require.Nil(t, expandIPInfos([]interface{}{}))
		require.Nil(t, expandTags([]interface{}{}))
		require.Nil(t, buildSLADetails([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		sla := expandPrimarySLA([]interface{}{map[string]interface{}{
			"sla_id":         "sla-1",
			"nx_cluster_ids": []interface{}{"c1"},
		}})
		require.Equal(t, "sla-1", utils.StringValue(sla.SLAID))
		ips := expandIPInfos([]interface{}{map[string]interface{}{
			"ip_type":      "static",
			"ip_addresses": []interface{}{"10.0.0.1"},
		}})
		require.Len(t, ips, 1)
		tags := expandTags([]interface{}{map[string]interface{}{"tag_name": "n", "value": "v"}})
		require.Len(t, tags, 1)
	})
}

func TestUnitNDBFlattenExpandHelpers(t *testing.T) {
	t.Parallel()

	t.Run("expandNodesProperties", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			s := schema.NewSet(schema.HashString, nil)
			_ = expandNodesProperties(s)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			s := schema.NewSet(schema.HashResource(&schema.Resource{Schema: map[string]*schema.Schema{"name": {Type: schema.TypeString}, "value": {Type: schema.TypeString}}}), []interface{}{map[string]interface{}{"name": "n", "value": "true"}})
			_ = expandNodesProperties(s)
		})
	})

	t.Run("flattenEntityCounts", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenEntityCounts(&era.EntityCounts{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenEntityCounts(&era.EntityCounts{})
		})
	})

	t.Run("flattenEngineCounts", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenEngineCounts(&era.EngineCounts{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenEngineCounts(&era.EngineCounts{})
		})
	})

	t.Run("flattenProfileTmsCount", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenProfileTmsCount(&era.ProfileTimeMachinesCount{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenProfileTmsCount(&era.ProfileTimeMachinesCount{})
		})
	})

	t.Run("flattenProfilesCount", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenProfilesCount(&era.ProfilesEntity{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenProfilesCount(&era.ProfilesEntity{})
		})
	})

	t.Run("flattenClustersResponse", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenClustersResponse(&era.ClusterListResponse{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClustersResponse(&era.ClusterListResponse{})
		})
	})

	t.Run("flattenClusterProperties", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterProperties([]*era.Properties{&era.Properties{}})
		})
	})

	t.Run("flattenResourceConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenResourceConfig(&era.Resourceconfig{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenResourceConfig(&era.Resourceconfig{})
		})
	})

	t.Run("flattenDBInstanceProperties", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBInstanceProperties([]*era.DBInstanceProperties{&era.DBInstanceProperties{}})
		})
	})

	t.Run("flattenDBInstanceMetadata", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBInstanceMetadata(&era.DBInstanceMetadata{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBInstanceMetadata(&era.DBInstanceMetadata{})
		})
	})

	t.Run("flattenDBNodes", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBNodes([]era.Databasenodes{{}})
		})
	})

	t.Run("flattenDBLinkedDbs", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBLinkedDbs([]era.Linkeddatabases{{}})
		})
	})

	t.Run("flattenDBProtectionDomain", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBProtectionDomain(&era.Protectiondomain{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBProtectionDomain(&era.Protectiondomain{})
		})
	})

	t.Run("flattenDBTags", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBTags([]*era.Tags{&era.Tags{}})
		})
	})

	t.Run("flattenDBInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBInfo(&era.Info{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBInfo(&era.Info{})
		})
	})

	t.Run("flattenBpgConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenBpgConfig(&era.BpgConfigs{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenBpgConfig(&era.BpgConfigs{})
		})
	})

	t.Run("flattenEraExpiryDetails", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenEraExpiryDetails(&era.DBExpiryDetails{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenEraExpiryDetails(&era.DBExpiryDetails{})
		})
	})

	t.Run("flattenEraRefreshDetails", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenEraRefreshDetails(&era.DBRefreshDetails{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenEraRefreshDetails(&era.DBRefreshDetails{})
		})
	})

	t.Run("flattenDBTimeMachine", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBTimeMachine(&era.TimeMachine{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBTimeMachine(&era.TimeMachine{})
		})
	})

	t.Run("flattenDBSLA", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBSLA(&era.ListSLAResponse{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBSLA(&era.ListSLAResponse{})
		})
	})

	t.Run("flattenSchedule", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenSchedule(&era.Schedule{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSchedule(&era.Schedule{})
		})
	})

	t.Run("flattenContinousSch", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenContinousSch(&era.Continuousschedule{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenContinousSch(&era.Continuousschedule{})
		})
	})

	t.Run("flattenMonthlySchedule", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenMonthlySchedule(&era.Monthlyschedule{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenMonthlySchedule(&era.Monthlyschedule{})
		})
	})

	t.Run("flattenQuartelySchedule", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenQuartelySchedule(&era.Quartelyschedule{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenQuartelySchedule(&era.Quartelyschedule{})
		})
	})

	t.Run("flattenSnapshotTimeOfDay", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnapshotTimeOfDay(&era.Snapshottimeofday{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnapshotTimeOfDay(&era.Snapshottimeofday{})
		})
	})

	t.Run("flattenWeeklySchedule", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenWeeklySchedule(&era.Weeklyschedule{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenWeeklySchedule(&era.Weeklyschedule{})
		})
	})

	t.Run("flattenYearlylySchedule", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenYearlylySchedule(&era.Yearlyschedule{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenYearlylySchedule(&era.Yearlyschedule{})
		})
	})

	t.Run("flattenDailySchedule", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDailySchedule(&era.Dailyschedule{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDailySchedule(&era.Dailyschedule{})
		})
	})

	t.Run("flattenTimeMachineMetadata", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenTimeMachineMetadata(&era.TimeMachineMetadata{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenTimeMachineMetadata(&era.TimeMachineMetadata{})
		})
	})

	t.Run("flattenDeRegiserInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDeRegiserInfo(&era.DeregisterInfo{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDeRegiserInfo(&era.DeregisterInfo{})
		})
	})

	t.Run("flattenDatabaseIntancesList", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDatabaseIntancesList(&era.ListDatabaseInstance{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDatabaseIntancesList(&era.ListDatabaseInstance{})
		})
	})

	t.Run("flattenDBServerMetadata", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBServerMetadata(&era.DBServerMetadata{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBServerMetadata(&era.DBServerMetadata{})
		})
	})

	t.Run("flattenDBServerVMInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBServerVMInfo(&era.VMInfo{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBServerVMInfo(&era.VMInfo{})
		})
	})

	t.Run("flattenVMNetworkInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenVMNetworkInfo([]*era.NetworkInfo{&era.NetworkInfo{}})
		})
	})

	t.Run("flattenAccessInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenAccessInfo([]*era.AccessInfo{&era.AccessInfo{}})
		})
	})

	t.Run("flattenDBServerVMResponse", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBServerVMResponse(&era.ListDBServerVMResponse{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBServerVMResponse(&era.ListDBServerVMResponse{})
		})
	})

	t.Run("flattenEntityTaskAssoc", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenEntityTaskAssoc([]*era.MaintenanceTasksResponse{&era.MaintenanceTasksResponse{}})
		})
	})

	t.Run("flattenEntityTaskPayload", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenEntityTaskPayload(&era.Payload{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenEntityTaskPayload(&era.Payload{})
		})
	})

	t.Run("flattenPrePostCommand", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenPrePostCommand(&era.PrePostCommand{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenPrePostCommand(&era.PrePostCommand{})
		})
	})

	t.Run("flattenMaintenanceWindowsResponse", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenMaintenanceWindowsResponse(&era.ListMaintenanceWindowResponse{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenMaintenanceWindowsResponse(&era.ListMaintenanceWindowResponse{})
		})
	})

	t.Run("flattenIPAddress", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPAddress([]*era.IPAddresses{&era.IPAddresses{}})
		})
	})

	t.Run("flattenNetworkListResponse", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenNetworkListResponse(&era.ListNetworkResponse{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenNetworkListResponse(&era.ListNetworkResponse{})
		})
	})

	t.Run("flattenNetworkProperties", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenNetworkProperties([]*era.Properties{&era.Properties{}})
		})
	})

	t.Run("flattenClusterAvailability", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterAvailability([]*era.Clusteravailability{&era.Clusteravailability{}})
		})
	})

	t.Run("flattenVersions", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenVersions([]*era.Versions{&era.Versions{}})
		})
	})

	t.Run("flattenProperties", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenProperties([]*era.ProfileProperties{&era.ProfileProperties{}})
		})
	})

	t.Run("flattenProfilesResponse", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenProfilesResponse(&era.ProfileListResponse{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenProfilesResponse(&era.ProfileListResponse{})
		})
	})

	t.Run("flattenClusterAssociation", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterAssociation([]*era.VersionClusterAssociation{&era.VersionClusterAssociation{}})
		})
	})

	t.Run("flattenAvailableIPsResponse", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenAvailableIPsResponse(&era.GetNetworkAvailableIPs{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenAvailableIPsResponse(&era.GetNetworkAvailableIPs{})
		})
	})

	t.Run("flattenSLAsResponse", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenSLAsResponse(&era.SLAResponse{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSLAsResponse(&era.SLAResponse{})
		})
	})

	t.Run("flattenClonedMetadata", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenClonedMetadata(&era.ClonedMetadata{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClonedMetadata(&era.ClonedMetadata{})
		})
	})

	t.Run("flattenSnapshotsList", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnapshotsList(&era.ListSnapshots{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnapshotsList(&era.ListSnapshots{})
		})
	})

	t.Run("flattenTagsList", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenTagsList(&era.ListTagsResponse{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenTagsList(&era.ListTagsResponse{})
		})
	})

	t.Run("flattenTmsCapability", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenTmsCapability([]*era.Capability{&era.Capability{}})
		})
	})

	t.Run("flattenContinousRegion", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenContinousRegion(&era.ContinuousRegion{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenContinousRegion(&era.ContinuousRegion{})
		})
	})

	t.Run("flattenDBLogs", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBLogs([]*era.DBLogs{&era.DBLogs{}})
		})
	})

	t.Run("flattenDBLogMetadata", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBLogMetadata(&era.DBLogsMetadata{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDBLogMetadata(&era.DBLogsMetadata{})
		})
	})

	t.Run("flattenLastDBLog", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenLastDBLog(&era.DBLogs{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenLastDBLog(&era.DBLogs{})
		})
	})

	t.Run("flattenLastContinousSnapshot", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenLastContinousSnapshot(&era.LastContinuousSnapshot{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenLastContinousSnapshot(&era.LastContinuousSnapshot{})
		})
	})

	t.Run("flattenLastContinousSnapshotMetadata", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenLastContinousSnapshotMetadata(&era.LastContinuousSnapshotMetadata{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenLastContinousSnapshotMetadata(&era.LastContinuousSnapshotMetadata{})
		})
	})

	t.Run("flattenProcessedRanges", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenProcessedRanges([]*era.ProcessedRanges{&era.ProcessedRanges{}})
		})
	})

	t.Run("flattenTimeMachines", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenTimeMachines(&era.ListTimeMachines{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenTimeMachines(&era.ListTimeMachines{})
		})
	})

	t.Run("expandClonesNodes", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClonesNodes(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClonesNodes([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandClonesNodes([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandPostgreSQLCloneActionArgs", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixNDBClone().Schema, map[string]interface{}{})
			_ = expandPostgreSQLCloneActionArgs(d, nil)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixNDBClone().Schema, map[string]interface{}{})
			_ = expandPostgreSQLCloneActionArgs(d, []interface{}{})
		})
	})

	t.Run("expandLCMConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandLCMConfig(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandLCMConfig([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandLCMConfig([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandDBExpiryDetails", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandDBExpiryDetails(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandDBExpiryDetails([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandDBExpiryDetails([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandDBRefreshDetails", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandDBRefreshDetails(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandDBRefreshDetails([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandDBRefreshDetails([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandCredentialInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandCredentialInfo(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandCredentialInfo([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandCredentialInfo([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandNetworkInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNetworkInfo(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNetworkInfo([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandNetworkInfo([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandClusterNetworkInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClusterNetworkInfo(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClusterNetworkInfo([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandClusterNetworkInfo([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandMaintenanceTasks", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandMaintenanceTasks(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandMaintenanceTasks([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandMaintenanceTasks([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandClusterInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClusterInfo(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClusterInfo([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandClusterInfo([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandClusterIPInfos", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClusterIPInfos(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClusterIPInfos([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandClusterIPInfos([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandDBServerPostgresInput", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandDBServerPostgresInput(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandDBServerPostgresInput([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandDBServerPostgresInput([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandPsRegisterDBServer", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandPsRegisterDBServer(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandPsRegisterDBServer([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandPsRegisterDBServer([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("flattenLinkedDBInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenLinkedDBInfo(era.Info{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenLinkedDBInfo(era.Info{})
		})
	})

	t.Run("flattenMaintenanceSchedule", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenMaintenanceSchedule(&era.MaintenaceSchedule{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenMaintenanceSchedule(&era.MaintenaceSchedule{})
		})
	})

	t.Run("expandMaintenanceSchdeule", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandMaintenanceSchdeule(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandMaintenanceSchdeule([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandMaintenanceSchdeule([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("flattenPropertiesMap", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenPropertiesMap(&era.NetworkPropertiesmap{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenPropertiesMap(&era.NetworkPropertiesmap{})
		})
	})

	t.Run("flattenIPPools", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPPools([]*era.IPPools{&era.IPPools{}})
		})
	})

	t.Run("expandSoftwareProfileProp", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandSoftwareProfileProp(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandSoftwareProfileProp([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandSoftwareProfileProp([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandNetworkSingleInstance", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNetworkSingleInstance(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNetworkSingleInstance([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandNetworkSingleInstance([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("flattenStretchedVlans", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenStretchedVlans([]*era.NetworkIntentResponse{&era.NetworkIntentResponse{}})
		})
	})

	t.Run("flattenStretchedVlanMetadata", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenStretchedVlanMetadata(&era.StretchedVlanMetadata{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenStretchedVlanMetadata(&era.StretchedVlanMetadata{})
		})
	})

}
