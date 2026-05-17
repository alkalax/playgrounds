from azure.identity import DefaultAzureCredential
from azure.mgmt.resource.subscriptions import SubscriptionClient
from azure.mgmt.compute import ComputeManagementClient

virtual_machines = {}

def load_virtual_machines():
    credential = DefaultAzureCredential()

    sub_client = SubscriptionClient(credential=credential)

    for sub in list(sub_client.subscriptions.list()):
        try:
            if not sub.subscription_id:
                continue
            compute_client = ComputeManagementClient(credential=credential, subscription_id=sub.subscription_id)

            for vm in compute_client.virtual_machines.list_all():
                virtual_machines[vm.name] = {
                    "subscription_id": sub.subscription_id,
                    "resource_group": vm.id.split("/")[4] if vm.id else ""
                }
        except Exception as e:
            print(f"Failed to read compute resources: {e}")


if __name__ == "__main__":
    load_virtual_machines()

    print(virtual_machines)
