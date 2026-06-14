#!/usr/bin/env bash

# set -x

CURL="curl"
CALL="grpcurl --plaintext"
URL="localhost:6090"

function call_CreateOriginalPost() {
    $CALL -d @ "$URL" post_processor.PostsProcessor/CreateOriginalPost
}

# function call_CreateOriginalPost_2() {
#     $CALL -d '{
#     "text": "Сегодня был прекрасный закат в горах. Природа напоминает, как важно иногда замедляться и наслаждаться моментом. 🏔️",
#     "link_original_post": "https://example.com/post/12345",
#     "original_channel": "telegram",
#     "theme": "test"
#     }' "$URL" post_processor.PostsProcessor/CreateOriginalPost
# }

function call_GetOriginalPostsFromDB(){
    $CALL -d @ "$URL" post_processor.PostsProcessor/GetOriginalPostsFromDB
}
