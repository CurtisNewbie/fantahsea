#!/bin/bash

curl -X POST http://localhost:8080/fantahsea/open/gallery/list -H "id: 1" -H "username: zhuangyongj" -H "userno: 123" -H "role: admin" -d '{ "pagingVo" : { "limit":10, "page":1} }' -v; echo
