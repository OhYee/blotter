import requests
from bs4 import BeautifulSoup
from utils import Site, Post, get, parseToUnix


class OmegaXYZ(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "https://www.omegaxyz.com/"

    def solver(self, url: str):
        res = get(url)
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("h3.rpwe-title"):
            link = item.select_one("a")
            posts.append(Post(
                link.get_text(),
                link.get("href"),
                parseToUnix(item.parent.select_one("time").get("datetime")),
            ))
        return posts


if __name__ == '__main__':
    t = OmegaXYZ()
    print(t.matcher("https://www.omegaxyz.com/"))
    print(t.solver("https://www.omegaxyz.com/"))
