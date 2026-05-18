from azure.identity import DefaultAzureCredential
from azure.mgmt.resource.subscriptions import SubscriptionClient
from azure.mgmt.compute import ComputeManagementClient
from azure.mgmt.monitor import MonitorManagementClient
from datetime import datetime, timedelta

virtual_machines = {}
actions = [
    'Microsoft.Compute/virtualMachines/start/action',
    'Microsoft.Compute/virtualMachines/deallocate/action'
]

def load_virtual_machines(credential):

    sub_client = SubscriptionClient(credential=credential)

    for sub in list(sub_client.subscriptions.list()):
        try:
            if not sub.subscription_id:
                continue
            compute_client = ComputeManagementClient(credential=credential, subscription_id=sub.subscription_id)

            for vm in compute_client.virtual_machines.list_all():
                virtual_machines[vm.name] = {
                    "resource_id": vm.id,
                    "subscription_id": sub.subscription_id,
                    "resource_group": vm.id.split("/")[4] if vm.id else None
                }
        except Exception as e:
            print(f"Failed to read compute resources: {e}")

def fetch_activity_logs(vm_name, credential):
    monitor_client = MonitorManagementClient(
        credential=credential, 
        subscription_id=virtual_machines[vm_name]["subscription_id"]
    )

    start = datetime.now() - timedelta(hours=24)
    end = datetime.now()
    filter_str = (
        f"resourceId eq '{virtual_machines[vm_name]['resource_id']}' and "
        f"eventTimestamp ge '{start.isoformat()}' and "
        f"eventTimestamp le '{end.isoformat()}'"
    )
    events = monitor_client.activity_logs.list(filter=filter_str)

    for event in events:
        if (
            event.operation_name 
            and event.operation_name.value in actions
            and event.status.value == "Succeeded"
        ):
            print(f"Time: {event.event_timestamp}")
            print(f"Event: {event.operation_name.value}")
            print(f"Event: {event.operation_name.localized_value}")
            print(f"Status: {event.status.value}")

if __name__ == "__main__":
    credential = DefaultAzureCredential()
    load_virtual_machines(credential=credential)

    fetch_activity_logs(vm_name="ubuntu-0", credential=credential)
