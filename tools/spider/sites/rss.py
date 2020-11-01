import feedparser
from dateutil.parser import parse
from .utils import Site, Post, get


class RSS(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return True

    def solver(self, url: str):
        res = get(url)
        file = feedparser.parse(res)
        entries = file.entries

        entries.sort(
            key=lambda x: parse(x.published).timestamp(),
            reverse=True,
        )

        posts = []
        for f in entries:
            posts.append(Post(f.title, f.link, parse(f.published).timestamp()))
        return posts


if __name__ == "__main__":
    t = RSS()
    print(t.matcher("http://127.0.0.1/rss.xml"))
    print(t.solver("http://127.0.0.1/rss.xml"))
