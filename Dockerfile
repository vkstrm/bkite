# Base iamge for pipeline
# Should contain all required dependencies for building, testing, etc.

FROM golang:1.26-alpine

RUN go install gotest.tools/gotestsum@latest
