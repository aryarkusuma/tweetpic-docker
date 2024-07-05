  GNU nano 5.4                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              main.go                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       package main
import (
        "fmt"
        "log"
        "os"
        "path/filepath"
        "time"

        "github.com/go-rod/rod"
        "github.com/go-rod/rod/lib/launcher"
        "github.com/go-rod/rod/lib/proto"
        "github.com/gofiber/fiber/v2"
        "github.com/gofiber/fiber/v2/middleware/cors"
        "github.com/gofiber/fiber/v2/middleware/limiter"
        "github.com/gofiber/fiber/v2/middleware/logger"
        "github.com/google/uuid"
        "github.com/valyala/fasthttp"

)

type Response struct {
        TweetPicUrl string `json:"TweetPicUrl"`
}

func main() {

        app := fiber.New()
        app.Use(logger.New())


        // Rate limiter middleware
        app.Use("/tweetpic", limiter.New(limiter.Config{
                Max:        2,
                Expiration: 2 * time.Second,
                KeyGenerator: func(c *fiber.Ctx) string {
                        return c.IP() // limit based on IP address
                },
                LimitReached: func(c *fiber.Ctx) error {
                        return c.Status(fiber.StatusTooManyRequests).JSON(map[string]interface{}{
                                "Error": "Too many requests, slow down!",
                        })
                },
        }))

        app.Use(cors.New(cors.Config{
                AllowOrigins: "*",
        }))

        app.Get("/health-check", func(c *fiber.Ctx) error {
                return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
                        "Status": 200,
                })
        })

        app.Get("/image", serveImage)
        app.Get("/tweetpic", func(c *fiber.Ctx) error {
                tweetId := c.Query("id")
                data, err := TweetPic(tweetId)
                if data == "" || err != nil {
                        return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
                                "Error": "Tweet Doesnt Exist",
                        })
                }
                return c.JSON(&Response{
                        TweetPicUrl: data,
                })
        })


        // app.Get("/file-with-url-chars", func(c *fiber.Ctx) error {
        //         c.Set("Content-Type", "text/plain")
        //         return c.Send(indexHtml)
        //       })

        // Start server
        log.Fatal(app.Listen(":8080"))
}

func serveImage(c *fiber.Ctx) error {
        // Extract the image name from the URL path
        imageName := c.Query("name")
        if imageName == "" {
                return fiber.NewError(fiber.StatusBadRequest, "Image name not specified.")
        }
        // Path to the image file in the tweetpic/images directory
        imagePath := "./images" + filepath.Clean("/" + imageName)

        return c.SendFile(imagePath)
}
func TweetPic(tweetId string) (string, error) {
        //check if the tweet is exist
        statusCode, err := TweetCheck(tweetId)
        if statusCode != 200 || err != nil {
                return "", err
        }

        //write/save that screenshoted tweet to a filesystem (local)
        uuidv7, err := uuid.NewV7()
        if err != nil {
                fmt.Println("Error generating UUIDv7:", err)
                return "", err
        }

        imageName := uuidv7.String() + ".jpeg"
        imageUrl :=  "/image?name=" + imageName
        go TweetPicWorker(tweetId, imageName)
        return imageUrl, nil
}

func TweetCheck(tweetId string) (int, error) {

        url := "https://cdn.syndication.twimg.com/tweet-result?id=" + tweetId + "&lang=en&token=4cxws3gzf&dg9pqp=98faifuxq3ga&zbov9y=48gs8r425yop&z7yj98=eumlq8i1u2vc&5pwt7f=1mr6v8rxsmw1"
        statusCode, _, err := fasthttp.Get(nil, url)
        if err != nil {
                return 0, err
        }

        return statusCode, nil
}

func TweetPicWorker(tweetId string, imgName string) {
        //Get running chrome dev remote address
        u := launcher.MustResolveURL("")

        //open,load and rendering a page of the embeded api
        page := rod.New().ControlURL(u).MustConnect().MustPage("https://platform.twitter.com/embed/Tweet.html?id=" + tweetId).MustSetViewport(1920, 2000, 1, false)
        page.MustWaitStable()
        //selecting the parent node elements of embeded tweet that we want to screenshot
        tweetpic := page.MustElement("#app > div > div > div")

        //lists of elements path to remove the elements
        elements := []string{
                "#app > div > div > div > article > div.css-1dbjc4n.r-kzbkwu.r-1h8ys4a",
                "#app > div > div > div > article > div.css-1dbjc4n.r-18u37iz.r-kzbkwu > div.css-1dbjc4n.r-eqz5dr.r-1777fci.r-4amgru.r-1kh6xel > div > div > div > div > div",
                "#app > div > div > div > article > div.css-1dbjc4n.r-1habvwh.r-1ets6dv.r-5kkj8d.r-18u37iz.r-14gqq1x.r-1h8ys4a > div",
                "#app > div > div > div > article > div.css-1dbjc4n.r-1awozwy.r-18u37iz.r-1bymd8e > a > svg",
                "#app > div > div > div > article > div.css-1dbjc4n.r-1habvwh.r-1ets6dv.r-5kkj8d.r-18u37iz.r-14gqq1x.r-1h8ys4a > a:nth-child(2)",
        }

        //Loop through the elements paths array and execute MustRemove, removing them elements form the node (tweetpic)
        for _, selector := range elements {
                tweetpic.MustElement(selector).MustRemove()
        }

        // //screenshot the node(tweetpic) and make a name file identifier for the image
        ss, err := tweetpic.Screenshot(proto.PageCaptureScreenshotFormatJpeg, 90)
        if err != nil {
                log.Fatalf("Failed to screenshot: %v", err)
        }

        filepath := "./images/" + imgName
        err = os.WriteFile(filepath, ss, 0644)
        if err != nil {
                log.Fatalf("Failed to save screenshot: %v", err)
        }

        page.MustClose()

}






