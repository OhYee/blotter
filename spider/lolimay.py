import requests
from bs4 import BeautifulSoup
from utils import Site, Post, get


class Lolimay(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "https://lolimay.cn/"

    def solver(self, url: str):
        res = get("https://lolimay.cn/archives/",)
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select(".post-title"):
            link = item.select_one("a")
            posts.append(Post(
                link.get_text(),
                "https://lolimay.cn%s" % link.get("href"),
            ))
        return posts
