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

def load_virtual_machines(credential, check_status=False):

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

    if check_status:
        for vm_name, _ in virtual_machines.items():
            virtual_machines[vm_name]["status"] = get_virtual_machine_status(vm_name, credential)

def get_virtual_machine_status(vm_name, credential):
    if vm_name not in virtual_machines:
        print(f"VM '{vm_name}' not found.")
        return None

    compute_client = ComputeManagementClient(
        credential=credential, 
        subscription_id=virtual_machines[vm_name]["subscription_id"]
    )
    resource_group = virtual_machines[vm_name]["resource_group"]
    vm_status = compute_client.virtual_machines.get(resource_group, vm_name, expand='instanceView')
    
    for status in vm_status.instance_view.statuses:
        if status.code.startswith('PowerState/'):
            return status.code.split('/')[1]
    
    return "unknown"

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

    return map(
        lambda e: {
            "timestamp": e.event_timestamp,
            "action": e.operation_name.value.split("/")[-2] if e.operation_name and e.operation_name.value else None
        },
        filter(
            lambda e: e.operation_name and e.operation_name.value in actions and e.status.value == "Succeeded",
            events
        )
    )

if __name__ == "__main__":
    credential = DefaultAzureCredential()
    load_virtual_machines(credential=credential, check_status=True)

    events = fetch_activity_logs(vm_name="ubuntu-0", credential=credential)
    for event in events:
        print(event)

    for vm_name, vm_info in virtual_machines.items():
        print(f"VM Name: {vm_name}, Status: {vm_info.get('status')}")