package data

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/Andrew-peng/go-dalle2/dalle2"
	"github.com/joho/godotenv"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var inputTextEdit *walk.TextEdit
	var runButton *walk.PushButton

	MainWindow{
		Title:   "DALLE Image Editor",
		MinSize: Size{400, 200},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: "Enter a new description:",
			},
			TextEdit{
				AssignTo: &inputTextEdit,
			},
			PushButton{
				AssignTo: &runButton,
				Text:     "Run",
				OnClicked: func() {
					description := inputTextEdit.Text()
					RunDalle(description)
				},
			},
		},
	}.Run()
}

func RunDalle(description string) {
	// get secrets
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	apiKey := os.Getenv("DALLE_API_KEY")

	// connect a client to dalle 2
	client, err := dalle2.MakeNewClientV1(apiKey)
	if err != nil {
		log.Fatalf("Error initializing client: %s", err)
	}

	// read in an image
	imagePath := "scratch/data/profile_pic_tiny_heads.png"
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		fmt.Printf("Error reading image file: %v\n", err)
		os.Exit(1)
	}
	maskPath := "scratch/data/profile_pic_tiny_heads_mask.png"
	maskData, err := ioutil.ReadFile(maskPath)
	if err != nil {
		fmt.Printf("Error reading image file: %v\n", err)
		os.Exit(1)
	}

	// resp, err := client.Create(
	// 	context.Background(),
	// 	"4 Frogs in 1800s with suits playing poker, smoking cigars, drinking whiskey, and laughing. The setting is a Old western saloon.",
	// 	dalle2.WithNumImages(1),
	// 	dalle2.WithSize(dalle2.LARGE),
	// 	dalle2.WithFormat(dalle2.URL),
	// )

	resp, err := client.Edit(
		context.Background(),
		imageData,
		maskData,
		description,
		dalle2.WithNumImages(10),
		dalle2.WithSize(dalle2.LARGE),
		dalle2.WithFormat(dalle2.URL),
	)

	if err != nil {
		log.Fatal(err)
	}
	for _, img := range resp.Data {
		openbrowser(img.Url)
		fmt.Println("%s", img.Url)
	}

}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
