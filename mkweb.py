#!/usr/bin/env python3

import os
import sys
import glob
from pprint import pprint

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
name, source = "sour.is/x/toolbox", "github.com/sour-is/go-toolbox"

subdirs = set()
for filename in glob.iglob('**/*.go', recursive=True):
    subdirs.add(os.path.dirname(filename))

    for d in subdirs:
        lis.append((name + "/" + d, name, source))

for path, name, source in lis:
    os.makedirs(path, exist_ok=True)
    with open(path + "/index.html", "w") as fw:
        print(":: " + name + " - " + path + " ::")
        print(tpl.format(**mk_github(name, source)), file=fw)
