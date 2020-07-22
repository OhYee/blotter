from bs4 import BeautifulSoup
from utils import Site, Post, get


class iMyShare(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'imyshare.com' in url

    def solver(self, url: str):
        res = get("https://imyshare.com/blog/")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select(".article-item"):
            posts.append(Post(
                item.select_one("h2").get_text(),
                "%s/%s" % (url.strip("/"),
                           item.select_one("a").get("href").strip("/"))
            ))
        return posts


if __name__ == '__main__':
    t = iMyShare()
    print(t.matcher("https://imyshare.com/blog/"))
    print(t.solver("https://imyshare.com/blog/"))
