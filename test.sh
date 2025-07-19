#!/bin/bash
# shellcheck disable=SC2086

declare -i EXECUTED=0
declare -i PASSED=0
declare -i FAILED=0

GOOS_ENV=${GOOS:-$(go env GOOS)}
GOARCH_ENV=${GOARCH:-$(go env GOARCH)}
BIN_PATH="./bin/goini-$GOOS_ENV-$GOARCH_ENV"

# Using simple indexed arrays instead of associative arrays to avoid subscript issues
declare -a test_names=(
  "TEST_01" "TEST_02" "TEST_03" "TEST_04" "TEST_05" "TEST_06" "TEST_07" "TEST_08"
  "TEST_09" "TEST_10" "TEST_11" "TEST_12" "TEST_13" "TEST_14" "TEST_15" "TEST_16"
  "TEST_17" "TEST_18" "TEST_19" "TEST_20"
)

declare -a test_commands=(
  "-section DEFAULT -has-section-key user"
  "-section DEFAULT -has-section-key non_existent_key"
  "-section DEFAULT -key user -has-section-key-value Yeshua"
  "-section DEFAULT -key user -has-section-key-value No"
  "-are-sections-present=DEFAULT,extra"
  "-are-sections-present=DEFAULT,non_existent_section"
  "-sections"
  "-sections -csv"
  "-sections -json"
  "-sections -yaml"
  "-section=DEFAULT -list-keys"
  "-section=DEFAULT -list-keys -json"
  "-section DEFAULT -list-key-values"
  "-section DEFAULT -list-key-values -json"
  "-add-section new_section"
  "-add-section new_section"
  "-section DEFAULT -key new_setting -value 123 -add-key"
  "-section DEFAULT -key user -value something -add-key"
  "-section DEFAULT -key user -value NewUserValue -modify-key"
  "-section DEFAULT -key non_existent_key -value some_value -modify-key"
)

declare -a test_descriptions=(
  "Relying on exit code, does 'section' have 'key'?"
  "Relying on exit code, does 'section' have 'key'?"
  "Relying on exit code, does 'section' 'key' have value 'value'?"
  "Relying on exit code, does 'section' 'key' have value 'value'?"
  "Relying on exit code, are 'section1' and 'section2' present?"
  "Relying on exit code, are 'section1' and 'section2' present?"
  "Using STDOUT, return list of 'sections' in ini file"
  "Using STDOUT, return list of 'sections' in ini file"
  "Using STDOUT, return list of 'sections' in ini file"
  "Using STDOUT, return list of 'sections' in ini file"
  "Using STDOUT, return a list of keys in 'section' (by name)"
  "Using STDOUT, return a list of keys in 'section' (by name)"
  "Using STDOUT, return a list of key/values in 'section' (by name)"
  "Using STDOUT, return a list of key/values in 'section' (by name)"
  "Using exit code for success status, add new section to ini file"
  "Try adding again (expecting exit 1)"
  "Using exit code for success status, in section 'section', add 'key' with value 'value'"
  "Try adding existing key (expecting exit 1)"
  "Using exit code for success status, in section 'section', modify 'key' with new value 'value'"
  "Try modifying non-existent key (expecting exit 1)"
)

declare -a test_wants=(
  0 1 0 1 0 1 0 0 0 0 0 0 0 0 0 1 0 1 0 1
# 1 2 3 4 5 6 7 8 9 A B C D E F G H I J K
#                  10  12  14  16  18  20
)
declare -a ini_reset=(
  1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 0 1 1 1 1
# 1 2 3 4 5 6 7 8 9 A B C D E F G H I J K
#                  10  12  14  16  18  20
)

function run_test(){
  local args="${1}"
  local -i want="${2:-0}"
  local -i reset_ini="${3:-0}"

  if [[ "${want}" == "3" ]]; then
    return 0
  fi

  if (( reset_ini == 1 )); then
    reset_sample_ini
  fi

  ((EXECUTED++))
  local tmp_file="/tmp/goini-test-${EXECUTED}"
  if $BIN_PATH -ini=sample.ini $args 1> /dev/null 2> $tmp_file ; then # test passed
    if [[ "${want}" == "0" ]]; then # wanted 0; got 0; we good
      ((PASSED++))
      echo "PASSED"
    else # wanted 1; got 0; el problemo
      echo "FAILED"
      ((FAILED++))
      cat $tmp_file
    fi
  else # test failed
    if [[ "${want}" == "1" ]]; then # wanted 1; got 1 ; we good
      ((PASSED++))
      echo "PASSED"
    else # wanted 0; got 1; el problemo
      echo "FAILED"
      ((FAILED++))
      stderr="$(cat $tmp_file)"
      if [[ -n "${stderr}" ]]; then
        printf "\tSTDERR = %s" "${stderr}"
        echo
      fi
    fi
  fi
  rm -rf $tmp_file
}

function reset_sample_ini(){
  {
    echo "[DEFAULT]"
    echo "user = Yeshua"
    echo "key = 369"
    echo "port = 1776"
    echo "country = ISRAEL"
    echo ""
    echo "[extra]"
    echo "ssh_key = ~/.ssh/id_rsa"
    echo "ssh_key_pub = ~/.ssh/id_rsa.pub"
  } | tee "sample.ini" > /dev/null
}

function build_goini_binary() {
  echo "Building goini binary..."

  [ -z "${BIN_PATH}" ] && { echo "invalid BIN_PATH in runtime"; exit 1; }

  [ -f "${BIN_PATH}" ] && rm -rf "${BIN_PATH}" && echo "Clean successful: ${BIN_PATH}"

  if ! go build -ldflags "-s -w" -o "$BIN_PATH" .; then
    echo "ERROR: Failed to build goini binary. Aborting tests."
    exit 1
  fi
  echo "Build successful: $BIN_PATH"
  echo
}

function main(){
  build_goini_binary

  # Iterate through tests using array indices
  for i in "${!test_names[@]}"; do
    echo -n "${test_names[i]}: ${test_descriptions[i]}..."
    printf "%s" ""
    echo
    run_test "${test_commands[i]}" "${test_wants[i]}" "${ini_reset[i]}"
  done

  echo && reset_sample_ini

  if (( FAILED == 0 )); then
    echo "All ${EXECUTED} tests passed. PASS"
    exit 0
  else
    echo "${FAILED} of ${EXECUTED} tests failed..."
    printf "\tTotal: %s\n\tPassed: %s\n\tFailed: %s" $EXECUTED $PASSED $FAILED
    echo
    exit 1
  fi
}

main "$@"
