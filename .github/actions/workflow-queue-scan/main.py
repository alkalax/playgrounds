import os
import sys
import github

API_TOKEN = os.getenv("GITHUB_TOKEN")
if not API_TOKEN:
    print("Invalid token")
    sys.exit(1)

labels = os.getenv("LABELS")
if not labels:
  print("Labels missing")
  sys.exit(1)

print(f"Labels: {labels}")

auth = github.Auth.Token(API_TOKEN)
gh = github.Github(auth=auth)

repo = gh.get_repo("alkalax/playgrounds")

runners = repo.get_self_hosted_runners()

matching_runners = []
for run in repo.get_workflow_runs():
    if run.status != "completed":
        print(f"Workflow id: {run.workflow_id}")
        print(f"Status: {run.status}")

        print("Fetching jobs for the workflow...")
        for job in run.jobs():
            if job.labels:
                print(f"Labels: {job.labels}")
                for runner in runners:
                    print(f"Runner: {runner.name} {runner.status}")
                    print(f"Labels: {[label['name'] for label in runner.labels]}")

                    if set(job.labels).issubset({label['name'] for label in runner.labels}):
                        print(f"Adding runner {runner.name}")
                        matching_runners.append(runner.name)
            if job.runner_name:
                print(f"Runner name: {job.runner_name}")
            if job.runner_group_name:
                print(f"Runner group name: {job.runner_group_name}")

print(f"Matching runners: {matching_runners}")

output_file = os.getenv("GITHUB_OUTPUT")

if output_file:
    with open(output_file, "a") as f:
        f.write(f"runners={','.join(matching_runners)}\n")
else:
    print("Couldn't open output file")
