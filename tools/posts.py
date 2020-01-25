import utils
import time

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
    "keywords": [],
    "published": x["published"] == "true",
    "head_image": x["img"],
}, old_data))

for d in new_data:
    print(d)

document = utils.mydb["posts"]

ids = document.insert_many(new_data)
