#!/bin/bash

# ============================================================
#  🚀 Database Migration Tool (golang-migrate wrapper)
# ============================================================
#  Usage:
#    ./scripts/migrate.sh up              - Apply all pending migrations
#    ./scripts/migrate.sh up <N>          - Apply N migrations
#    ./scripts/migrate.sh down            - Rollback last migration (with confirmation)
#    ./scripts/migrate.sh down <N>        - Rollback N migrations (with confirmation)
#    ./scripts/migrate.sh down all        - Rollback ALL migrations (with confirmation)
#    ./scripts/migrate.sh create <name>   - Create new migration files (up & down)
#    ./scripts/migrate.sh status          - Show current migration version
#    ./scripts/migrate.sh version         - Show current migration version
#    ./scripts/migrate.sh force <V>       - Force set version (fix dirty state)
#    ./scripts/migrate.sh goto <V>        - Migrate to a specific version
#    ./scripts/migrate.sh redo            - Rollback last then re-apply (with confirmation)
#    ./scripts/migrate.sh fresh           - Drop ALL & re-apply all (⚠️ DESTRUCTIVE)
#    ./scripts/migrate.sh list            - List all migration files
#    ./scripts/migrate.sh doctor          - Check system requirements
# ============================================================

set -euo pipefail

# ── Colors ───────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

# ── Paths ────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MIGRATIONS_DIR="$PROJECT_ROOT/migrations"
LOG_FILE="$PROJECT_ROOT/migrations/.migration_history.log"

# ── Load .env ────────────────────────────────────────────────
if [ -f "$PROJECT_ROOT/.env" ]; then
    set -a
    source "$PROJECT_ROOT/.env"
    set +a
fi

# ── Validate ─────────────────────────────────────────────────
if [ -z "${DATABASE_URL:-}" ]; then
    echo -e "${RED}✗ DATABASE_URL is not set!${NC}"
    echo -e "  Set it in ${CYAN}.env${NC} or export it."
    exit 1
fi

mkdir -p "$MIGRATIONS_DIR"

# ── Helpers ──────────────────────────────────────────────────

print_banner() {
    echo ""
    echo -e "${CYAN}${BOLD}  ╔══════════════════════════════════════╗${NC}"
    echo -e "${CYAN}${BOLD}  ║       Database Migration Tool        ║${NC}"
    echo -e "${CYAN}${BOLD}  ╚══════════════════════════════════════╝${NC}"
    echo ""
}

print_usage() {
    echo -e "${BOLD}Commands:${NC}"
    echo ""
    echo -e "  ${GREEN}up${NC}              Apply all pending migrations"
    echo -e "  ${GREEN}up ${DIM}<N>${NC}          Apply next N migrations"
    echo -e "  ${RED}down${NC}            Rollback last migration"
    echo -e "  ${RED}down ${DIM}<N>${NC}        Rollback N migrations"
    echo -e "  ${RED}down all${NC}        Rollback ALL migrations"
    echo -e "  ${BLUE}create ${DIM}<name>${NC}  Create new migration files"
    echo -e "  ${MAGENTA}redo${NC}            Rollback last & re-apply"
    echo -e "  ${RED}${BOLD}fresh${NC}           Drop ALL & re-apply (⚠️  destructive)"
    echo -e "  ${CYAN}status${NC}          Show current version & state"
    echo -e "  ${CYAN}version${NC}         Show current version number"
    echo -e "  ${YELLOW}force ${DIM}<V>${NC}      Force set version (fix dirty state)"
    echo -e "  ${YELLOW}goto ${DIM}<V>${NC}       Migrate to specific version"
    echo -e "  ${CYAN}list${NC}            List all migration files"
    echo -e "  ${CYAN}doctor${NC}          Check system requirements"
    echo ""
    echo -e "${DIM}Examples:${NC}"
    echo "  ./scripts/migrate.sh create create_users_table"
    echo "  ./scripts/migrate.sh up"
    echo "  ./scripts/migrate.sh down 2"
    echo ""
}

log_action() {
    local action="$1"
    local detail="${2:-}"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $action $detail" >> "$LOG_FILE"
}

