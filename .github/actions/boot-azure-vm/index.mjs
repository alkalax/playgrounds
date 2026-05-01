import * as core from "@actions/core"
import { DefaultAzureCredential } from "@azure/identity"
import { ComputeManagementClient } from "@azure/arm-compute"

try {
  const vmName = core.getInput('name')
  const subId = core.getInput('subscription-id')
  const client = new azCompute.ComputeManagementClient(new DefaultAzureCredential(), subId)

  console.log("Passed.")
} catch (err) {
  core.setFailed(err.message)
}
