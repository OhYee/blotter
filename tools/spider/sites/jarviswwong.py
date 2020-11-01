import datetime
import re
import json
from bs4 import BeautifulSoup

if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


class Jarviswwong(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'jarviswwong.com' in url

    def solver(self, url: str):
        res = get("https://jarviswwong.com/")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select(".ajaxcard"):
            y, m, d = map(int, re.findall(
                r'(\d+)', item.select_one(".info-date").select_one("span").get_text()))
            posts.append(Post(
                item.select_one(".card-title").get_text(),
                item.select_one("a").get("href"),
                datetime.datetime(y, m, d).timestamp(),
            ))
        return posts


if __name__ == '__main__':
    t = Jarviswwong()
    print(t.matcher("https://jarviswwong.com/"))
    print(t.solver("https://jarviswwong.com/"))
