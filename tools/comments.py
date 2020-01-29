import utils
import time
import re
import requests

old_data = utils.query("select * from comments;")


new_data = list(map(lambda x: {
    "id": int(x["id"]),
    "email": x["email"],
    "recv": x["sendemail"] == "true",
    "avatar": "",
    "time": int(time.mktime(time.strptime(x["time"], "%Y-%m-%d %H:%M:%S"))),
    "raw": re.sub(r'@[0-9]+#', "", x["raw"]),
    "content": requests.post("http://127.0.0.1:50000/api/markdown", {"source": re.sub(r'@[0-9]+#', "", x["raw"])}).json()["html"],
    "url": "/"+("post/" + x["url"][5:].replace("/", "_") if x["url"][0:5] == "post/" else x["url"]).lower(),
    "show": x["show"] == "true",
    "ad": x["ad"] == "true",
    "reply": int(temp[0]) if (temp:= re.findall(r'@([0-9]+)#', x["raw"])) and len(temp) > 0 else -1,
}, old_data))

for d in new_data:
    print(d)

document = utils.mydb["comments"]
document.delete_many({})
ids = document.insert_many(new_data)
