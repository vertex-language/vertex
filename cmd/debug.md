

sudo apt update
sudo apt install wabt

wasm2wat qsort.wasm -o out.wat

wasm-objdump -x output.wasm

go mod edit -replace github.com/vertex-language/ir=/Users/galaxy/Desktop/ir


gdb ./tcp_server

run

bt

(gdb) info registers

(gdb) x/20i $pc-20




# Install ltrace
sudo apt-get install -y ltrace

# Trace only the network calls — shows exact args to inet_pton, connect, socket
ltrace -e "socket+inet_pton+connect+htons+write+read" ./tcp_client2 2>&1

# Also check if the data is now correct in the new binary
strings ./tcp_client2 | grep -E "127\.0\.0\.1|GET /"