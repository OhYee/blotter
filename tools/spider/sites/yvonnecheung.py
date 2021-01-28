from bs4 import BeautifulSoup
import re

if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


class Yvonnecheung(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'yvonnecheung.cn' in url

    def solver(self, url: str):
        res = get(url)
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("div.abstract-item"):
            posts.append(Post(
                item.select_one("a").get_text(),
                "%s/%s" % (url.strip("/"),
                           item.select_one("a").get("href").strip("/")),
                parseToUnix(item.select_one("i.date").get_text().replace("Spt", "Sep")),
            ))
        return posts


if __name__ == '__main__':
    t = Yvonnecheung()
    print(t.matcher("https://www.yvonnecheung.cn/"))
    print(t.solver("https://www.yvonnecheung.cn/"))
