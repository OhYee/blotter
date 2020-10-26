import requests
from bs4 import BeautifulSoup
from utils import Site, Post, get, post


class Taifua(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "https://taifua.com/"

    def solver(self, url: str):
        posts = []

        res = get("https://taifua.com/")
        soup = BeautifulSoup(res, features="lxml")
        for item in soup.select(".list-title"):
            link = item.select_one("a")
            posts.append(Post(link.get_text(), link.get("href"), 0))

        res = post(
            "https://taifua.com/wp-admin/admin-ajax.php",
            {
                "append": "list-home",
                "paged": 2,
                "action": "ajax_load_posts",
                "query": "",
                "page": "home",
            },
        )
        soup = BeautifulSoup(res, features="lxml")
        for item in soup.select(".list-title"):
            link = item.select_one("a")
            posts.append(Post(link.get_text(), link.get("href"), 0))

        return posts


if __name__ == "__main__":
    t = Taifua()
    print(t.matcher("https://taifua.com/"))
    print(t.solver("https://taifua.com/"))
