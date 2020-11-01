from bs4 import BeautifulSoup

if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


class Lylblog(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'lylblog.com' in url

    def solver(self, url: str):
        res = get("https://www.lylblog.com/")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("article.post"):
            link = item.select_one("a.post-title")
            posts.append(
                Post(
                    link.get_text(),
                    link.get("href"),
                    parseToUnix(item.select_one("time").get_text())
                ))
        return posts


if __name__ == '__main__':
    t = Lylblog()
    print(t.matcher("https://www.lylblog.com/"))
    print(t.solver("https://www.lylblog.com/"))
