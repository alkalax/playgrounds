import * as core from "@actions/core";

try {
  const nameToGreet = core.getInput('who-to-greet');
  console.log(`Hello ${nameToGreet}!`);
  console.log("Extra log")
} catch (err) {
  core.setFailed(err.message);
}
