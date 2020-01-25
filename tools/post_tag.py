import utils


posts = {
    post["url"].lower(): post["tags"].split(",")
    for post in utils.query("select * from posts;")
}


document = utils.mydb["tags"]
tags = {data["name"]: data["_id"] for data in document.find({})}

document = utils.mydb["posts"]

data = []
for post in document.find({}):
    for tag in posts[post["url"]]:
        data.append({
            "tag": tags[tag],
            "post": post["_id"],
        })

for d in data:
    print(d)

document = utils.mydb["post_tag"]
ids = document.insert_many(data)
