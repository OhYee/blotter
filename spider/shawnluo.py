import requests
from bs4 import BeautifulSoup
from utils import Site, Post, get


class ShawnLuo(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "https://www.shawnluo.com/"

    def solver(self, url: str):
        res = get(url)
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("h2.kratos-entry-title-new"):
            link = item.select_one("a")
            posts.append(Post(link.get_text(), link.get("href")))
        return posts
