import * as core from "@actions/core"
import { DefaultAzureCredential } from "@azure/identity"
import { ComputeManagementClient } from "@azure/arm-compute"

try {
  const vmName = core.getInput('name')
  const subId = core.getInput('subscription-id')
  const client = new ComputeManagementClient(new DefaultAzureCredential(), subId)

  console.log("Passed:", vmName)
} catch (err) {
  core.setFailed(err.message)
}
