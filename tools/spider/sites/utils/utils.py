
import requests
import urllib3
import datetime
from dateutil.parser import parse

urllib3.disable_warnings()

_headers = {
    'User-Agent': 'OhYee Spider',
}


def auto_retry(method: str, url: str, timeout=30, verify=False, headers=_headers, **kwargs,):
    retry = 10
    rep = None

    def gotException():
        nonlocal retry
        print("Get page error: {}, {} times left...".format(e, retry))
        retry -= 1

    while retry > 0:
        try:
            rep = requests.request(
                method,
                url,
                timeout=timeout,
                verify=verify,
                headers=headers,
                **kwargs,
            )
            break
        except requests.exceptions.ContentDecodingError as e:
            # 当站点编码方式设置错误时，可能会导致无法正确解码
            try:
                rep = requests.request(
                    url,
                    timeout=timeout,
                    verify=verify,
                    headers={**headers, **{"Accept-Encoding": ""}}
                )
                break
            except Exception as e:
                gotException()
        except Exception as e:
            gotException()

    rep.encoding = 'utf-8'
    return rep.text


def get(url: str, headers={}):
    return auto_retry(
        "get",
        url,
        headers={**_headers, **headers},
    )


def post(url: str, data: object, headers={}):
    return auto_retry(
        "post",
        url,
        data=data,
        headers={**_headers, **headers},
    )


class Site:
    def __init__(self):
        pass

    def matcher(self, url: str):
        return False

    def solver(self, url: str):
        return []


class Post:
    def __init__(self, title: str, link: str, time: int):
        self.title = title.strip()
        self.link = link
        self.time = time

    def __repr__(self):

        return "(%s - %s - %s)" % (
            self.title,
            self.link,
            datetime.datetime.fromtimestamp(
                self.time
            ).strftime("%Y-%m-%d %H:%M:%S"),
        )


def parseToUnix(timeStr: str):
    return parse(timeStr).timestamp()
