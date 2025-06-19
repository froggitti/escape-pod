
ssh escapepod@zecred.net "mkdir -p ~/escape-pod/bin"
ssh escapepod@zecred.net "mkdir -p ~/escape-pod/coqui"
scp -r modules/escape-pod-ui/dist escapepod@zecred.net:~/escape-pod/dist

# # Escape Pod Env
# scp image-builder/files/escape-pod.conf escapepod@zecred.net:~/escape-pod

# # Escape Pod Systemd Service Config
# scp image-builder/files/escape_pod.service escapepod@zecred.net:~/escape-pod

# # Escape Pod Binary
scp bin/escape-pod-linux-arm64 escapepod@zecred.net:~/escape-pod/bin/escape-pod-linux-arm64

# scp -r coqui/linux-tflite-aarch64 escapepod@zecred.net:~/escape-pod/coqui/linux-tflite-aarch64
# scp coqui/large_vocabulary.scorer escapepod@zecred.net:~/escape-pod/coqui
# scp coqui/model.tflite escapepod@zecred.net:~/escape-pod/coqui
# scp default-intents/default-intent-list.json escapepod@zecred.net:~