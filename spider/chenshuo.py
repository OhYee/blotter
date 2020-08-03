from bs4 import BeautifulSoup

from utils import Site, Post, get, parseToUnix


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
            posts.append(
                Post(
                    item.get_text(),
                    "%s/%s" % (url.strip("/"), item.get("href").strip("/")),
                    parseToUnix(item.parent.parent.select_one(
                        "time").get_text())
                ))
        return posts


if __name__ == '__main__':
    t = Chenshuo()
    print(t.matcher("http://blog.chenshou.top/"))
    print(t.solver("http://blog.chenshou.top/"))
