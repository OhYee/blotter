from datetime import datetime
import requests
import re
from bs4 import BeautifulSoup
from utils import Site, Post, get, parseToUnix

regxp = re.compile("(\d+)")


class ShawnLuo(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "https://www.shawnluo.com/"

    def solver(self, url: str):
        res = get(url)
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("article.kratos-hentry"):
            link = item.select_one(".kratos-entry-title-new").select_one("a")

            timeTuple = list(map(int, regxp.findall(
                item.select_one("i.fa-calendar").parent.get_text())))

            posts.append(Post(
                link.get_text(), link.get("href"),
                datetime(timeTuple[0], timeTuple[1], timeTuple[2]).timestamp(),
            ))
        return posts


if __name__ == '__main__':
    t = ShawnLuo()
    print(t.matcher("https://www.shawnluo.com/"))
    print(t.solver("https://www.shawnluo.com/"))
