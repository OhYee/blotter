import utils
import time

document = utils.mydb["tags"]
tags = {tag["name"]: tag["_id"] for tag in document.find({})}

old_data = utils.query("select * from posts;")


new_data = list(map(lambda x: {
    "title": x["title"],
    "abstract": x["abstruct"],
    "view": int(x["view"]),
    "url": x["url"].lower(),
    "publish_time": int(time.mktime(time.strptime(x["time"], "%Y-%m-%d %H:%M:%S"))),
    "edit_time": int(time.mktime(time.strptime(x["updatetime"], "%Y-%m-%d %H:%M:%S"))),
    "content": x["html"],
    "raw": x["raw"],
    "tags": [tags[tag] for tag in x["tags"].split(",")],
    "keywords": [],
    "published": x["published"] == "true",
    "head_image": x["img"],
}, old_data))


document = utils.mydb["posts"]
document.delete_many({})
ids = document.insert_many(new_data)
