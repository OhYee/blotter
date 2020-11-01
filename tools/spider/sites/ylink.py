from datetime import datetime
from bs4 import BeautifulSoup
if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


class YLink(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'ylinknest.top' in url

    def solver(self, url: str):
        res = get("http://ylinknest.top")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select(".post-title"):
            timeTuple = list(map(int, item.parent.select_one(
                ".fa-calendar").parent.get_text().split("/")))

            posts.append(Post(
                item.get_text(),
                item.select_one("a").get("href"),
                datetime(2000+timeTuple[0], timeTuple[1],
                         timeTuple[2]).timestamp(),
            ))
        return posts


if __name__ == '__main__':
    t = YLink()
    print(t.matcher("http://ylinknest.top/"))
    print(t.solver("http://ylinknest.top/"))
