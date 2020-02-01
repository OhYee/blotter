import utils
import time
import re
import requests
import bson

data = sorted(utils.query("select * from comments;"),
              key=lambda x: int(x["id"]))


m = {}

defaultObjectID = bson.ObjectId("000000000000000000000000")


def transfer(x):
    temp = re.findall(r'@([0-9]+)#', x["raw"])
    url = ""
    if x["url"][0:5] == "post/":
        url = "/post/" + x["url"][5:].replace("/", "_")
    elif x["url"] == "pages/comments":
        url = "/comment"
    else:
        url = "/"+x["url"]
    url = url.lower()

    return {
        # "id": int(x["id"]),
        "email": x["email"],
        "recv": x["sendemail"] == "true",
        "avatar": requests.get("http://127.0.0.1:50000/api/avatar", {"email": x["email"]}).json()["avatar"],
        "time": int(time.mktime(time.strptime(x["time"], "%Y-%m-%d %H:%M:%S"))),
        "raw": re.sub(r'@[0-9]+#', "", x["raw"]),
        "content": requests.post("http://127.0.0.1:50000/api/markdown", {"source": re.sub(r'@[0-9]+#', "", x["raw"])}).json()["html"],
        "url": url,
        "show": x["show"] == "true",
        "ad": x["ad"] == "true",
        "reply": m.get(int(temp[0]), defaultObjectID) if len(temp) > 0 else defaultObjectID,
    }


document = utils.mydb["comments"]
document.delete_many({})


for i, d in enumerate(data):
    print(i, d)
    objId = document.insert_one(transfer(d))
    m[i+1] = objId.inserted_id
    print(type(objId.inserted_id))
