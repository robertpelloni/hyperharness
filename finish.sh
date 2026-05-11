#!/bin/bash
git rm -r --cached .
git add .
git commit -m "chore: snapshot tree" || true
