import * as core from "@actions/core"
import { DefaultAzureCredential } from "@azure/identity"
import { ComputeManagementClient } from "@azure/arm-compute"

const subId = 'a2b28c85-1948-4263-90ca-bade2bac4df4';

try {
  const client = new azCompute.ComputeManagementClient(new DefaultAzureCredential(), subId)
} catch (err) {
  core.setFailed(err.message)
}
