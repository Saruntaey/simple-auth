#!/bin/bash
GO_ENV=development \
PORT=8080 \
MONGO_URI=localhost:27017 \
MONGO_DB=simple-auth \
SESSION_EXPIRED=10 \
./main