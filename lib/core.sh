#!/usr/bin/env bash

CORE_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

source "$CORE_SCRIPT_DIR/logger.sh"
source "$CORE_SCRIPT_DIR/journal.sh"
source "$CORE_SCRIPT_DIR/script_metadata.sh"
source "$CORE_SCRIPT_DIR/utils.sh"
source "$CORE_SCRIPT_DIR/gallery.sh"

execute_script_function() {
    local script_file=$1
    local function_name=$2
    local script_id=$3

    log "DEBUG" "Executing function '$function_name' from script '$script_id'"

    (
        set -e
        # shellcheck source=/dev/null
        source "$script_file"
        
        if [ "$(type -t "$function_name")" = "function" ]; then
            "$function_name"
        else
            log "ERROR" "Function '$function_name' not found in script '$script_id'"
            exit 1
        fi
    )
    return $?
}

command_apply() {
    local id=$1
    local mode=$2
    local filter=$3
    local local_path=$4
    local journal=$id
    
    local scripts_dir=$SCRIPTS_DIR

    if [ "$local_path" != "" ]; then
        log "INFO" "Starting local 'apply' command (id: ${id}, mode: ${mode}, filter: ${filter:-none}, local_path: ${local_path})"

        rm -rf "${SCRIPTS_DIR:?}/$id"
        mkdir -p "$SCRIPTS_DIR/$id"
        cp -rf "$local_path/$id" "$SCRIPTS_DIR"

        scripts_dir="$SCRIPTS_DIR/$id"
    else
        log "INFO" "Starting 'apply' command (id: ${id}, mode: ${mode}, filter: ${filter:-none})"
        download_scripts_from_gallery "$id"
    fi
    
    local scripts_to_scan
    scripts_to_scan=$(find "$scripts_dir" -maxdepth 1 -type f -name "*.sh" | sort)

    if [ -z "$scripts_to_scan" ]; then
        log "WARN" "No scripts found in '$scripts_dir'"
        log "INFO" "Use 'migrator install <script_id>' to install scripts"
        return
    fi
    
    local scripts_to_apply=()
    local any_applied=0
    
    for script_file in $scripts_to_scan; do
        local migration_id
        migration_id=$(basename "$script_file" .sh)

        if is_script_applied "$migration_id" "$id"; then
            log "DEBUG" "Skipping script '${migration_id}': Already applied in journal '$id'"
            continue
        fi

        scripts_to_apply+=("$script_file")
    done
    
    if [ ${#scripts_to_apply[@]} -eq 0 ]; then
        log "SUCCESS" "System is up-to-date for the given filter in journal '$journal'"
        return
    fi
    
    echo ""
    echo -e "${COLOR_BLUE}━━━ Scripts to Apply ━━━${COLOR_NC}"
    for script_file in "${scripts_to_apply[@]}"; do
        local migration_id
        migration_id=$(basename "$script_file" .sh)
        # desc=$(get_metadata_value "$script_file" "description")
        echo -e "  ${COLOR_CYAN}${migration_id}${COLOR_NC}"
    done
    echo -e "${COLOR_BLUE}━━━━━━━━━━━━━━━━━━━━━━━━${COLOR_NC}"
    echo ""
    
    if ! confirm_action "Apply these ${#scripts_to_apply[@]} script(s) to journal '$journal'?"; then
        log "INFO" "Apply operation cancelled"
        return
    fi
    
    for script_file in "${scripts_to_apply[@]}"; do
        local migration_id
        migration_id=$(basename "$script_file" .sh)

        echo ""

        log "INFO" "Applying script '${migration_id}'"
        if execute_script_function "$script_file" "apply_changes" "$migration_id"; then
            log "SUCCESS" "Successfully applied '${migration_id}'"
            record_script_in_journal "$migration_id" "$journal"
            any_applied=1
        else
            log "ERROR" "Failed to apply script '${migration_id}'"
            if confirm_action "Continue with remaining scripts?"; then
                continue
            else
                log "ERROR" "Apply operation halted by user"
                exit 1
            fi
        fi
    done
    
    if [ $any_applied -eq 1 ]; then
        log "SUCCESS" "Apply command finished successfully"
    else
        log "INFO" "No scripts were applied"
    fi
}

command_remove() {
    local id=$1
    local mode=$2
    local filter=$3
    local journal=$id

    local script_dir=$SCRIPTS_DIR/$id
    
    local journal_contents=$(get_journal_file "$journal")

    scripts=""
    while read -r script || [[ -n "$script" ]]; do
        scripts="$script"$'\n'"$scripts"
    done < "$journal_contents"

    local any_removed=0
    local scripts_to_remove=()

    while read -r script; do
        if [ -z "$script" ]; then
            continue
        fi
        scripts_to_remove+=("$script_dir/$script.sh")
    done <<< "$scripts"

    if [ ${#scripts_to_remove[@]} -eq 0 ]; then
        log "SUCCESS" "System is up-to-date for the given filter in journal '$journal'"
        return
    fi
    
    echo ""
    echo -e "${COLOR_BLUE}━━━ Scripts to Remove ━━━${COLOR_NC}"
    for script_file in "${scripts_to_remove[@]}"; do
        local migration_id
        migration_id=$(basename "$script_file" .sh)
        # desc=$(get_metadata_value "$script_file" "description")
        echo -e "  ${COLOR_CYAN}${migration_id}${COLOR_NC}"
    done
    echo -e "${COLOR_BLUE}━━━━━━━━━━━━━━━━━━━━━━━━${COLOR_NC}"
    echo ""

    if ! confirm_action "Remove these ${#scripts_to_remove[@]} script(s) from journal '$journal'?"; then
        log "INFO" "Remove operation cancelled"
        return
    fi
    

    for script_file in "${scripts_to_remove[@]}"; do
        local migration_id
        migration_id=$(basename "$script_file" .sh)

        echo ""

        log "INFO" "Removing script '${migration_id}'"
        if execute_script_function "$script_file" "remove_changes" "$migration_id"; then
            log "SUCCESS" "Successfully removed '${migration_id}'"
            remove_script_from_journal "$migration_id" "$journal"
            any_removed=1
        else
            log "ERROR" "Failed to remove script '${migration_id}'"
            if confirm_action "Continue with remaining scripts?"; then
                continue
            else
                log "ERROR" "Remove operation halted by user"
                exit 1
            fi
        fi
    done
    
    if [ "$any_removed" -eq 1 ]; then
        log "SUCCESS" "Apply command finished successfully"
    else
        log "INFO" "No scripts were applied"
    fi
}