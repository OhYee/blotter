from bs4 import BeautifulSoup
from utils import Site, Post, get


class Chenshuo(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "http://blog.chenshou.top/"

    def solver(self, url: str):
        res = get(url)
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("a.article-title"):
            posts.append(Post(item.get_text(), item.get("href")))
        return posts
