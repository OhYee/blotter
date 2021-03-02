from bs4 import BeautifulSoup
import datetime

if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


class SummerFire(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return "summerfire.cn" in url

    def solver(self, url: str):
        res = get(url)
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("article.summary"):
            title = item.select_one("h1.single-title")
            y, m, d = map(int, item.select_one(
                "time").get_text().split("-"))
            posts.append(
                Post(
                    title.get_text(),
                    "%s/%s" % (url.strip("/"),
                               title.select_one("a").get("href").strip("/")),
                    datetime.datetime(y, m, d).timestamp(),
                ))
        return posts


if __name__ == '__main__':
    t = SummerFire()
    print(t.matcher("https://www.summerfire.cn/"))
    print(t.solver("https://www.summerfire.cn/"))
