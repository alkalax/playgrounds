import * as core from "@actions/core"
import { DefaultAzureCredential } from "@azure/identity"
import { ComputeManagementClient } from "@azure/arm-compute"

async function main() {
  const vmName = core.getInput('name')
  const subId = core.getInput('subscription-id')
  const client = new ComputeManagementClient(new DefaultAzureCredential(), subId)

  const vms = client.virtualMachines.listAll()
  for await (const vm of vms) {
    if (vm.name === vmName) {
      console.log(`Found VM ${vm.name} in ${vm.id}`)
      return
    }
  }

  console.log(`VM ${vmName} not found.`)
}

main().catch((err) => core.setFailed(err.message))
