#!/usr/bin/env python3

import os
import sys
import glob
from pprint import pprint

pfx = "sour.is/go/"

tpl = """
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{name} git {git}">
<meta name="go-source" content="{name} {http} {dir} {file}">
<meta http-equiv="refresh" content="0; url={http}">
</head>
<body>
Nothing to see here; <a href="{http}">move along</a>.
</body>
"""

def mk_github(name, source):
    return {
        "name": name,
        "git":  "https://" + source + ".git",
        "http": "https://" + source,
        "dir":  "https://" + source + "/tree/master{/dir}",
        "file": "https://" + source + "/blob/master{/dir}/{file}#L{line}",
    }


lis = []
with open("CANNONICAL_NAMES.txt", 'r') as fd:
    for line in fd.readlines():
        name, source = line.split()
        subdirs = set()
        for filename in glob.iglob(name[len(pfx):] + '/**/*.go', recursive=True):
            subdirs.add(os.path.dirname(filename))

        for d in subdirs:
            lis.append((pfx + d, name, source))

for path, name, source in lis:
    os.makedirs(path, exist_ok=True)
    with open(path + "/index.html", "w") as fw:
        print(":: " + name + " - " + path + " ::")
        print(tpl.format(**mk_github(name, source)), file=fw)
