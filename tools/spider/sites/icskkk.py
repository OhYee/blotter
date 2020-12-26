from bs4 import BeautifulSoup
import re

if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


class ICSKKK(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'icskkk.com' in url

    def solver(self, url: str):
        res = get("http://icskkk.com/archive")
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select(".ant-list-item"):
            posts.append(Post(
                item.select_one("a").get_text(),
                "%s/%s" % (url.strip("/"),
                           item.select_one("a").get("href").strip("/")),
                parseToUnix(re.findall(r'\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}', item.select_one(
                    ".ant-list-item-meta-description"
                ).get_text())[0]),
            ))
        return posts


if __name__ == '__main__':
    t = ICSKKK()
    print(t.matcher("https://icskkk.com/"))
    print(t.solver("https://icskkk.com/"))
