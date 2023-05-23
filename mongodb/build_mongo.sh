#!/bin/bash 
docker build . -t mongodb 
docker run --name mongodb -d mongodb