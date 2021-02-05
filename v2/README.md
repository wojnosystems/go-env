# Overview

Reads environment variables and unmarshalls them into your structures.

## Install

`go get github.com/wojnosystems/go-env/v2`

This was originally intended to be used for reading environment variables into a structure and then validating those structures for a more intelligent configuration system for Go apps.

Because of the limitations of environment variables from the shell's perspective, we're limited to the following character set for valid environment variable names:

`[a-zA-Z_]+[a-zA-Z0-9_]*`

 POSIX says it doesn't care what the name is, as long as it doesn't contain an equal sign:
https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap08.html#tag_08 so we're limited by the shell
restrictions.

Thus, the naming conventions for environment variables being mapped to members in structures is structured:

Example variable naming scheme:

```go
type myStruct struct {
  Name      string     `env:"name"`
  PetNames  []string   `env:"pet_names"`
  Addresses addrStruct // no tag, assumes "Addresses"
}
type addrStruct struct {
  Street string     `env:"street"`
}

s := myStruct
Env{}.Unmarshall(&s)
```

The following environment variables can create the following data structure:

```bash
name=Chris pet_names_0_=Frankie pet_names_4_=Charlie Addresses_0_street="742 Evergreen Terrace" \
  Addresses_1_street="2001 Creaking Oak Drive" ./my-app
```

```go
s := myStruct{
  Name: "Chris",
  PetNames: []string{
    "Frankie", // index:0
    "",        // index:1
    "",        // Charlie was index 4, so we created 3 blank strings in between as this is an array
    "",        // index:3
    "Charlie", // index:4
  },
  Addresses: []addrStruct{
    {
      Street: "742 Evergreen Terrace",
    },
    {
      Street: "2001 Creaking Oak Drive",
    },
  }
}
```

Because the structure enforces the name and prevents duplicate names from being used, this scheme guarantees unique names for structure members and child members as long as you don't use tags to create collisions.

Every environment variable is a unique member of the destination. You cannot nest objects in the environment variable value.

# Naming

By default, variable names are assumed by the name provided in the destination GoLang structure. If you proivide a tag for "env" and set its value to a new name, then only the tag's env: name value will work.

## Separators

An underscore (_) separates the name of a field and its containing structure. "container_field".

Indexes are numbers and are always surrounded by underscores: "_5_"

## Tags

You can override the name of any field by using tags. You cannot, however, modify the separator between indices or fields.
