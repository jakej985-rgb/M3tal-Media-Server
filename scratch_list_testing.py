import re
import json

with open("/home/m3tal/.gemini/antigravity-ide/brain/4a1ef673-af49-4b65-9649-d43e1bbd4b5c/.system_generated/steps/225/content.md", "r") as f:
    content = f.read()

# Find the memex-paginated-items-data script content
match = re.search(r'<script type="application/json" id="memex-paginated-items-data">(.*?)</script>', content, re.DOTALL)
if not match:
    print("Could not find memex-paginated-items-data script tag")
    exit(1)

data = json.loads(match.group(1))

# Find the group value for 'Testing'
testing_group = None
for node in data["groups"]["nodes"]:
    if node["groupValue"] == "Testing":
        testing_group = node
        break

if not testing_group:
    print("Could not find Testing group")
    exit(1)

print(f"Testing Group ID: {testing_group['groupId']}")
print(f"Total Count in Testing: {testing_group['totalCount']['value']}")

# Find all nodes in groupedItems belonging to this group
for group in data["groupedItems"]:
    if group["groupId"] == testing_group["groupId"]:
        print("Issues in Testing:")
        for node in group["nodes"]:
            # Find the title/number
            title = ""
            number = None
            state = node.get("state", "unknown")
            url = ""
            for val in node["memexProjectColumnValues"]:
                if val.get("memexProjectColumnId") == "Title":
                    title_info = val.get("value", {}).get("title", {})
                    title = title_info.get("raw", "")
                    number = val.get("value", {}).get("number")
                    url = val.get("value", {}).get("url", "")
            print(f" - #{number}: {title} (State: {state})")
