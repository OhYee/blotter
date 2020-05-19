from bs4 import BeautifulSoup
from utils import Site, Post, get


class YLink(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'ylink.top' in url

    def solver(self, url: str):
        res = get("http://ylinknest.top")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select(".post-title"):
            posts.append(Post(
                item.get_text(),
                item.select_one("a").get("href")
            ))
        return posts


if __name__ == '__main__':
    t = YLink()
    print(t.matcher("http://ylink.top/"))
    print(t.solver("http://ylink.top/"))
