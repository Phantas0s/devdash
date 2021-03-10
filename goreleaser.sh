#!/bin/bash

goreleaser && rm -rf dist/ && cd -
