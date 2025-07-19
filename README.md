# GO INI

This utility is called `goini` and its designed to help you interact with `.ini` files. 

## Install 

You can use the `go install ...` method to install `goini` on your system:

```bash
go install github.com/andreimerlescu/goini@latest
```

Or you can use the binary release of the `goini` directly.

```bash
curl -sL https://github.com/andreimerlescu/goini/releases/download/v1.0.0/goini-linux-amd64 \
         --ouput /tmp/goini
chmod +x /tmp/goini
sudo mv /tmp/goini /usr/local/bin/goini
which goini
```

## Configuration

| Property Name                                  |   Type   | Default Value | Description                                                                                                                |
|------------------------------------------------|:--------:|:-------------:|:---------------------------------------------------------------------------------------------------------------------------|
| `-ini` <sup style="color: red;">REQUIRED</sup> | `String` |      ``       | Path to ini file to process                                                                                                |
| `-section`                                     | `String` |      ``       | Specify section name for operations like --has-section-key, --list-keys, --list-key-values, --add-key, --modify-key.       |
| `-sections`                                    |  `Bool`  |    `false`    | If true, only the sections from the ini file will be displayed                                                             |
| `-add-section`                                 | `String` |      ``       | If set, adds a new section to the ini file. Exit code 0 on success.                                                        |
| `-has-section`                                 | `String` |      ``       | If set, exit code will respond if the section exists or not in the --ini file                                              |
| `-has-section-key`                             | `String` |      ``       | If set, exit code will respond if the section has the key in the --ini file. Requires --section.                           |
| `-has-section-key-value`                       | `String` |      ``       | If set, exit code will respond if the section key has the specified value in the --ini file. Requires --section and --key. |
| `-are-sections-present`                        | `String` |      ``       | Comma-separated list of section names. Exit code will be 0 if all are present, 1 otherwise.                                |
| `-key`                                         | `String` |      ``       | Specify key name for operations like --has-section-key-value, --add-key, --modify-key.                                     |
| `-value`                                       | `String` |      ``       | Specify value for operations like --has-section-key-value, --add-key, --modify-key.                                        |
| `-add-key`                                     | `String` |      ``       | If true, adds a new key-value pair to the specified --section. Requires --key and --value. Exit code 0 on success.         |
| `-modify-key`                                  |  `Bool`  |    `false`    | If true, modifies an existing key's value in the specified --section. Requires --key and --value. Exit code 0 on success.  |
| `-list-keys`                                   |  `Bool`  |    `false`    | If true, returns a list of keys in the specified --section using STDOUT.                                                   |
| `-list-key-values`                             |  `Bool`  |    `false`    | If true, returns a list of key/values in the specified --section using STDOUT.                                             |
| `-csv`                                         |  `Bool`  |    `false`    | Output as CSV.                                                                                                             |
| `-json`                                        |  `Bool`  |    `false`    | Output as JSON.                                                                                                            |
| `-yaml`                                        |  `Bool`  |      ``       | Output as YAML.                                                                                                            |


## Usage Examples

This is a called `sample.ini` and it's going to be used for the **Usage Examples** of _goini_. 

```yaml
[default]
user = Yeshua
key = 369
port = 1776
country = ISRAEL

[extra]
ssh_key = ~/.ssh/id_rsa
ssh_key_pub = ~/.ssh/id_rsa.pub
```

### Relying on exit code, does "section" have "key"?

```bash
./goini --ini sample.ini --section default --has-section-key user
echo $? # Output: 0 (success, key 'user' exists in 'default')

./goini --ini sample.ini --section default --has-section-key non_existent_key
echo $? # Output: 1 (failure, key 'non_existent_key' does not exist)
```

### Relying on exit code, does "section" "key" have value "value"?

```bash
./goini --ini sample.ini --section default --key user --has-section-key-value Yeshua
echo $? # Output: 0

./goini --ini sample.ini --section default --key user --has-section-key-value No
echo $? # Output: 1
```

### Task 5 - Relying on exit code, are "section1" and "section2" present?

```bash
./goini --ini sample.ini --are-sections-present default,extra
echo $? # Output: 0

./goini --ini sample.ini --are-sections-present default,non_existent_section
echo $? # Output: 1
```

