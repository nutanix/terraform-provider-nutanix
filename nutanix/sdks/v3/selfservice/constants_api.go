package selfservice

const (
	// Blueprint API
	listBlueprintAPI             = "/blueprints/list"
	launchBlueprintAPI           = "/blueprints/%s/simple_launch"       // uses Blueprint UUID
	getBlueprintAPI              = "/blueprints/%s"                     // uses Blueprint UUID
	pendingLaunchBlueprintAPI    = "/blueprints/%s/pending_launches/%s" // uses Blueprint UUID, Launch UUID
	getBlueprintRuntimeEditables = "/blueprints/%s/runtime_editables"   // uses Blueprint UUID

	// Application API
	listApplicationAPI            = "/apps/list"
	getApplicationAPI             = "/apps/%s"                                                                  // uses App UUID
	softDeleteApplicationAPI      = getApplicationAPI + "?type=soft"                                            // uses App UUID
	runApplicationSystemActionAPI = "/apps/%s/actions/run"                                                      // uses App UUID
	runApplicationCustomActionAPI = "/apps/%s/actions/%s/run"                                                   // uses App UUID
	getAppRunlogOutputAPI         = "/apps/%s/app_runlogs/%s/output"                                            // uses App UUID, Runlog UUID
	runPatchActionAPI             = "/apps/%s/patch/%s/run"                                                     // uses App UUID, Patch Runlog UUID
	listAppProtectionPolicyAPI    = "/blueprints/%s/app_profile/%s/config_spec/%s/app_protection_policies/list" // uses Blueprint UUID, App UUID, Config UUID
	listAppRecoveryPointsAPI      = "/apps/%s/recovery_groups/list"                                             // uses App UUID
	deleteRecoveryPointsAPI       = "/apps/%s/recovery_group_delete"                                            // uses App UUID

	// Accounts API
	listAccountsAPI = "/accounts/list"
)
