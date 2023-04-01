package scratch

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
	decl "github.com/lxn/walk/declarative"
)

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

func RunDalleTest() {
	var inputTextEdit *walk.TextEdit
	var runButton *walk.PushButton
	var openFileDialogButton *walk.PushButton
	var imagePath, maskPath string

	windowSize := decl.Size{
		Width:  50,
		Height: 50,
	}

	enableRunButton := func() {
		runButton.SetEnabled(len(inputTextEdit.Text()) > 0 && len(imagePath) > 0 && len(maskPath) > 0)
	}

	decl.MainWindow{
		Title:   "DALLE Image Editor",
		MinSize: windowSize,
		Layout:  decl.VBox{},
		Children: []decl.Widget{
			decl.Label{
				Text: "Enter a new description:",
			},
			decl.TextEdit{
				AssignTo:      &inputTextEdit,
				OnTextChanged: enableRunButton,
			},
			decl.PushButton{
				AssignTo: &openFileDialogButton,
				Text:     "Select Image and Mask",
				OnClicked: func() {
					imageFilter := "Image Files (*.png, *.jpg, *.jpeg)|*.png;*.jpg;*.jpeg"

					dlg := new(walk.FileDialog)
					dlg.Title = "Select Image"
					dlg.Filter = imageFilter

					if ok, err := dlg.ShowOpen(nil); err != nil {
						log.Printf("Error opening image file: %v", err)
					} else if !ok {
						log.Println("Image selection canceled")
					} else {
						imagePath = dlg.FilePath
					}

					dlg.Title = "Select Mask"
					if ok, err := dlg.ShowOpen(nil); err != nil {
						log.Printf("Error opening mask file: %v", err)
					} else if !ok {
						log.Println("Mask selection canceled")
					} else {
						maskPath = dlg.FilePath
					}
				},
			},
			decl.PushButton{
				AssignTo: &runButton,
				Text:     "Run",
				Enabled:  false,
				OnClicked: func() {
					description := inputTextEdit.Text()
					RunDalle(description, imagePath, maskPath)
				},
			},
		},
	}.Run()
}

func RunDalle(description string, imagePath string, maskPath string) {
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
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		fmt.Printf("Error reading image file: %v\n", err)
		os.Exit(1)
	}
	maskData, err := ioutil.ReadFile(maskPath)
	if err != nil {
		fmt.Printf("Error reading image file: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.Edit(
		context.Background(),
		imageData,
		maskData,
		description,
		dalle2.WithNumImages(1),
		dalle2.WithSize(dalle2.LARGE),
		dalle2.WithFormat(dalle2.URL),
	)

	if err != nil {
		log.Fatal(err)
	}
	for _, img := range resp.Data {
		openbrowser(img.Url)
	}
}
