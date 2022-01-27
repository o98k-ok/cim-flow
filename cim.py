#!/usr/bin/python
import sys
import os
import json
import urllib2
import datetime
import hashlib

indx = sys.argv[1]

url = "http://www.bing.com/HPImageArchive.aspx?format=js&idx=%s&n=1&mkt=en-US" % indx
BaseUrl = "http://www.bing.com"
BasePath = os.environ.get('image_base_path')

if BasePath is None:
    print (BasePath)
    os._exit(-1)

if not os.path.exists(BasePath):
    os.makedirs(BasePath)


# get image url
resp = urllib2.urlopen(url)
data = resp.read()
jsondata = json.loads(data)
imageUrl = BaseUrl + jsondata['images'][0]['url']

# get md5
m = hashlib.md5()
m.update(imageUrl)
code = m.hexdigest()

# for filepath
filename = "{}-{}.jpg".format(code, datetime.date.today().strftime("%Y%m%d"))
filepath = os.path.join(BasePath, filename)

# downloads picture
idata = urllib2.urlopen(imageUrl).read()
fd = open(filepath, "wb")
fd.write(idata)
fd.close()

# parsing files
print(filepath)
