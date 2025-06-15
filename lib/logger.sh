#!/usr/bin/env bash

log() {
    local level=$1; shift
    local message="$*"
    local color=$COLOR_NC
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    case "$level" in
        "INFO") color=$COLOR_CYAN ;;
        "SUCCESS") color=$COLOR_GREEN ;;
        "WARN") color=$COLOR_YELLOW ;;
        "ERROR") color=$COLOR_RED ;;
        "DEBUG") color=$COLOR_PURPLE ;;
    esac
    
    if [ "$level" == "DEBUG" ] && [ "$DEBUG" != "true" ]; then return; fi
    
    echo -e "${color}[$level]${COLOR_NC} $message" >&2
    
    if [ "$DEBUG" == "true" ]; then
        echo "[$timestamp] [$level] $message" >> "$MIGRATOR_HOME/migrator.log"
    fi
}