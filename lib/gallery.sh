#!/usr/bin/env bash

download_scripts_from_gallery() {
  local GITHUB_USER="joshm998"
  local GITHUB_REPO="dugong-gallery"
  local GITHUB_BRANCH="main"
  local script_id="$1"

  for cmd in curl jq; do
    if ! command_exists "$cmd"; then
      echo "Error: Required command '$cmd' is not installed."
      return 1
    fi
  done

  local CACHE_BASE_DIR="$MIGRATOR_HOME/cache"
  local SCRIPT_CACHE_DIR="$CACHE_BASE_DIR/$script_id"
  local LOCAL_COMMIT_FILE="$SCRIPT_CACHE_DIR/.last_commit"
  mkdir -p "$SCRIPT_CACHE_DIR"

  echo "--- Fetching Script from Gallery: $script_id ---"

  local commit_api_url="https://api.github.com/repos/$GITHUB_USER/$GITHUB_REPO/commits?path=$script_id&page=1&per_page=1"
  local remote_commit
  remote_commit=$(curl -sL "$commit_api_url" | jq -r '.[0].sha // empty')

  if [[ -z "$remote_commit" ]]; then
      echo "Error: Could not find script_id '$script_id' or repo is empty."
      return 1
  fi

  local local_commit=""
  if [[ -f "$LOCAL_COMMIT_FILE" ]]; then
    local_commit=$(cat "$LOCAL_COMMIT_FILE")
  fi

  echo "Remote commit: $remote_commit"
  echo "Local commit:  $local_commit"

  if [[ "$remote_commit" != "$local_commit" ]]; then
    echo "Remote has been updated. Refreshing cache..."

    # TODO Delete any existing SH files

    local contents_api_url="https://api.github.com/repos/$GITHUB_USER/$GITHUB_REPO/contents/$script_id?ref=$GITHUB_BRANCH"
    
    local file_urls
    file_urls=$(curl -sL "$contents_api_url" | jq -r '.[] | select(.type == "file") | .download_url')
    
    if [[ -z "$file_urls" ]]; then
      echo "Warning: No files found in the repo directory."
    else
      for url in $file_urls; do
        local filename
        filename=$(basename "$url")
        echo "Downloading $filename..."
        curl -sL "$url" -o "$SCRIPT_CACHE_DIR/$filename"
      done
    fi

    echo "$remote_commit" > "$LOCAL_COMMIT_FILE"
    echo "Cache refresh complete."
  else
    echo "Cache is up to date."
  fi

  echo ""

  local details_file="$SCRIPT_CACHE_DIR/details.json"
  if [[ -f "$details_file" ]]; then
    echo "--- details.json (from cache) ---"
    cat "$details_file"
    echo
  fi

  shopt -s nullglob # Ensure loop doesn't run if no files match
  for script_file in "$SCRIPT_CACHE_DIR"/[0-9][0-9][0-9]_*.sh; do
    echo "--- Script: $(basename "$script_file") (from cache) ---"
    cat "$script_file"
    echo
  done
}


check_gallery_connectivity() {    
    local test_url="$MIGRATOR_GALLERY_URL"
    if curl -sSL --connect-timeout 5 "$test_url" >/dev/null 2>&1; then
        log "DEBUG" "Gallery connectivity: OK"
        return 0
    else
        log "WARN" "Gallery not accessible: $test_url"
        return 1
    fi
}


show_script_info() {
    local script_id=$1
    
    echo "-----------------------------"
    echo -e "${COLOR_BLUE}Application Name:${COLOR_NC}"
    echo -e "  ${COLOR_YELLOW}○${COLOR_NC} ${script_id}"

    echo -e "${COLOR_BLUE}Gallery URL:${COLOR_NC}"
    echo -e "  ${COLOR_YELLOW}○${COLOR_NC} https://dugong.dev/gallery/${script_id}"
    
    echo -e "${COLOR_BLUE}Application Status:${COLOR_NC}"
    local found_in_journal=false
    
    for journal_file in "$JOURNALS_DIR"/*.journal; do
        if [ -f "$journal_file" ]; then
            local journal_name
            journal_name=$(basename "$journal_file" .journal)
            
            if is_script_applied "$script_id" "$journal_name"; then
                echo -e "  ${COLOR_GREEN}✓${COLOR_NC} Applied in journal: $journal_name"
                found_in_journal=true
            fi
        fi
    done
    
    if [ "$found_in_journal" == "false" ]; then
        echo -e "  ${COLOR_YELLOW}○${COLOR_NC} Not applied in any journal"
    fi
    
    echo "-----------------------------"
    echo ""
}