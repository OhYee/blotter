from bs4 import BeautifulSoup
from utils import Site, Post, get
import json


class JSPang(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'jspang.com' in url

    def solver(self, url: str):
        res = get("https://jspang.com/")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("div.list-title"):
            link = item.select_one("a")
            posts.append(Post(
                link.get_text(),
                link.get("href"),
            ))
        return posts


if __name__ == '__main__':
    t = JSPang()
    print(t.matcher("https://jspang.com/"))
    print(t.solver("https://jspang.com/"))
