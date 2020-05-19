
import requests
import urllib3

urllib3.disable_warnings()

headers = {
    'User-Agent': 'OhYee Spider',
}


def get(url: str):
    rep = requests.get(url, verify=False, headers=headers)
    rep.encoding = 'utf-8'
    return rep.text


def post(url: str, data: object):
    rep = requests.post(url, data, verify=False, headers=headers)
    rep.encoding = 'utf-8'
    return rep.text


class Site:
    def __init__(self):
        pass

    def matcher(self, url: str):
        return False

    def solver(self, url: str):
        return []


class Post:
    def __init__(self, title: str, link: str):
        self.title = title.strip()
        self.link = link

    def __repr__(self):
        return "%s - %s" % (self.title, self.link)
