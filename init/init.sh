#!/bin/bash

echo "Creating mongodb indexes......"
mongo fse --eval "db.fileInfo.createIndex({'FileId': 1}, {'unique': true});"

