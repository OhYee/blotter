from rss import RSS
from utils import Site, Post, get


class CSDN(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url[:22] == "https://blog.csdn.net/"

    def solver(self, url: str):
        return RSS().solver("%s/rss/list" % (url.strip("/")))
