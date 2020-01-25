import utils

old_data = utils.query("select * from pages;")

new_data = list(map(lambda x: {
    "icon": "",
    "name": x["title"],
    "link": x["url"],
    "index": int(x["idx"]),
}, old_data))

for d in new_data:
    print(d)

document = utils.mydb["pages"]

ids = document.insert_many(new_data)
