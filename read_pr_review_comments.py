import urllib.request
import json
url = "https://api.github.com/repos/julesklord/tinyfetch/pulls/2/comments"
req = urllib.request.Request(url)
try:
    with urllib.request.urlopen(req) as response:
        data = json.loads(response.read().decode())
        print(json.dumps(data, indent=2))
except Exception as e:
    print(f"Error: {e}")
