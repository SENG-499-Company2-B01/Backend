#!/bin/bash 
docker build . -t mongodb 
docker run -d mongodb