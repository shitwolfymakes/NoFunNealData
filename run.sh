#!/bin/bash

go run main.go

git add export.json
git commit -m "updated data"
git push
