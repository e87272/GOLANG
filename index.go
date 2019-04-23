package main
import( "fmt"
		"github.com/kataras/iris"
		"github.com/kataras/iris/websocket"
		"./im"
		"./gm"
)

func main() {
	
    ws := websocket.New(websocket.Config{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
	})
	
    ws.OnConnection(handleConnection)

    app := iris.New()
    // register the server on an endpoint.
    // see the inline javascript code in the websockets.html, this endpoint is used to connect to the server.
    app.Get("/echo", ws.Handler())

	app.StaticWeb("/", "./client/index.html")
	
    // serve the javascript built'n client-side library,
    // see websockets.html script tags, this path is used.
    app.Get("/iris-ws.js", func(ctx iris.Context) {
        ctx.Write(websocket.ClientSource)
	})
    app.Get("./client/im.js", func(ctx iris.Context) {
		ctx.ServeFile("im.js", false)
	})
	
	app.Run(iris.Addr(":8080"))
	
	fmt.Printf("End\n")
}

func handleConnection(c websocket.Connection) {
    // Read events from browser
    c.On("gmsocket", func(msg string) {
        // Print the message to the console, c.Context() is the iris's http context.
		//fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		//fmt.Printf("%+v\n", c.Context())
        // Write message back to the client message owner:
		gm.ReceiveCustomPacket(c,msg);
		
	})
	
    c.On("imsocket", func(msg string) {
		//fmt.Printf("%s : %s\n", "ReceiveImPacket", msg)
		im.ReceiveImPacket(c,msg);
    })
}
