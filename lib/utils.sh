#!/usr/bin/env bash

init_migrator() {
    log "INFO" "Initializing dugong environment"
    
    mkdir -p "$MIGRATOR_HOME"
    mkdir -p "$SCRIPTS_DIR"
    mkdir -p "$JOURNALS_DIR"
    mkdir -p "$CACHE_DIR"\
    
    log "DEBUG" "Dugong home: $MIGRATOR_HOME"
    log "DEBUG" "Scripts directory: $SCRIPTS_DIR"
    log "DEBUG" "Journals directory: $JOURNALS_DIR"
}

install_dependencies() {
    local dependencies=("curl" "git" "jq")
    local missing_deps=()
    
    for dep in "${dependencies[@]}"; do
        if ! command_exists "$dep"; then
            missing_deps+=("$dep")
        fi
    done
    
    if [ ${#missing_deps[@]} -eq 0 ]; then
        log "DEBUG" "All dependencies are satisfied"
        return 0
    fi
    
    log "WARN" "Missing dependencies: ${missing_deps[*]}"
    
    if command_exists "apt-get"; then
        log "INFO" "Installing dependencies with apt-get"
        for dep in "${missing_deps[@]}"; do
            sudo apt-get install -y "$dep"
        done
    elif command_exists "yum"; then
        log "INFO" "Installing dependencies with yum"
        for dep in "${missing_deps[@]}"; do
            sudo yum install -y "$dep"
        done
    elif command_exists "dnf"; then
        log "INFO" "Installing dependencies with dnf"
        for dep in "${missing_deps[@]}"; do
            sudo dnf install -y "$dep"
        done
    else
        log "ERROR" "No supported package manager found. Please install: ${missing_deps[*]}"
        return 1
    fi
}

confirm_action() {
    local prompt="$1"
    local default=${2:-"n"}
    
    echo -n "$prompt"
    if [ "$default" == "y" ]; then
        echo -n " [Y/n]: "
    else
        echo -n " [y/N]: "
    fi
    
    read -r response
    response=${response:-$default}
    
    case "$response" in
        [yY]|[yY][eE][sS])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

get_required_arg() {
    local option_name="$1"
    local arg_value="$2"

    if [[ -n "$arg_value" && "$arg_value" != -* ]]; then
        echo "$arg_value"
    else
        log "ERROR" "Option '$option_name' requires an argument."
        usage
        exit 1
    fi
}