import boto3
import datetime
import os

from dateutil.tz import tzutc


def check_conditions(stack: dict, target_stacks: dict, time_diff: datetime.timedelta):
    current_time = datetime.datetime.now(tz=tzutc())

    for k in target_stacks:
        if k in stack["StackName"] and current_time - stack["CreationTime"] > time_diff:
            target_stacks[k].append(stack["StackName"])


if __name__ == "__main__":
    client = boto3.client("cloudformation")
    response = client.describe_stacks()
    difference = datetime.timedelta(hours=6)

    targetStacks = {'nightly':[], 'e2e-gpu':[], 'parallel':[]}

    for each in response["Stacks"]:
        check_conditions(each, targetStacks, difference)

    for k, lists in targetStacks.items():
        for each in lists:
            print("det-deploy aws down --cluster-id "+each)
            # os.system("det-deploy aws down --cluster-id " + each)
