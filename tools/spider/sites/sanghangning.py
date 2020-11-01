from bs4 import BeautifulSoup
if __name__ == "__main__":
    from utils import *
else:
    from .utils import *
import json


class Sanghangning(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return 'sanghangning.cn' in url

    def solver(self, url: str):
        res = get("https://www.sanghangning.cn/json/blog.json")
        data = json.loads(res)
        posts = []
        for post in data["blog"]:
            posts.append(Post(
                post['title'],
                "%s/%s" % (
                    "https://www.sanghangning.cn".strip("/"),
                    post['url'].strip("/")
                ),
                parseToUnix(post["date"]),
            ))
        return posts


if __name__ == '__main__':
    t = Sanghangning()
    print(t.matcher("https://www.sanghangning.cn/"))
    print(t.solver("https://www.sanghangning.cn/"))
