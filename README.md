A very lazy renderfarm!

### testing

```go get github.com/kybin/lazyfarm # It will say some error messages, but it's OK.
cd lazyfarm
./build```

Then you wil have three excutables. "lazyfarm", "worker", "client".

`./lazyfarm`

It will print lazyfarm address("ip:port"). If it returns "192.168.0.4:8081".

With a new terminal, you can add workers to lazyfarm.

```# Add a worker
./worker -server="192.168.0.4:8081"```

And with a new terminal,

```./client -server="192.168.0.4:8081" -command="ls -al"

# or

./client -server=192.168.0.4:8081 -cmd="hython ./scene/houdini/test.hipnc -c hou.node('/out/mantra1').render(frame_range=({frame},{frame},{frame}))" -frames="1-24"```
