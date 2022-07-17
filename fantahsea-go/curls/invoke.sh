#!/bin/bash

curl -X POST http://localhost:8081/fantahsea/open/gallery/list -H "id: 1" -H "username: zhuangyongj" -H "userno: 123" -H "role: admin" -d '{ "pagingVo" : { "limit":10, "page":1} }' -v; echo

curl -X POST http://localhost:8081/fantahsea/open/gallery/images -H "id: 1" -H "username: zhuangyongj" -H "userno: 123" -H "role: admin" -d '{ "pagingVo" : { "limit":10, "page":1}, "galleryNo" : "GAL111233" }' -v; echo

curl -X POST "http://localhost:8081/fantahsea/open/gallery/image/transfer" -H "id: 3" -H "username: zhuangyongj" -H "userno: 123" -H "role: admin" -d '{ "galleryNo" : "GAL111233", "name" : "yoyoyo", "fileKey" : "eb6bc04f-15c5-4f85-a84d-be3d5a7236d8" }' -v; echo

curl -X GET "http://localhost:8081/fantahsea/open/gallery/image/download?imageNo=IMG5FBBD6P84OE7UT4QT" -H "id: 1" -H "username: zhuangyongj" -H "userno: 123" -H "role: admin" -v






