#!/usr/bin/env bash

is_script_applied() {
    local script_id=$1
    local journal=$2
    local journal_file="$JOURNALS_DIR/${journal}.journal"
    
    [ -f "$journal_file" ] && grep -qx "$script_id" "$journal_file"
}

record_script_in_journal() {
    local script_id=$1
    local journal=$2
    local journal_file="$JOURNALS_DIR/${journal}.journal"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    echo "$script_id" >> "$journal_file"
    echo "[$timestamp] Applied: $script_id" >> "$JOURNALS_DIR/${journal}.log"
    
    log "DEBUG" "Recorded $script_id in journal $journal"
}

remove_script_from_journal() {
    local script_id=$1
    local journal=$2
    local journal_file="$JOURNALS_DIR/${journal}.journal"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    if [ -f "$journal_file" ]; then
    echo $journal_file

        while IFS= read -r line || [[ -n "$line" ]]; do
            [[ "$line" != "$script_id" ]] && echo "$line"
        done < "$journal_file" > "${journal_file}.tmp"

        mv "${journal_file}.tmp" "${journal_file}"

        echo "[$timestamp] Removed: $script_id" >> "$JOURNALS_DIR/${journal}.log"
        log "DEBUG" "Removed $script_id from journal $journal"
    fi
}

get_journal_file() {
    local journal=$1
    echo "$JOURNALS_DIR/${journal}.journal"
}

list_journals() {
    echo ""
    echo -e "${COLOR_BLUE}Available Journals:${COLOR_NC}"
    echo ""
    
    local found_journals=false
    for journal_file in "$JOURNALS_DIR"/*.journal; do
        if [ -f "$journal_file" ]; then
            found_journals=true
            local journal_name
            journal_name=$(basename "$journal_file" .journal)
            local count
            count=$(wc -l < "$journal_file" 2>/dev/null || echo "0")
            
            echo -e "  ${COLOR_CYAN}${journal_name}${COLOR_NC} (${count} applied scripts)"
        fi
    done
    
    if [ "$found_journals" == "false" ]; then
        echo "  No journals found"
    fi
    echo ""
}