confirm_action() {
    local message="$1"
    local require_yes="${2:-false}"

    echo ""
    echo -e "${RED}┌─────────────────────────────────────────────┐${NC}"
    echo -e "${RED}│  ⚠️   WARNING: DESTRUCTIVE ACTION            │${NC}"
    echo -e "${RED}└─────────────────────────────────────────────┘${NC}"
    echo ""
    echo -e "${YELLOW}$message${NC}"
    echo -e "${RED}This action may result in ${BOLD}DATA LOSS${NC}${RED}!${NC}"
    echo ""

    if [ "$require_yes" = "true" ]; then
        read -p "$(echo -e "${YELLOW}Type '${BOLD}yes${NC}${YELLOW}' to confirm: ${NC}")" confirm
        if [ "$confirm" != "yes" ]; then
            echo -e "${BLUE}✗ Action cancelled.${NC}"
            exit 0
        fi
    else
        read -p "$(echo -e "${YELLOW}Are you sure? (y/N): ${NC}")" confirm
        if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
            echo -e "${BLUE}✗ Action cancelled.${NC}"
            exit 0
        fi
    fi
    echo ""
}

get_current_version() {
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version 2>&1 || echo "no migration"
}

count_migration_files() {
    find "$MIGRATIONS_DIR" -maxdepth 1 -name "*.up.sql" 2>/dev/null | wc -l
}

check_migrate_installed() {
    if ! command -v migrate &> /dev/null; then
        echo -e "${RED}✗ 'migrate' CLI is not installed!${NC}"
        echo ""
        echo -e "${YELLOW}Install it with:${NC}"
        echo ""
        echo -e "  ${DIM}# Using Go${NC}"
        echo -e "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        echo ""
        echo -e "  ${DIM}# Using curl (Linux)${NC}"
        echo -e "  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz"
        echo -e "  sudo mv migrate /usr/local/bin/"
        echo ""
        exit 1
    fi
}

show_elapsed() {
    local start=$1
    local end
    end=$(date +%s%3N)
    local elapsed=$(( end - start ))

    if [ "$elapsed" -lt 1000 ]; then
        echo -e "${DIM}Done in ${elapsed}ms${NC}"
    else
        local seconds=$(( elapsed / 1000 ))
        echo -e "${DIM}Done in ${seconds}s${NC}"
    fi
}

# ── Commands ─────────────────────────────────────────────────

cmd_up() {
    local count="${1:-}"
    local start
    start=$(date +%s%3N)

    echo -e "${GREEN}▸ Running migrations UP...${NC}"
    echo ""

    if [ -n "$count" ]; then
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up "$count"
        log_action "UP" "applied $count migration(s)"
    else
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
        log_action "UP" "applied all pending migrations"
    fi

    echo ""
    echo -e "${GREEN}✓ Migrations applied successfully!${NC}"
    show_elapsed "$start"
}

cmd_down() {
    local count="${1:-1}"

    if [ "$count" = "all" ]; then
        local total
        total=$(count_migration_files)
        confirm_action "Rolling back ALL ${total} migration(s)." "true"

        local start
        start=$(date +%s%3N)

        echo -e "${RED}▸ Rolling back ALL migrations...${NC}"
        echo ""
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down -all
        log_action "DOWN" "rolled back ALL migrations"

        echo ""
        echo -e "${GREEN}✓ All migrations rolled back.${NC}"
        show_elapsed "$start"
    else
        confirm_action "Rolling back ${count} migration(s)."

        local start
        start=$(date +%s%3N)

        echo -e "${RED}▸ Rolling back ${count} migration(s)...${NC}"
        echo ""
        echo "y" | migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down "$count"
        log_action "DOWN" "rolled back $count migration(s)"

        echo ""
        echo -e "${GREEN}✓ Rolled back ${count} migration(s).${NC}"
        show_elapsed "$start"
    fi
}

