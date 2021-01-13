if __name__ == "__main__":
    from utils import *
    from rss import RSS
else:
    from .utils import *
    from .rss import RSS


class CSDN(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url[:22] == "https://blog.csdn.net/"

    def solver(self, url: str):
        return RSS().solver("%s/rss/list" % (url.strip("/")))


if __name__ == '__main__':
    t = CSDN()
    print(t.matcher("https://blog.csdn.net/qq_42673093"))
    print(t.solver("https://blog.csdn.net/qq_42673093"))
