import requests
import pymongo

from rss import RSS
from taifua import Taifua
from kingmo import Kingmo
from lolimay import Lolimay
from shawnluo import ShawnLuo
from omegaxyz import OmegaXYZ
from chenshuo import Chenshuo
from csdn import CSDN

sites = [
    RSS(),
    Taifua(),
    Kingmo(),
    Lolimay(),
    ShawnLuo(),
    OmegaXYZ(),
    Chenshuo(),
    CSDN(),
]


def getSitePosts(url: str):
    for site in sites:
        if site.matcher(url):
            return site.solver(url)
    return []


def err(e: str):
    print(e)


if __name__ == '__main__':
    myclient = pymongo.MongoClient("mongodb://localhost:27017/")
    mydb = myclient["blotter"]
    document = mydb["friends"]

    for item in document.find({}):
        if item["rss"] != "":
            try:
                posts = getSitePosts(item["rss"])
                if len(posts) == 0:
                    err("No posts in %s" % item['rss'])
                else:
                    document.update_one(
                        {"_id": item["_id"]},
                        {
                            "$set": {
                                "posts": [{"title": post.title, "link": post.link} for post in posts[:3]]
                            }
                        }
                    )

            except Exception as e:
                err("%s %s" % (item["rss"], str(e)))