from bs4 import BeautifulSoup
if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


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
            posts.append(Post(
                item.get_text(),
                item.get("href"),
                parseToUnix(item.parent.parent.parent.select_one(
                    "time").get("datetime")),
            ))
        return posts


if __name__ == '__main__':
    t = Kingmo()
    print(t.matcher("https://www.lizenghai.com/"))
    print(t.solver("https://www.lizenghai.com/"))
