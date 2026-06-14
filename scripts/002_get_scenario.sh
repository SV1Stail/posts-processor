#!/usr/bin/env bash

CURL="curl"
CALL="grpcurl --plaintext"
URL="localhost:6090"

function call_GetOriginalPostsFromDB(){
    $CALL -d @ "$URL" post_processor.PostsProcessor/GetOriginalPostsFromDB
}

call_GetOriginalPostsFromDB <<EOF
{
    "limit": 10,
    "offset": 0
}
EOF
