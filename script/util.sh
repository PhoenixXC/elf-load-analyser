#!/usr/bin/env bash

function rm_suffix() {
    ls -- *."$1" >/dev/null 2>&1 && rm -- *."$1"
}