### Using STDOUT, return list of "sections" in ini file

```bash
./goini --ini sample.ini --sections
# Output:
# default
# extra

./goini --ini sample.ini --sections --csv
# Output: default,extra

./goini --ini sample.ini --sections --json
# Output:
# [
#   "default",
#   "extra"
# ]

./goini --ini sample.ini --sections --yaml
# Output:
# - default
# - extra
```

### Using STDOUT, return a list of keys in "section" (by name)

```bash
./goini --ini sample.ini --section default --list-keys
# Output:
# user
# key
# port
# country

./goini --ini sample.ini --section default --list-keys --json
# Output:
# [
#   "user",
#   "key",
#   "port",
#   "country"
# ]
```

### Using STDOUT, return a list of key/values in "section" (by name)

```bash
./goini --ini sample.ini --section default --list-key-values
# Output:
# user = Yeshua
# key = 369
# port = 1776
# country = ISRAEL

./goini --ini sample.ini --section default --list-key-values --json
# Output:
# {
#   "country": "ISRAEL",
#   "key": "369",
#   "port": "1776",
#   "user": "Yeshua"
# }
```

### Using exit code for success status, add new section to ini file

```bash
./goini --ini sample.ini --add-section new_section
echo $? # Output: 0 (and new_section will be added to sample.ini)

# Try adding again (should fail)
./goini --ini sample.ini --add-section new_section
echo $? # Output: 1
```

### Using exit code for success status, in section "section", add "key" with value "value"

```bash
./goini --ini sample.ini --section default --key new_setting --value 123 --add-key
echo $? # Output: 0 (and new_setting=123 will be added to [default])

# Try adding existing key (should fail)
./goini --ini sample.ini --section default --key user --value something --add-key
echo $? # Output: 1
```

### Using exit code for success status, in section "section", modify "key" with new value "value"

```bash
./goini --ini sample.ini --section default --key user --value NewUserValue --modify-key
echo $? # Output: 0 (and 'user' in [default] will be updated to 'NewUserValue')

# Try modifying non-existent key (should fail)
./goini --ini sample.ini --section default --key non_existent_key --value some_value --modify-key
echo $? # Output: 1
```

## Testing

There is a [test.sh](/test.sh) script that produces the following output that tests each of these examples.

```log
./test.sh
Building goini binary...
Clean successful: ./bin/goini-darwin-arm64
Build successful: ./bin/goini-darwin-arm64

TEST_01: Relying on exit code, does 'section' have 'key'?...PASSED
TEST_02: Relying on exit code, does 'section' have 'key'?...PASSED
TEST_03: Relying on exit code, does 'section' 'key' have value 'value'?...PASSED
TEST_04: Relying on exit code, does 'section' 'key' have value 'value'?...PASSED
TEST_05: Relying on exit code, are 'section1' and 'section2' present?...PASSED
TEST_06: Relying on exit code, are 'section1' and 'section2' present?...PASSED
TEST_07: Using STDOUT, return list of 'sections' in ini file...PASSED
TEST_08: Using STDOUT, return list of 'sections' in ini file...PASSED
TEST_09: Using STDOUT, return list of 'sections' in ini file...PASSED
TEST_10: Using STDOUT, return list of 'sections' in ini file...PASSED
TEST_11: Using STDOUT, return a list of keys in 'section' (by name)...PASSED
TEST_12: Using STDOUT, return a list of keys in 'section' (by name)...PASSED
TEST_13: Using STDOUT, return a list of key/values in 'section' (by name)...PASSED
TEST_14: Using STDOUT, return a list of key/values in 'section' (by name)...PASSED
TEST_15: Using exit code for success status, add new section to ini file...PASSED
TEST_16: Try adding again (expecting exit 1)...PASSED
TEST_17: Using exit code for success status, in section 'section', add 'key' with value 'value'...PASSED
TEST_18: Try adding existing key (expecting exit 1)...PASSED
TEST_19: Using exit code for success status, in section 'section', modify 'key' with new value 'value'...PASSED
TEST_20: Try modifying non-existent key (expecting exit 1)...PASSED

All 20 tests passed. PASS
```