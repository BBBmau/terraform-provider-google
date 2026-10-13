/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package builds

// Custom runners must consume TeamCity's exclusions themselves.
fun parallelTestsBatchScript(): String = """
    EXCLUDES_FILE="%system.teamcity.build.parallelTests.excludesFile%"
    if test -n "${'$'}EXCLUDES_FILE"; then
        if ! test -r "${'$'}EXCLUDES_FILE"; then
            echo "Cannot read Parallel Tests exclusions: ${'$'}EXCLUDES_FILE" >&2
            exit 1
        fi
        awk '!/^#/ && NF { sub(/\r${'$'}/, ""); print }' "${'$'}EXCLUDES_FILE" > batch-excludes.txt || exit 1
    else
        : > batch-excludes.txt
    fi
    BATCH_SKIP=${'$'}(awk 'BEGIN { printf "^(" } { if (NR > 1) printf "|"; gsub(/[][\\.^${'$'}*+?(){}|]/, "\\\\&"); printf "%s", ${'$'}0 } END { print ")${'$'}" }' batch-excludes.txt) || exit 1

    runBatchGoTests() {
        local user_skip="" skip="${'$'}BATCH_SKIP"
        local args=()
        while test "${'$'}#" -gt 0; do
            case "${'$'}1" in
                -skip=*|-test.skip=*) user_skip="${'$'}{1#*=}" ;;
                -skip|-test.skip)
                    if test "${'$'}#" -lt 2; then
                        echo "Missing value for ${'$'}1" >&2
                        return 1
                    fi
                    shift
                    user_skip="${'$'}1"
                    ;;
                *) args+=("${'$'}1") ;;
            esac
            shift
        done
        if test -n "${'$'}user_skip"; then
            if ! test -s batch-excludes.txt; then
                skip="${'$'}user_skip"
            elif [[ "${'$'}user_skip" == */* ]]; then
                echo "Cannot combine a subtest -skip filter with Parallel Tests exclusions" >&2
                return 1
            else
                skip="${'$'}skip|(${'$'}user_skip)"
            fi
        fi
        go test "${'$'}{args[@]}" -skip="${'$'}skip"
    }
""".trimIndent()

fun parallelTestsListScript(): String = """
    ./test-binary -test.list="%TEST_PREFIX%" > all-tests.txt || exit 1
    awk 'FILENAME == ARGV[1] { excluded[${'$'}0] = 1; next } !(${'$'}0 in excluded)' batch-excludes.txt all-tests.txt > batch-tests.txt || exit 1
""".trimIndent()

fun parallelTestsReportingScript(): String = """
    set -o pipefail
    go build -o ./teamcity-report-tests "%teamcity.build.checkoutDir%/scripts/teamcityreporter/main.go" || exit 1
    reportParallelTests() {
        ./teamcity-report-tests -format="${'$'}1"
    }
""".trimIndent()
