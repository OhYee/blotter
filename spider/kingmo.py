from bs4 import BeautifulSoup
from utils import Site, Post, get


class Kingmo(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "https://www.lizenghai.com/"

    def solver(self, url: str):
        res = get("https://www.lizenghai.com/user/1/posts")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.find_all("a", rel="bookmark"):
            posts.append(Post(item.get_text(), item.get("href")))
        return posts
