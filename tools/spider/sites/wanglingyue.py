from bs4 import BeautifulSoup
import json

if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


class Wanglingyue(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'wanglingyue.com' in url

    def solver(self, url: str):
        res = post(
            "https://www.wanglingyue.com/api/article/list?page=1&limit=10&orderBy=undefined",
            "{}",
            headers={'Content-Type': 'application/json;charset=UTF-8'},
        )
        posts = []
        for p in json.loads(res)["data"]["rows"]:
            posts.append(Post(
                p["title"],
                "https://www.wanglingyue.com/archives/%d" % p["id"],
                parseToUnix(p["createTime"])
            ))
        return posts


if __name__ == '__main__':
    t = Wanglingyue()
    print(t.matcher("https://www.wanglingyue.com"))
    print(t.solver("https://www.wanglingyue.com"))
