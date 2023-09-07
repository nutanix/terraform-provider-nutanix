variable "user" {
  type = string
}
variable "password" {
  type = string
}
variable "endpoint" {
  type = string
}
variable "insecure" {
  type = bool
}
variable "port" {
  type = number
}

variable "k8s_infra_provision_permissions" {
  type = list(string)
  default = [
    "Create_Category_Mapping",
    "Create_Image",
    "Create_Or_Update_Name_Category",
    "Create_Or_Update_Value_Category",
    "Create_Virtual_Machine",
    "Delete_Category_Mapping",
    "Delete_Image",
    "Delete_Name_Category",
    "Delete_Value_Category",
    "Delete_Virtual_Machine",
    "Update_Category_Mapping",
    "Update_Virtual_Machine_Project",
    "Update_Virtual_Machine",
    "View_Category_Mapping",
    "View_Cluster",
    "View_Image",
    "View_Name_Category",
    "View_Project",
    "View_Subnet",
    "View_Value_Category",
    "View_Virtual_Machine"
  ]
}

variable "csi_system_role_permissions" {
  type = list(string)
  default = [
    "Create_Volume_Group_Disk",
    "Delete_Volume_Group_Disk",
    "Update_Volume_Group_Disk_Internal",
    "View_Project",
    "View_Task",
    "Create_Or_Update_Value_Category",
    "Create_Category",
    "View_Name_Category",
    "View_Category",
    "View_External_iSCSI_Client",
    "View_VM_Recovery_Point",
    "View_Virtual_Machine",
    "View_Volume_Group_Details",
    "View_Volume_Group_Disks",
    "View_Volume_Group_iSCSI_Attachments",
    "View_Volume_Group_VM_Attachments",
    "View_Volume_Group_Category_Associations",
    "View_Volume_Group_Metadata",
    "Create_Virtual_Machine",
    "Restore_VM_Recovery_Point",
    "Delete_Image",
    "Associate_Volume_Group_Categories",
    "Disassociate_Volume_Group_Categories",
    "Update_Virtual_Machine_Project",
    "Update_Container_Disks",
    "View_Image",
    "Create_Category_Mapping",
    "Create_Volume_Group",
    "Delete_Category_Mapping",
    "Update_Category_Mapping",
    "View_Category_Mapping",
    "View_Subnet",
    "Delete_Availability_Zone",
    "Create_Or_Update_Name_Category",
    "Delete_Volume_Group",
    "View_Cluster",
    "View_Value_Category",
    "Delete_Category",
    "Create_Image",
    "Delete_Virtual_Machine",
    "View_Container",
    "View_Storage_Container",
    "View_Any_Virtual_Machine",
    "Create_Job",
    "Update_Virtual_Machine",
    "Update_Network_Function_Chain",
    "Delete_Name_Category",
    "Create_Vm_Snapshot",
    "Update_Account",
    "Delete_Value_Category",
    "Update_Category",
    "Update_Remote_Connection",
    "Attach_Volume_Group_To_External_iSCSI_Client",
    "Detach_Volume_Group_From_External_iSCSI_Client",
    "Create_Consistency_Group",
    "Update_Consistency_Group",
    "View_Consistency_Group",
    "Create_Recovery_Point",
    "View_Recovery_Point",
    "Delete_Recovery_Point",
    "Set_Expiration_Time_Recovery_Point",
    "View_Container_Datastore",
    "View_Container_Stats",
    "Update_Volume_Group_Details_Internal",
    "Update_External_iSCSI_Client_Internal"
  ]
}

variable "k8s_data_services_system_role_permissions" {
  type = list(string)
  default = [
    "Create_Volume_Group_Disk",
    "Delete_Volume_Group_Disk",
    "Update_Volume_Group_Disk_Internal",
    "View_Project",
    "View_Task",
    "Create_Or_Update_Value_Category",
    "Create_Category",
    "View_Name_Category",
    "View_Category",
    "View_External_iSCSI_Client",
    "View_VM_Recovery_Point",
    "View_Virtual_Machine",
    "View_Volume_Group_Details",
    "View_Volume_Group_Disks",
    "View_Volume_Group_iSCSI_Attachments",
    "View_Volume_Group_VM_Attachments",
    "View_Volume_Group_Category_Associations",
    "View_Volume_Group_Metadata",
    "Create_Virtual_Machine",
    "Restore_VM_Recovery_Point",
    "Delete_Image",
    "Associate_Volume_Group_Categories",
    "Disassociate_Volume_Group_Categories",
    "Update_Virtual_Machine_Project",
    "Update_Container_Disks",
    "View_Image",
    "Create_Category_Mapping",
    "Create_Volume_Group",
    "Delete_Category_Mapping",
    "Update_Category_Mapping",
    "View_Category_Mapping",
    "View_Subnet",
    "Delete_Availability_Zone",
    "Create_Or_Update_Name_Category",
    "Delete_Volume_Group",
    "View_Cluster",
    "View_Value_Category",
    "Delete_Category",
    "Create_Image",
    "Delete_Virtual_Machine",
    "View_Container",
    "View_Storage_Container",
    "View_Any_Virtual_Machine",
    "Create_Job",
    "Update_Virtual_Machine",
    "Update_Network_Function_Chain",
    "Delete_Name_Category",
    "Create_Vm_Snapshot",
    "Update_Account",
    "Delete_Value_Category",
    "Update_Category",
    "Update_Remote_Connection",
    "Attach_Volume_Group_To_External_iSCSI_Client",
    "Detach_Volume_Group_From_External_iSCSI_Client",
    "Create_Consistency_Group",
    "Update_Consistency_Group",
    "View_Consistency_Group",
    "Create_Recovery_Point",
    "View_Recovery_Point",
    "Delete_Recovery_Point",
    "Set_Expiration_Time_Recovery_Point",
    "View_Container_Datastore",
    "View_Container_Stats",
    "Update_Volume_Group_Details_Internal",
    "Update_External_iSCSI_Client_Internal"
  ]
}
