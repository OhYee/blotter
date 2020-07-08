from bs4 import BeautifulSoup
from utils import Site, Post, get
import json


class SecNews(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'wiki.ioin.in' in url

    def solver(self, url: str):
        res = get("https://wiki.ioin.in/")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("tr"):
            print(item)
            link = item.select_one("a")
            if link == None:
                continue
            posts.append(Post(
                link.get_text(),
                "%s/%s" % (url.strip("/"), link.get("href").strip("/"))),
            )
        return posts


if __name__ == '__main__':
    t = SecNews()
    print(t.matcher("https://wiki.ioin.in/"))
    print(t.solver("https://wiki.ioin.in/"))
