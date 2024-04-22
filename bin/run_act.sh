#!/bin/bash
#
# act is a local github actions workflow runner
#
# the note about the artifact server is based on the error (and fix)
# reported at https://www.thegoatinthemachine.net/2022/09/04/act_for_local_development.html

act --artifact-server-path /tmp/artifacts
