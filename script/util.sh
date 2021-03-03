#!/usr/bin/env bash

function rm_file() {
    ls -- "$1" >/dev/null 2>&1 && rm -- "$1"
}
