#!/bin/bash

go run main.go

git add export.csv
git commit -m "updated data"
git push
