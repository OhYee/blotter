import utils
import time
import requests

document = utils.mydb["tags"]
tags = {tag["name"]: tag["_id"] for tag in document.find({})}

old_data = utils.query("select * from posts;")


new_data = list(map(lambda x: {
    "title": x["title"],
    "abstract": x["abstruct"],
    "view": int(x["view"]),
    "url": x["url"].lower().replace("/", "_"),
    "publish_time": int(time.mktime(time.strptime(x["time"], "%Y-%m-%d %H:%M:%S"))),
    "edit_time": int(time.mktime(time.strptime(x["updatetime"], "%Y-%m-%d %H:%M:%S"))),
    "content":  requests.post("http://127.0.0.1:50000/api/markdown", {"source": x["raw"]}).json()["html"],
    "raw": x["raw"],
    "tags": [tags[tag] for tag in x["tags"].split(",")],
    "keywords": [],
    "published": x["published"] == "true",
    "head_image": x["img"],
}, old_data))


document = utils.mydb["posts"]
document.delete_many({})
ids = document.insert_many(new_data)