cmd_create() {
    local name="$1"

    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$name"

    echo ""
    echo -e "${GREEN}✓ Created migration: ${BOLD}${name}${NC}"
    echo ""

    # Show created files
    local latest_files
    latest_files=$(ls -t "$MIGRATIONS_DIR"/*.sql 2>/dev/null | head -2)
    if [ -n "$latest_files" ]; then
        echo -e "${BLUE}Files created:${NC}"
        echo "$latest_files" | while IFS= read -r f; do
            echo -e "  ${CYAN}→${NC} $(basename "$f")"
        done
    fi

    log_action "CREATE" "$name"
    echo ""
    echo -e "${DIM}Edit the files in ${CYAN}migrations/${NC}${DIM} to define your schema changes.${NC}"
    echo ""
}

cmd_status() {
    echo -e "${BOLD}Migration Status${NC}"
    echo -e "${CYAN}──────────────────────────────────────────${NC}"

    local total
    total=$(count_migration_files)
    echo -e "  ${BLUE}Total migrations :${NC} $total"

    local version_output
    version_output=$(get_current_version)

    if echo "$version_output" | grep -q "no migration"; then
        echo -e "  ${BLUE}Current version  :${NC} ${YELLOW}none (no migrations applied)${NC}"
    elif echo "$version_output" | grep -q "dirty"; then
        echo -e "  ${BLUE}Current version  :${NC} ${RED}${version_output} (⚠️  DIRTY STATE)${NC}"
        echo ""
        echo -e "  ${YELLOW}Fix with:${NC} ./scripts/migrate.sh force <version>"
    else
        echo -e "  ${BLUE}Current version  :${NC} ${GREEN}${version_output}${NC}"
    fi

    echo -e "  ${BLUE}Database         :${NC} ${DIM}$(echo "$DATABASE_URL" | sed 's|://[^:]*:[^@]*@|://***:***@|')${NC}"
    echo -e "${CYAN}──────────────────────────────────────────${NC}"

    # Show recent history
    if [ -f "$LOG_FILE" ]; then
        echo ""
        echo -e "${BOLD}Recent History${NC} ${DIM}(last 5)${NC}"
        echo -e "${CYAN}──────────────────────────────────────────${NC}"
        tail -5 "$LOG_FILE" | while IFS= read -r line; do
            if echo "$line" | grep -q "UP"; then
                echo -e "  ${GREEN}↑${NC} $line"
            elif echo "$line" | grep -q "DOWN\|FRESH"; then
                echo -e "  ${RED}↓${NC} $line"
            else
                echo -e "  ${BLUE}•${NC} $line"
            fi
        done
        echo -e "${CYAN}──────────────────────────────────────────${NC}"
    fi
    echo ""
}

cmd_version() {
    local version_output
    version_output=$(get_current_version)
    echo -e "${BLUE}Current version:${NC} ${BOLD}${version_output}${NC}"
}

cmd_force() {
    local version="$1"
    echo -e "${YELLOW}▸ Forcing version to: ${version}${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$version"
    echo -e "${GREEN}✓ Version forced to ${version}${NC}"
    log_action "FORCE" "version set to $version"
}

cmd_goto() {
    local version="$1"
    confirm_action "Migrating to version ${version}. This may rollback migrations."

    local start
    start=$(date +%s%3N)

    echo -e "${BLUE}▸ Migrating to version ${version}...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" goto "$version"

    echo -e "${GREEN}✓ Migrated to version ${version}${NC}"
    log_action "GOTO" "migrated to version $version"
    show_elapsed "$start"
}

cmd_redo() {
    confirm_action "This will rollback the last migration and re-apply it."

    local start
    start=$(date +%s%3N)

    echo -e "${MAGENTA}▸ Redo: rolling back last migration...${NC}"
    echo "y" | migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down 1

    echo -e "${MAGENTA}▸ Redo: re-applying migration...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up 1

    echo ""
    echo -e "${GREEN}✓ Redo completed!${NC}"
    log_action "REDO" "rolled back and re-applied last migration"
    show_elapsed "$start"
}

cmd_fresh() {
    local total
    total=$(count_migration_files)

    echo -e "${RED}${BOLD}"
    echo "  ██████╗  █████╗ ███╗   ██╗ ██████╗ ███████╗██████╗ "
    echo "  ██╔══██╗██╔══██╗████╗  ██║██╔════╝ ██╔════╝██╔══██╗"
    echo "  ██║  ██║███████║██╔██╗ ██║██║  ███╗█████╗  ██████╔╝"
    echo "  ██║  ██║██╔══██║██║╚██╗██║██║   ██║██╔══╝  ██╔══██╗"
    echo "  ██████╔╝██║  ██║██║ ╚████║╚██████╔╝███████╗██║  ██║"
    echo "  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝"
    echo -e "${NC}"
    echo -e "${RED}${BOLD}  This will DROP ALL tables and re-run all ${total} migration(s)!${NC}"
    echo ""
    echo -e "${YELLOW}Type the database name to confirm:${NC}"

    local db_name
    db_name=$(echo "$DATABASE_URL" | sed -n 's|.*/\([^?]*\).*|\1|p')
    read -p "$(echo -e "${YELLOW}Database name: ${NC}")" input_name

    if [ "$input_name" != "$db_name" ]; then
        echo -e "${BLUE}✗ Database name doesn't match. Action cancelled.${NC}"
        exit 0
    fi

    local start
    start=$(date +%s%3N)

    echo ""
    echo -e "${RED}▸ Dropping all migrations...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" drop -f

    echo -e "${GREEN}▸ Re-applying all migrations...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up

    echo ""
    echo -e "${GREEN}✓ Fresh migration completed!${NC}"
    log_action "FRESH" "dropped all and re-applied $total migration(s)"
    show_elapsed "$start"
}

cmd_list() {
    echo -e "${BOLD}Migration Files${NC}"
    echo -e "${CYAN}──────────────────────────────────────────${NC}"

    local files
    files=$(find "$MIGRATIONS_DIR" -maxdepth 1 -name "*.sql" -printf "%f\n" 2>/dev/null | sort)

    if [ -z "$files" ]; then
        echo -e "  ${YELLOW}No migration files found.${NC}"
        echo ""
        echo -e "  Create one with: ${GREEN}./scripts/migrate.sh create <name>${NC}"
    else
        echo "$files" | while IFS= read -r f; do
            if echo "$f" | grep -q "\.up\.sql$"; then
                echo -e "  ${GREEN}↑${NC} $f"
            elif echo "$f" | grep -q "\.down\.sql$"; then
                echo -e "  ${RED}↓${NC} $f"
            else
                echo -e "  ${BLUE}•${NC} $f"
            fi
        done
    fi

    echo -e "${CYAN}──────────────────────────────────────────${NC}"
    echo ""
}

cmd_doctor() {
    echo -e "${BOLD}System Check${NC}"
    echo -e "${CYAN}──────────────────────────────────────────${NC}"

    # Check migrate CLI
    if command -v migrate &> /dev/null; then
        local ver
        ver=$(migrate -version 2>&1 || echo "unknown")
        echo -e "  ${GREEN}✓${NC} migrate CLI       ${DIM}($ver)${NC}"
    else
        echo -e "  ${RED}✗${NC} migrate CLI       ${RED}not installed${NC}"
    fi

    # Check psql
    if command -v psql &> /dev/null; then
        local psql_ver
        psql_ver=$(psql --version 2>&1 | head -1)
        echo -e "  ${GREEN}✓${NC} psql              ${DIM}($psql_ver)${NC}"
    else
        echo -e "  ${YELLOW}○${NC} psql              ${YELLOW}not installed (optional)${NC}"
    fi

    # Check DATABASE_URL
    if [ -n "${DATABASE_URL:-}" ]; then
        echo -e "  ${GREEN}✓${NC} DATABASE_URL      ${DIM}set${NC}"
    else
        echo -e "  ${RED}✗${NC} DATABASE_URL      ${RED}not set${NC}"
    fi

    # Check .env
    if [ -f "$PROJECT_ROOT/.env" ]; then
        echo -e "  ${GREEN}✓${NC} .env file         ${DIM}found${NC}"
    else
        echo -e "  ${YELLOW}○${NC} .env file         ${YELLOW}not found${NC}"
    fi

    # Check migrations dir
    local total
    total=$(count_migration_files)
    echo -e "  ${GREEN}✓${NC} migrations/       ${DIM}${total} migration(s)${NC}"

    # Check DB connectivity
    echo -e ""
    echo -ne "  ${BLUE}◌${NC} Database connection... "
    if migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version &> /dev/null; then
        echo -e "\r  ${GREEN}✓${NC} Database          ${GREEN}connected${NC}      "
    else
        # Could be "no migration" which is still a successful connection
        local result
        result=$(migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version 2>&1 || true)
        if echo "$result" | grep -q "no migration"; then
            echo -e "\r  ${GREEN}✓${NC} Database          ${GREEN}connected (no migrations)${NC}      "
        else
            echo -e "\r  ${RED}✗${NC} Database          ${RED}connection failed${NC}      "
        fi
    fi

    echo -e "${CYAN}──────────────────────────────────────────${NC}"
    echo ""
}

# ── Main ─────────────────────────────────────────────────────

print_banner

# Doctor doesn't need migrate CLI check
if [ "${1:-}" = "doctor" ]; then
    cmd_doctor
    exit 0
fi

check_migrate_installed

case "${1:-}" in
    up)
        cmd_up "${2:-}"
        ;;
    down)
        cmd_down "${2:-1}"
        ;;
    create)
        if [ -z "${2:-}" ]; then
            echo -e "${RED}✗ Please provide a migration name.${NC}"
            echo -e "  Example: ${GREEN}./scripts/migrate.sh create create_users_table${NC}"
            exit 1
        fi
        cmd_create "$2"
        ;;
    status)
        cmd_status
        ;;
    version)
        cmd_version
        ;;
    force)
        if [ -z "${2:-}" ]; then
            echo -e "${RED}✗ Please provide a version number.${NC}"
            exit 1
        fi
        cmd_force "$2"
        ;;
    goto)
        if [ -z "${2:-}" ]; then
            echo -e "${RED}✗ Please provide a version number.${NC}"
            exit 1
        fi
        cmd_goto "$2"
        ;;
    redo)
        cmd_redo
        ;;
    fresh)
        cmd_fresh
        ;;
    list)
        cmd_list
        ;;
    *)
        print_usage
        exit 1
        ;;
esac
