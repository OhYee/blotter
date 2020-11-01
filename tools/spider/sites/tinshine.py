from bs4 import BeautifulSoup
from .utils import Site, Post, get


class Tinshine(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'www.tinshine.cn' in url

    def solver(self, url: str):
        res = get(
            "%s/%s" %
            (url.strip("/"), "front/articles.action?pageCount=0")
        )
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select(".card"):
            posts.append(Post(
                item.select_one(".card-header").get_text(),
                "%s/%s" % (url.strip("/"),
                           item.parent.get("href").strip("/"))
            ))
        return posts


if __name__ == '__main__':
    t = Tinshine()
    print(t.matcher("http://www.tinshine.cn/"))
    print(t.solver("http://www.tinshine.cn/"))
