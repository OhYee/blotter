import requests
from bs4 import BeautifulSoup
from utils import Site, Post, get


class Taifua(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "https://taifua.com/"

    def solver(self, url: str):
        res = get("https://taifua.com/")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select(".list-title"):
            link = item.select_one("a")
            posts.append(Post(link.get_text(), link.get("href")))
        return posts


if __name__ == "__main__":
    t = Taifua()
    print(t.matcher("https://taifua.com/"))
    print(t.solver("https://taifua.com/"))
