import utils

old_data = utils.query("select * from friends;")

new_data = list(map(lambda x: {
    "index": int(x["idx"]),
    "image": '',
    "link": x["url"],
    "name": x["name"],
    "description": '',
    "posts": [],
}, old_data))

for d in new_data:
    print(d)

document = utils.mydb["friends"]

ids = document.insert_many(new_data)
