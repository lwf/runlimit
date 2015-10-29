# runlimit
runlimit is a rate limiter for process invocations, primarily designed to be used with [daemontools](http://cr.yp.to/daemontools.html) and similar process supervisors (think runit, s6, etc).

## Building
Running `make` will produce a statically linked binary in the `target` directory. Debian infrastructure is provided in `debian/` and may be used to produce packages by way of your preferred method (debuild, pbuilder, etc).

## Usage
```shell
runlimit -sv-cmd cmd [-window-size duration] [-max-restarts restarts] [-metadata-dir dir] [-metadata-key key] prog [args...]
```

runlimit will allow `max-restarts` invocations of `prog` within a moving window of `window-size` duration. This is similar to how Upstart's respawn directive works. runlimit will invoke the command specified by the `sv-cmd` flag once the restart limit is reached. This command is expected to send SIGTERM to the runlimit process, as is conventional among daemontools and its clones.

Per-process state information is stored in `metadata-dir`. File names are determined by the `metadata-key` flag if specified or based on the current working directory of the runlimit process.

## Examples
Usage with daemontools:

```shell
# cat > /service/mydaemon/run <<EOF
#!/bin/sh

exec runlimit -window-size 10m -max-restarts 3 \
  -sv-cmd "svc -d `pwd`" mydaemon
EOF
```

Usage with runit:

```shell
# cat > /etc/service/mydaemon/run <<EOF
#!/bin/sh

exec runlimit -window-size 10m -max-restarts 3 \
  -sv-cmd "sv stop `pwd`" mydaemon
EOF
```

## License
Copyright 2015 Torbjörn Norinder.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
