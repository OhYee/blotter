import requests
import pymongo
import threading
import datetime

from sites import sites

threadLock = threading.Lock()
myclient = pymongo.MongoClient("mongodb://localhost:27017/")
mydb = myclient["blotter"]
document = mydb["friends"]


class Worker(threading.Thread):
    def __init__(self, fid, rss):
        # super()
        threading.Thread.__init__(self)
        self.fid = fid
        self.rss = rss

    def run(self):
        print("Start", self.rss)

        posts = []
        try:
            posts = getSitePosts(self.rss)
            posts.sort(key=lambda x: x.time, reverse=True)
            if len(posts) > 5:
                posts = posts[:5]
        except Exception as e:
            err("%s %s" % (self.rss, str(e)))

        threadLock.acquire()
        document.update_one(
            {"_id": self.fid},
            {
                "$set": {
                    "error": True,
                } if len(posts) == 0 else {
                    "error": False,
                    "posts": [{"title": post.title, "link": post.link, "time": post.time} for post in posts]
                }
            }
        )
        threadLock.release()
        print("Finish", self.rss)


def getSitePosts(url: str):
    for site in sites:
        if site.matcher(url):
            return site.solver(url)
    return []


def err(e: str):
    print(e)


if __name__ == '__main__':
    print(datetime.datetime.now())
    threads = [
        Worker(item["_id"], item["rss"])
        for item in document.find({})
        if item["rss"] != ""
    ]

    for t in threads:
        t.start()
    for t in threads:
        t.join()
