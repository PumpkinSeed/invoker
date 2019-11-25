# Invoker

[![Go Report Card](https://goreportcard.com/badge/github.com/PumpkinSeed/invoker)](https://goreportcard.com/report/github.com/PumpkinSeed/invoker)

#### Run pre-defined commands inside your container/containers.

The main concept behind the solution is simplify and make faster command executions in docker containers.

#### Usage

```
go install github.com/PumpkinSeed/invoker
```

Create settings.json and move to a common location

```
/home/username/invoker-settings.json
{
  "containers": {
    "couchbase": [
      "api_couchbase_1"
    ]
  },
  "commands": {
    "ls": [
      "ls -ll",
      "echo end"
    ]
  }
}
```

Then add the `INVOKER_SETTINGS` environment variable with the location of settings.json

```
export INVOKER_SETTINGS=/home/username/invoker-settings.json
```

Use it like, where the first parameter always the container group and the second is the command group

```
invoker couchbase ls

// Output:
api_couchbase_1 [cc3a6f47227c8b084c8bb6546eaa0165a5631dee078c3197c2ea601f1042c290] ~ ls -ll
Htotal 100
drwxr-xr-x.   1 root root 4096 Sep 19 00:07 bin
drwxr-xr-x.   2 root root 4096 Apr 12  2016 boot
drwxr-xr-x.   5 root root  340 Nov 25 09:26 dev
-rwxrwxr-x.   1 root root 1930 Sep 19 00:06 entrypoint.sh
drwxr-xr-x.   1 root root 4096 Nov 25 09:26 etc
drwxr-xr-x.   2 root root 4096 Apr 12  2016 home
-rw-r--r--.   1 root root 9452 Nov 25 09:26 index.html
-rwxrwxr-x.   1 root root 2857 Oct 29 13:56 init.sh
drwxr-xr-x.   1 root root 4096 Sep 19 00:07 lib
drwxr-xr-x.   2 root root 4096 Sep  4 18:49 lib64
drwxr-xr-x.   2 root root 4096 Sep  4 18:47 media
drwxr-xr-x.   2 root root 4096 Sep  4 18:47 mnt
drwxr-xr-x.   1 root root 4096 Sep 19 00:07 opt
dr-xr-xr-x. 361 root root    0 Nov 25 09:26 proc
drwx------.   1 root root 4096 Nov 25 09:26 root
drwxr-xr-x.   1 root root 4096 Sep  4 18:49 run
drwxr-xr-x.   1 root root 4096 Sep 19 00:07 sbin
drwxr-xr-x.   2 root root 4096 Sep  4 18:47 srv
dr-xr-xr-x.  13 root root    0 Nov 25 07:15 sys
drwxrwxrwt.   1 root root 4096 Nov 25 09:26 tmp
drwxr-xr-x.   1 root root 4096 Sep  4 18:47 usr
drwxr-xr-x.   1 root root 4096 Sep  4 18:49 var
api_couchbase_1 [cc3a6f47227c8b084c8bb6546eaa0165a5631dee078c3197c2ea601f1042c290] ~ echo end
end
```

#### Settings file

- You can create container groups which can contain one or more container based on container's name. Get it with `docker ps`.
- You will use the alias (as the example shown `couchbase`) for call the invoker.

```
"containers": {
    "couchbase": [
      "api_couchbase_1"
    ]
  }
```

- You can create command groups which can have one or more commands.
- You will use the alias (as the example shown `ls`) for call the invoker.

```
"commands": {
    "ls": [
      "ls -ll",
      "echo end"
    ]
  }
```

#### Additional

There are extra environment variables to manipulate the behaviour of the program:

- `INVOKER_VERBOSE`: true // add detailed logs
- `INVOKER_SKIP_OUTPUT`: true // skip the print of commands output