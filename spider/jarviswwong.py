from bs4 import BeautifulSoup
from utils import Site, Post, get
import json


class Jarviswwong(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'jarviswwong.com' in url

    def solver(self, url: str):
        res = get("https://jarviswwong.com/")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("li.masonry-item"):
            posts.append(Post(
                item.select_one("h1").get_text(),
                item.select_one("a").get("href"),
            ))
        return posts


if __name__ == '__main__':
    t = Jarviswwong()
    print(t.matcher("https://jarviswwong.com/"))
    print(t.solver("https://jarviswwong.com/"))
