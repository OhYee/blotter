import pymongo
import sqlite3

def query(query_string):
    with sqlite3.connect("./database.db") as db:
        cur = db.execute(query_string)
        res = cur.fetchall()

    return [
        {
            name[0]: data[idx].replace(r"$double-quote;", r'"')
            for idx, name in enumerate(cur.description)
        }
        for data in res
    ]

myclient = pymongo.MongoClient("mongodb://localhost:27017/")
mydb = myclient["blotter"]