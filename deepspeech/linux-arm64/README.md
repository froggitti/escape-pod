To build this..

see [this](https://deepspeech.readthedocs.io/en/v0.9.2/BUILDING.html)

```sh
# bazel build --workspace_status_command="bash native_client/bazel_workspace_status_cmd.sh" --config=monolithic --config=rpi3-armv8 --config=rpi3-armv8_opt -c opt --copt=-O3 --copt=-fvisibility=hidden --define=no_aws_support=true --define=no_gcp_support=true //native_client:libdeepspeech.so
```