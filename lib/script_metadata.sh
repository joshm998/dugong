#!/usr/bin/env bash

get_metadata_value() {
    local script_file=$1
    local key=$2
    
    if [ ! -f "$script_file" ]; then
        return 1
    fi
    
    grep "^# $key:" "$script_file" 2>/dev/null | head -n 1 | cut -d':' -f2- | sed 's/^[[:space:]]*//;s/[[:space:]]*$//'
}

display_script_details() {
    local script_file=$1
    local is_gallery=${2:-false}
    
    local id desc version author tags optional
    id=$(get_metadata_value "$script_file" "id")
    desc=$(get_metadata_value "$script_file" "description")
    version=$(get_metadata_value "$script_file" "version")
    author=$(get_metadata_value "$script_file" "author")
    tags=$(get_metadata_value "$script_file" "tags")
    optional=$(get_metadata_value "$script_file" "optional")
    
    echo ""
    echo -e "${COLOR_BLUE}━━━ Script Details ━━━${COLOR_NC}"
    echo -e "${COLOR_CYAN}ID:${COLOR_NC} $id"
    echo -e "${COLOR_CYAN}Description:${COLOR_NC} $desc"
    [ -n "$version" ] && echo -e "${COLOR_CYAN}Version:${COLOR_NC} $version"
    [ -n "$author" ] && echo -e "${COLOR_CYAN}Author:${COLOR_NC} $author"
    [ -n "$tags" ] && echo -e "${COLOR_CYAN}Tags:${COLOR_NC} $tags"
    [ "$optional" == "true" ] && echo -e "${COLOR_CYAN}Optional:${COLOR_NC} Yes"
    
    if [ "$is_gallery" == "true" ]; then
        echo -e "${COLOR_CYAN}Gallery URL:${COLOR_NC} https://a.com/gallery/$id"
    fi
    echo -e "${COLOR_BLUE}━━━━━━━━━━━━━━━━━━━━━━━${COLOR_NC}"
    echo ""
}

validate_script_file() {
    local script_file=$1
    
    if [ ! -f "$script_file" ]; then
        log "ERROR" "Script file not found: $script_file"
        return 1
    fi
    
    local id desc
    id=$(get_metadata_value "$script_file" "id")
    desc=$(get_metadata_value "$script_file" "description")
    
    if [ -z "$id" ]; then
        log "ERROR" "Script missing required metadata: id"
        return 1
    fi
    
    if [ -z "$desc" ]; then
        log "ERROR" "Script missing required metadata: description"
        return 1
    fi
    
    if ! grep -q "^apply_changes()" "$script_file"; then
        log "ERROR" "Script missing required function: apply_changes()"
        return 1
    fi
    
    if ! grep -q "^remove_changes()" "$script_file"; then
        log "ERROR" "Script missing required function: remove_changes()"
        return 1
    fi
    
    return 0
}

script_has_tag() {
    local script_file=$1
    local tag_to_find=$2
    
    local tags
    tags=$(get_metadata_value "$script_file" "tags")
    
    if [ -z "$tags" ]; then
        return 1
    fi
    
    # Convert comma-separated tags to space-separated and check
    echo "$tags" | tr ',' ' ' | grep -q -w "$tag_to_find"
}

find_script_by_id() {
    local id_to_find=$1
    
    for script in "$SCRIPTS_DIR"/*.sh; do
        if [ -f "$script" ]; then
            local script_id
            script_id=$(get_metadata_value "$script" "id")
            if [ "$script_id" == "$id_to_find" ]; then
                echo "$script"
                return 0
            fi
        fi
    done
    return 1
}