#!/bin/bash
CLIENT_IMAGE=${CLIENT_IMAGE:-secret-server}
docker build . -t $CLIENT_IMAGE
