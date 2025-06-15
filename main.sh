#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

: "${MIGRATOR_HOME:="$HOME/.dugong"}"
: "${MIGRATOR_GALLERY_URL:="https://raw.githubusercontent.com/joshm998/dugong-gallery/main"}"
: "${DEBUG:=false}"

# Colors
readonly COLOR_GREEN='\033[0;32m'
readonly COLOR_YELLOW='\033[0;33m'
readonly COLOR_RED='\033[0;31m'
readonly COLOR_CYAN='\033[0;36m'
readonly COLOR_BLUE='\033[0;34m'
readonly COLOR_PURPLE='\033[0;35m'
readonly COLOR_NC='\033[0m'

# Directories
readonly SCRIPTS_DIR="$MIGRATOR_HOME/scripts"
readonly JOURNALS_DIR="$MIGRATOR_HOME/journals"
readonly CACHE_DIR="$MIGRATOR_HOME/cache"

# Core utility functions
# shellcheck source=./lib/core.sh
source "$SCRIPT_DIR/lib/core.sh"

usage() {
    echo "Usage: $0 [command] [options]"
    echo ""
    echo "██████╗ ██╗   ██╗ ██████╗  ██████╗ ███╗   ██╗ ██████╗"
    echo "██╔══██╗██║   ██║██╔════╝ ██╔═══██╗████╗  ██║██╔════╝"
    echo "██║  ██║██║   ██║██║  ███╗██║   ██║██╔██╗ ██║██║  ███╗"
    echo "██║  ██║██║   ██║██║   ██║██║   ██║██║╚██╗██║██║   ██║"
    echo "██████╔╝╚██████╔╝╚██████╔╝╚██████╔╝██║ ╚████║╚██████╔╝"
    echo "╚═════╝  ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝  "
    echo ""
    echo "The missing package manager for bash scripts"
    echo "Commands:"
    echo "  apply [id]                 - Apply unapplied scripts"
    echo "    --local <dir>            - Apply scripts from a local project"
    echo "    --tag <tag>              - Apply only scripts with specified tag"
    echo "    --all                    - Apply all unapplied scripts (default)"
    echo ""
    echo "  remove [id]                - Remove applied scripts"
    echo "    --tag <tag>              - Remove all scripts with specified tag"
    echo "    --all                    - Remove all applied scripts"
    echo ""
    echo "  info [id]                  - Show detailed information about a script from the gallery"
    echo ""
    echo "  list [options]             - List all available journals"
    echo ""
    echo "  help                       - Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  MIGRATOR_HOME:             Home directory for migrator data (default: ~/.migrator)"
    echo "  MIGRATOR_GALLERY_URL:      Base URL for script gallery"
    echo "  DEBUG:                     Set to 'true' for verbose logging"
}

parse_args() {
    local command="$1"
    shift

    case "$command" in
        "apply")
            cmd_apply "$@"
            ;;
        "remove")
            cmd_remove "$@"
            ;;
        "info")
            cmd_info "$@"
            ;;
        "list")
            cmd_list "$@"
            ;;
        "help"|"-h"|"--help")
            usage
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            usage
            exit 1
            ;;
    esac
}

cmd_apply() {
    local mode="--all"
    local filter=""
    local journal="default"
    local local_path=""

    if [[ $# -eq 0 || "$1" == -* ]]; then
        log "ERROR" "A script ID is required for the apply command."
        usage
        exit 1
    fi

    local script_id="$1"
    shift

    while [[ $# -gt 0 ]]; do
        case $1 in
            --tag)
                mode="--tag"
                filter=$(get_required_arg "$1" "$2")
                shift 2
                ;;
            --all)
                mode="--all"
                shift
                ;;
            --local)
                local_path=$(get_required_arg "$1" "$2")
                shift 2
                ;;
            -*)
                log "ERROR" "Unknown option for apply: $1"
                usage
                exit 1
                ;;
            *)
                log "ERROR" "Unexpected argument: $1"
                usage
                exit 1
                ;;
        esac
    done

    command_apply "$script_id" "$mode" "$filter" "$local_path" 
}

cmd_remove() {
    local mode="--all"
    local filter=""
    local journal="default"

    if [[ $# -eq 0 || "$1" == -* ]]; then
        log "ERROR" "A script ID is required for the apply command."
        usage
        exit 1
    fi

    local script_id="$1"
    shift
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --id)
                mode="--id"
                filter="$2"
                shift 2
                ;;
            --tag)
                mode="--tag"
                filter="$2"
                shift 2
                ;;
            --all)
                mode="--all"
                shift
                ;;
            -*)
                log "ERROR" "Unknown option for remove: $1"
                usage
                exit 1
                ;;
            *)
                log "ERROR" "Unexpected argument: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    command_remove "$script_id" "$mode" "$filter"
}

cmd_info() {
    if [ $# -eq 0 ]; then
        log "ERROR" "Script ID is required for info command"
        usage
        exit 1
    fi
    
    show_script_info "$1"
}

cmd_list() {
    list_journals
}

main() {
    init_migrator
    
    if [ $# -eq 0 ]; then
        usage
        exit 0
    fi
    
    parse_args "$@"
}

main "$@"